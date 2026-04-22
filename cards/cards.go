package cards

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/slikasp/dbmanfrags/config"
	"github.com/slikasp/dbmanfrags/database"
)

func makeCardURL(cardID int32) string {
	url := fmt.Sprintf("https://fimgs.net/mdimg/perfume-social-cards/en-p_c_%d.jpeg", cardID)
	return url
}

func makeFilePath(cardID int32) string {
	path := fmt.Sprintf("cards/en/p_c_%d.jpeg", cardID)
	return path
}

func downloadCard(cardID int32) (database.AddCardParams, error) {
	// Start with card ID, URL, file path and HasCard as false
	card := database.AddCardParams{
		FragranticaID: cardID,
	}
	card.Url = makeCardURL(cardID)
	card.Image = makeFilePath(cardID)
	card.HasCard = false

	// Try getting the image, return initial details if we can't
	resp, err := http.Get(card.Url)
	if err != nil {
		return card, fmt.Errorf("Get request failed for %s: %w", card.Url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return card, fmt.Errorf("No card for ID %d", cardID)
	}

	// If response is ok we set HasCard to true, after this we don't expect more errors so need to handle that
	card.HasCard = true

	file, err := os.Create(card.Image)
	if err != nil {
		return card, fmt.Errorf("Could not create image %s: %w", card.Image, err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println(err)
		return card, fmt.Errorf("Could not write image %s: %w", card.Image, err)
	}

	// log.Printf("Card downloaded for ID %d\n", cardID)

	return card, nil
}

// Only use this for a fresh database because it will redownload all exisiting cards too.
// Might be used to force update and find new cards that weren't there before.
// Populates DB even with non existing fragrances.
//
// TODO: try to make this faster
func DownloadAllCards(frags *config.Frags) error {
	startCardID := int32(1)
	endCardID := frags.LastID

	for id := startCardID; id <= endCardID; id++ {
		// Try downloading the card (existing or not)
		card, err := downloadCard(id)
		// If card was found, but error still returned stop execution
		if err != nil && card.HasCard {
			return fmt.Errorf("Card download failed for ID %d: %w", id, err)
		}

		// Check if card already exists
		_, err = frags.DB.GetCard(context.Background(), id)
		if err != nil {
			// ID doesn't exist
			if errors.Is(err, sql.ErrNoRows) {
				// Add new card to the database
				_, err = frags.DB.AddCard(context.Background(), card)
				if err != nil {
					return fmt.Errorf("Adding card to database failed for ID %d: %w", id, err)
				}
				log.Printf("New card added, ID:%d, URL:%s", id, card.Image)
			} else {
				// Real error
				return fmt.Errorf("Could not get card from the database with ID %d: %w", id, err)
			}
		} else {
			// Update card if exists
			_, err = frags.DB.UpdateCard(context.Background(), database.UpdateCardParams{
				FragranticaID: card.FragranticaID,
				Image:         card.Image,
				HasCard:       card.HasCard,
			})
			if err != nil {
				return fmt.Errorf("Could not update card with ID %d: %w", id, err)
			}
		}
	}

	return nil
}

// Goes through all cards marked with HasCard == false and retries them
func CheckMissingCards(frags *config.Frags) error {
	cardIDs, err := frags.DB.GetMissingCardIDs(context.Background())
	if err != nil {
		return fmt.Errorf("Failed getting missing cards from database: %w", err)
	}

	log.Printf("Missing cards: %d", len(cardIDs))
	cardsAdded := 0

	for _, id := range cardIDs {
		card, err := downloadCard(id)
		if card.HasCard {
			// Card found, but not downloaded -> return download error
			if err != nil {
				return fmt.Errorf("Card download failed for ID %d: %w", id, err)
			} else {
				// Card found -> update DB
				_, err = frags.DB.UpdateCard(context.Background(), database.UpdateCardParams{
					FragranticaID: card.FragranticaID,
					Image:         card.Image,
					HasCard:       card.HasCard,
				})
				if err != nil {
					return fmt.Errorf("Could not update card with ID %d: %w", id, err)
				}
				cardsAdded += 1
				log.Printf("New card added, ID:%d, URL:%s", id, card.Image)
			}
		} else {
			// TODO: create a custom error type for card downloads and log only actual errors
			// these are 100% no card found so far
			// Log error (most likely card not found + ID)
			// log.Println(err)
		}

		// Proceed to next id if no card found
	}

	log.Printf("New cards added: %d / %d", cardsAdded, len(cardIDs))
	return nil
}

// TODOS

// Need to remove all cards from the database that have no card after last ID with a card
// Then run this, because it will create card entries in DB whether they are available or not
func FindNewCards(frags *config.Frags, cardsToCheck int) error {
	id, err := frags.DB.GetLastCardID(context.Background())
	if err != nil {
		return fmt.Errorf("Failed getting last card ID from database: %w", err)
	}

	cardsFound := 0
	currentID := id + 1
	lastFound := id

	// Check cardsToCheck number of cards after the last found
	for currentID <= lastFound+int32(cardsToCheck) {
		// Try downloading the card (existing or not)
		card, err := downloadCard(currentID)

		// log.Printf("checking: %d", id)

		if card.HasCard {
			// Card found -> look for 100 more cards
			cardsFound += 1
			lastFound = currentID
			// If card was found, but error still returned -> stop execution
			if err != nil {
				return fmt.Errorf("Card download failed for ID %d: %w", currentID, err)
			}

			// If found, no errors -> Check if card already exists
			_, err = frags.DB.GetCard(context.Background(), currentID)
			if err != nil {
				// doesn't exist -> Add new card to the database
				if errors.Is(err, sql.ErrNoRows) {
					_, err = frags.DB.AddCard(context.Background(), card)
					if err != nil {
						return fmt.Errorf("Adding card to database failed for ID %d: %w", currentID, err)
					}
					log.Printf("New card added, ID:%d, URL:%s", currentID, card.Image)
				} else {
					// Real error -> stop execution
					return fmt.Errorf("Could not get card from the database with ID %d: %w", currentID, err)
				}
			} else {
				// Update card if it exists in the database, card might have appeared (and we just downloaded it)
				_, err = frags.DB.UpdateCard(context.Background(), database.UpdateCardParams{
					FragranticaID: card.FragranticaID,
					Image:         card.Image,
					HasCard:       card.HasCard,
				})
				if err != nil {
					return fmt.Errorf("Could not update card with ID %d: %w", currentID, err)
				}
			}

			// Card already exists in database -> do nothing, CheckMissingCards will recheck
		}
		// Card not foun -> proceed to next ID
		currentID += 1

	}

	log.Printf("Last found card: %d. Cards found: %d\n", lastFound, cardsFound)

	return nil
}
