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
func DownloadAllCards(state *config.State) error {
	startCardID := int32(1)
	endCardID := state.LastID

	for id := startCardID; id <= endCardID; id++ {
		// Try downloading the card (existing or not)
		card, err := downloadCard(id)
		// If card was found, but error still returned stop execution
		if err != nil && card.HasCard {
			return fmt.Errorf("Card download failed for ID %d: %w", id, err)
		}

		// Check if card already exists
		_, err = state.DB.GetCard(context.Background(), id)
		if err != nil {
			// ID doesn't exist
			if errors.Is(err, sql.ErrNoRows) {
				// Add new card to the database
				_, err = state.DB.AddCard(context.Background(), card)
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
			_, err = state.DB.UpdateCard(context.Background(), database.UpdateCardParams{
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
func CheckMissingCards(state *config.State) error {
	cardIDs, err := state.DB.GetMissingCardIDs(context.Background())
	if err != nil {
		return fmt.Errorf("Failed getting missing cards from database: %w", err)
	}

	for _, id := range cardIDs {
		card, err := downloadCard(id)
		if card.HasCard {
			// Card found, but not downloaded -> return download error
			if err != nil {
				return fmt.Errorf("Card download failed for ID %d: %w", id, err)
			} else {
				// Card found -> update DB
				_, err = state.DB.UpdateCard(context.Background(), database.UpdateCardParams{
					FragranticaID: card.FragranticaID,
					Image:         card.Image,
					HasCard:       card.HasCard,
				})
				if err != nil {
					return fmt.Errorf("Could not update card with ID %d: %w", id, err)
				}
				log.Printf("New card added, ID:%d, URL:%s", id, card.Image)
			}
		} else {
			// Log error (most likely card not found + ID)
			log.Println(err)
		}

		// Proceed to next id if no card found
	}

	return nil
}

// TODOS

// Need to remove all cards from the database that have no card after last ID with a card
// Then run this, because it will create card entries in DB whether they are available or not
func FindNewCards(state *config.State, cardsToCheck int) error {
	id, err := state.DB.GetLastCardID(context.Background())
	if err != nil {
		return fmt.Errorf("Failed getting last card ID from database: %w", err)
	}

	cardsFound := 0
	lastFound := id

	// Check 100 cards after the last one
	for id <= lastFound+int32(cardsToCheck) {
		// Try downloading the card (existing or not)
		card, err := downloadCard(id)

		if card.HasCard {
			// Card found -> look for 100 more cards
			cardsFound++
			lastFound = id
			// If card was found, but error still returned -> stop execution
			if err != nil {
				return fmt.Errorf("Card download failed for ID %d: %w", id, err)
			}

			// If found, no errors -> Check if card already exists
			_, err = state.DB.GetCard(context.Background(), id)
			if err != nil {
				// doesn't exist -> Add new card to the database
				if errors.Is(err, sql.ErrNoRows) {
					_, err = state.DB.AddCard(context.Background(), card)
					if err != nil {
						return fmt.Errorf("Adding card to database failed for ID %d: %w", id, err)
					}
					log.Printf("New card added, ID:%d, URL:%s", id, card.Image)
				} else {
					// Real error -> stop execution
					return fmt.Errorf("Could not get card from the database with ID %d: %w", id, err)
				}
			} else {
				// Update card if it exists in the database, card might have appeared (and we just downloaded it)
				_, err = state.DB.UpdateCard(context.Background(), database.UpdateCardParams{
					FragranticaID: card.FragranticaID,
					Image:         card.Image,
					HasCard:       card.HasCard,
				})
				if err != nil {
					return fmt.Errorf("Could not update card with ID %d: %w", id, err)
				}
			}

			// Card already exists in database -> do nothing, CheckMissingCards will recheck
		}
		// Card not foun -> proceed to next ID
		id++

	}

	log.Printf("Last found card: %d. Cards found: %d\n", lastFound, cardsFound)

	return nil
}
