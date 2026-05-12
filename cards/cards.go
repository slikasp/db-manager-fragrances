package cards

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/slikasp/dbmanfrags/config"
	"github.com/slikasp/dbmanfrags/database"
)

// Goes through all cards marked with HasCard == false and retries them
func CheckMissingCards(db *config.Database) error {
	cardIDs, err := db.Queries.GetMissingCardIDs(context.Background())
	if err != nil {
		db.Logger.Error("get missing cards from database", "error", err)
		return err
	}

	db.Logger.Info("missing cards",
		"number", len(cardIDs),
	)
	cardsAdded := 0

	for _, id := range cardIDs {
		card, err := downloadCard(id)
		if card.HasCard {
			if err != nil {
				// Card found but not downloaded -> return download error
				db.Logger.Error("downloadCard", "error", err)
				return err
			} else {
				// Card found -> update DB
				_, err = db.Queries.UpdateCard(context.Background(), database.UpdateCardParams{
					FragranticaID: card.FragranticaID,
					Image:         card.Image,
					HasCard:       card.HasCard,
				})
				if err != nil {
					db.Logger.Error("update card", "id", id, "error", err)
					return err
				}
				cardsAdded += 1
				db.Logger.Info("card added",
					"id", id,
					"location", card.Image,
				)
			}
			// } else {
			// TODO: create a custom error type for card downloads and log only actual errors
			// these are 100% no card found so far
			// Log error (most likely card not found + ID)
		}

		// Proceed to next id if no card found
	}
	db.Logger.Info("cards added",
		"number", cardsAdded,
	)
	return nil
}

// Goes through all cards marked with HasCard == true and checks if they exist in local storage
func CheckExistingCards(db *config.Database) error {
	cardIDs, err := db.Queries.GetExistingCardIDs(context.Background())
	if err != nil {
		db.Logger.Error("get existing cards from database", "error", err)
		return err
	}

	localCards, err := listCardFiles()

	diff := listDiff(cardIDs, localCards)
	if len(diff) > 0 {
		db.Logger.Info("cards in database", "expected", len(cardIDs), "found", len(localCards))

		for _, id := range diff {
			card, err := downloadCard(id)
			if card.HasCard {
				// Card found
				if err != nil {
					// but not downloaded -> return download error
					db.Logger.Error("downloadCard", "error", err)
					return err
				} else {
					// Card found -> update DB
					_, err = db.Queries.UpdateCard(context.Background(), database.UpdateCardParams{
						FragranticaID: card.FragranticaID,
						Image:         card.Image,
						HasCard:       card.HasCard,
					})
					if err != nil {
						db.Logger.Error("update card", "id", id, "error", err)
						return err
					}
					db.Logger.Info("card redownloaded", "path", makeFilePath(id))
				}
			}
		}
	}

	return nil
}

// TODO
// Need to remove all cards from the database that have no card after last ID with a card
// Then run this, because it will create card entries in DB whether they are available or not
func FindNewCards(db *config.Database, cardsToCheck int) error {
	id, err := db.Queries.GetLastCardID(context.Background())
	if err != nil {
		db.Logger.Error("get last card ID from database", "error", err)
		return err
	}

	cardsFound := 0
	currentID := id + 1
	lastFound := id

	// Check cardsToCheck number of cards after the last found
	for currentID <= lastFound+int32(cardsToCheck) {
		// Try downloading the card (existing or not)
		card, err := downloadCard(currentID)
		db.Logger.Debug("checking card", "id", id, "url", card.Url)
		if card.HasCard {
			// Card found -> look for 100 more cards
			cardsFound += 1
			lastFound = currentID
			// If card was found, but error still returned -> stop execution
			if err != nil {
				db.Logger.Error("downloadCard", "error", err)
				return err
			}

			// If found, no errors -> Check if card already exists
			_, err = db.Queries.GetCard(context.Background(), currentID)
			if err != nil {
				// doesn't exist -> Add new card to the database
				if errors.Is(err, sql.ErrNoRows) {
					_, err = db.Queries.AddCard(context.Background(), card)
					if err != nil {
						db.Logger.Error("add card to database", "id", currentID, "error", err)
						return err
					}
					db.Logger.Info("new card added", "id", currentID, "path", card.Image)
				} else {
					// Real error -> stop execution
					db.Logger.Error("get card from database", "id", currentID, "error", err)
					return err
				}
			} else {
				// Update card if it exists in the database, card might have appeared (and we just downloaded it)
				_, err = db.Queries.UpdateCard(context.Background(), database.UpdateCardParams{
					FragranticaID: card.FragranticaID,
					Image:         card.Image,
					HasCard:       card.HasCard,
				})
				if err != nil {
					db.Logger.Error("update card", "id", currentID, "error", err)
					return err
				}
			}
			db.Logger.Debug("card already in database", "id", id)
		}
		// Card not foun -> proceed to next ID
		currentID += 1

	}
	db.Logger.Info("cards found", "last_found", lastFound, "found", cardsFound)
	return nil
}

