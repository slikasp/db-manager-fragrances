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
		fmt.Println(err)
		return card, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("No card for ID %d", cardID)
		fmt.Println(err)
		return card, err
	}
	// If response is ok we set HasCard to true, after this we don't expect more errors so need to handle that
	card.HasCard = true

	file, err := os.Create(card.Image)
	if err != nil {
		fmt.Println(err)
		return card, err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		fmt.Println(err)
		return card, err
	}

	fmt.Printf("Card downloaded for ID %d\n", cardID)

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
			return err
		}

		// Check if card already exists
		_, err = state.DB.GetCard(context.Background(), id)
		if err != nil {
			// ID doesn't exist
			if errors.Is(err, sql.ErrNoRows) {
				// Add new card to the database
				_, err = state.DB.AddCard(context.Background(), card)
				if err != nil {
					return err
				}
				log.Printf("New card added, ID:%d, URL:%s", id, card.Image)
			} else {
				// Real error
				return err
			}
		} else {
			// Update card if exists
			_, err = state.DB.UpdateCard(context.Background(), database.UpdateCardParams{
				FragranticaID: card.FragranticaID,
				Image:         card.Image,
				HasCard:       card.HasCard,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Goes through all cards marked with HasCard == false and retries them
func CheckMissingCards(state *config.State) error {
	cardIDs, err := state.DB.GetMissingCardIDs(context.Background())
	if err != nil {
		return err
	}

	for _, id := range cardIDs {
		card, err := downloadCard(id)
		if card.HasCard {
			// Card found, but not downloaded -> return download error
			if err != nil {
				return err
			} else {
				// Card found -> update DB
				_, err = state.DB.UpdateCard(context.Background(), database.UpdateCardParams{
					FragranticaID: card.FragranticaID,
					Image:         card.Image,
					HasCard:       card.HasCard,
				})
				log.Printf("New card added, ID:%d, URL:%s", id, card.Image)
				if err != nil {
					return err
				}
			}
		}
		// Proceed to next id if no card found
	}

	return nil
}

// Only use this on a fresh database.
// Not useful otherwise, will only check existing links
// TODO: rework this so only non existent frags and existing cards are decoded
func CheckAllLinks(state *config.State) error {
	ids, err := state.DB.GetExistingCardIDs(context.Background())
	if err != nil {
		return fmt.Errorf("Failed getting IDs from database: %s", err)
	}

	for _, id := range ids {
		card, err := state.DB.GetCard(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Could not get card by ID from database: %s", err)
		}

		urlCard, err := getLinkFromCard(card.Image)
		if err != nil {
			return fmt.Errorf("Failed parsing QR from image: %s", err)
		}

		urlFrag, err := state.DB.GetFragranceLink(context.Background(), id)
		if err != nil {
			// No fragrance with this ID -> add new
			if errors.Is(err, sql.ErrNoRows) {
				state.DB.AddFragranceLink(context.Background(), database.AddFragranceLinkParams{
					FragranticaID: id,
					Url: sql.NullString{
						String: urlCard,
						Valid:  true,
					},
				})
				log.Printf("Added new fragrance, ID:%d, URL:%s", id, urlCard)
			} else {
				// Real error
				return fmt.Errorf("Could not get fragrance link from database: %s", err)
			}
		} else {
			// Compare links if fragrance is already in database
			if urlCard != urlFrag.String {
				return fmt.Errorf("URL mismatch (card:frag): %s:%s", urlCard, urlFrag.String)
			}
			log.Printf("Decoded link matches existing fragrance, ID:%d", id)
		}
	}

	return nil
}

// TODOS

// func - look for new cards: get max number with card, check for next ~100 cards, if found, update the max number
