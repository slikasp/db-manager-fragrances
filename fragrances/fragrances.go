package fragrances

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/slikasp/dbmanfrags/cards"
	"github.com/slikasp/dbmanfrags/config"
	"github.com/slikasp/dbmanfrags/database"
)

// Only use this on a fresh database.
// Not useful otherwise, will only check existing links.
func CheckAllLinks(state *config.State) error {
	ids, err := state.DB.GetExistingCardIDs(context.Background())
	if err != nil {
		return fmt.Errorf("Failed getting IDs from database: %w", err)
	}

	for _, id := range ids {
		card, err := state.DB.GetCard(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Could not get card by ID %d from database: %w", id, err)
		}

		urlCard, err := cards.GetLinkFromCard(card.Image)
		if err != nil {
			return fmt.Errorf("Failed parsing QR from image %s: %w", card.Image, err)
		}

		urlFrag, err := state.DB.GetFragranceLink(context.Background(), id)
		if err != nil {
			// No fragrance with this ID -> add new
			if errors.Is(err, sql.ErrNoRows) {
				err = state.DB.AddFragranceLink(context.Background(), database.AddFragranceLinkParams{
					FragranticaID: id,
					Url: sql.NullString{
						String: urlCard,
						Valid:  true,
					},
				})
				if err != nil {
					return fmt.Errorf("Failed adding fragrance with ID %d: %w", id, err)
				}
				log.Printf("Added new fragrance, ID:%d, URL:%s", id, urlCard)
			} else {
				// Real error
				return fmt.Errorf("Could not get fragrance link from database: %w", err)
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

// Checks card that doesn't have a linked fragrance to it's ID
func AddMissingFragrances(state *config.State) error {
	ids, err := state.DB.GetMissingFragranceIDs(context.Background())
	if err != nil {
		return fmt.Errorf("Failed getting IDs from database: %w", err)
	}

	log.Printf("New cards found: %d", len(ids))
	fragrancesAdded := 0

	for _, id := range ids {
		card, err := state.DB.GetCard(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Could not get card by ID %d from database: %w", id, err)
		}

		urlCard, err := cards.GetLinkFromCard(card.Image)
		if err != nil {
			// If QR decoding fails, set has_card to false, the card will be redownloaded on the next check
			// This is required because some cards are generated with empty QR codes
			log.Printf("Failed decoding card ID %d", id)
			err = state.DB.InvalidateCard(context.Background(), id)
			if err != nil {
				return fmt.Errorf("Failed setting has_card for card ID %d to false: %w", id, err)
			}
			continue
		}

		// Should we need a check for existing fragrances?

		err = state.DB.AddFragranceLink(context.Background(), database.AddFragranceLinkParams{
			FragranticaID: id,
			Url: sql.NullString{
				String: urlCard,
				Valid:  true,
			},
		})
		if err != nil {
			return fmt.Errorf("Failed adding fragrance with ID %d: %w", id, err)
		}
		fragrancesAdded += 1
		log.Printf("Added new fragrance, ID:%d, URL:%s", id, urlCard)
	}

	log.Printf("Fragrances added: %d / %d", len(ids), fragrancesAdded)

	return nil
}