// Redownloads a card (to be used as part of fragrance update)
func RedownloadCard(db *config.Database, id int32) error {
	card, err := downloadCard(id)
	if card.HasCard {
		if err != nil {
			// card found,  but not downloaded -> return download error
			db.Logger.Error("downloadCard", "error", err)
			return err
		} else {
			// card found & downloaded -> update DB
			_, err = db.Queries.UpdateCard(context.Background(), database.UpdateCardParams{
				FragranticaID: card.FragranticaID,
				Image:         card.Image,
				HasCard:       card.HasCard,
			})
			if err != nil {
				db.Logger.Error("update card", "id", id, "error", err)
				return err
			}
			db.Logger.Info("card updated", "id", id, "path", card.Image)
		}
	}
	return nil
}

// Only use this for a fresh database because it will redownload all exisiting cards too.
// Might be used to force update and find new cards that weren't there before.
// Populates DB even with non existing fragrance cards.
func DownloadAllCards(db *config.Database) error {
	startCardID := int32(1)
	// Set ID of last card from the database
	endCardID, err := db.Queries.GetLastCardID(context.Background())
	if err != nil {
		return err
	}

	for id := startCardID; id <= endCardID; id++ {
		// Try downloading the card (existing or not)
		card, err := downloadCard(id)
		// If card was found, but error still returned stop execution
		if err != nil && card.HasCard {
			db.Logger.Error("downloadCard", "error", err)
			return err
		}
		// Check if card already exists
		_, err = db.Queries.GetCard(context.Background(), id)
		if err != nil {
			// ID doesn't exist
			if errors.Is(err, sql.ErrNoRows) {
				// Add new card to the database
				_, err = db.Queries.AddCard(context.Background(), card)
				if err != nil {
					db.Logger.Error("add card to database", "id", id, "error", err)
					return err
				}
				db.Logger.Info("new card added", "id", id, "path", card.Image)
			} else {
				// Real error
				db.Logger.Error("get card from database", "id", id, "error", err)
				return err
			}
		} else {
			// Update card if exists
			_, err = db.Queries.UpdateCard(context.Background(), database.UpdateCardParams{
				FragranticaID: card.FragranticaID,
				Image:         card.Image,
				HasCard:       card.HasCard,
			})
			if err != nil {
				db.Logger.Error("update card", "id", id, "error", err)
				return err
			}
		}
	}

	return nil
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

func makeCardURL(cardID int32) string {
	url := fmt.Sprintf("https://fimgs.net/mdimg/perfume-social-cards/en-p_c_%d.jpeg", cardID)
	return url
}

func makeFilePath(cardID int32) string {
	path := fmt.Sprintf("cards/en/p_c_%d.jpeg", cardID)
	return path
}

func listCardFiles() ([]int32, error) {
	path := "cards/en/"
	var ids []int32

	files, err := os.ReadDir(path)
	if err != nil {
		return ids, err
	}

	for _, f := range files {
		name := f.Name()

		// Skip directories just in case
		if f.IsDir() {
			continue
		}

		// Expected format: p_c_<ID>.jpeg
		if !strings.HasPrefix(name, "p_c_") || !strings.HasSuffix(name, ".jpeg") {
			continue // skip unexpected files
		}

		idStr := strings.TrimPrefix(name, "p_c_")
		idStr = strings.TrimSuffix(idStr, ".jpeg")

		id, err := strconv.ParseInt(idStr, 10, 32)
		if err != nil {
			continue // skip invalid IDs
		}

		ids = append(ids, int32(id))
	}

	return ids, nil
}

func listDiff(a, b []int32) []int32 {
	setB := make(map[int32]struct{}, len(b))
	for _, id := range b {
		setB[id] = struct{}{}
	}

	result := make([]int32, 0, len(a))
	for _, id := range a {
		if _, exists := setB[id]; !exists {
			result = append(result, id)
		}
	}

	return result
}
