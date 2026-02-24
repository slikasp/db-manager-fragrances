package cards

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/slikasp/dbmanfrags/config"
	"github.com/slikasp/dbmanfrags/database"
)

type Log struct {
	Checked  int32
	Found    int32
	NotFound int32
}

func (l *Log) card(id int32, found bool) {
	l.Checked = id
	if found {
		l.Found += 1
	} else {
		l.NotFound += 1
	}
}

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

// Only use this for a fresh database because it will redownload all exisiting cards too
// Might be used to force update and find new cards that weren't there before
func DownloadAllCards(state *config.State) Log {
	log := Log{}

	startCardID := state.CurrentID
	endCardID := state.LastID

	for cardID := startCardID; cardID <= endCardID; cardID++ {
		// Keep track of card being worked on
		state.CurrentID = cardID

		// Try downloading the card (existing or not)
		card, err := downloadCard(cardID)
		log.card(card.FragranticaID, card.HasCard)
		// If card was found, but error still returned stop execution
		if err != nil && card.HasCard {
			fmt.Println(err)
			return log
		}

		// Check if card already exists
		_, err = state.DB.GetCard(context.Background(), cardID)
		if err != nil {
			// ID doesn't exist
			if errors.Is(err, sql.ErrNoRows) {
				// Try adding new card to the database
				_, err = state.DB.AddCard(context.Background(), card)
				if err != nil {
					fmt.Println(err)
					return log
				}
			} else {
				// Real error
				fmt.Println(err)
				return log
			}
		}

		// Card exists, update instead
		_, err = state.DB.UpdateCard(context.Background(), database.UpdateCardParams{
			FragranticaID: card.FragranticaID,
			Image:         card.Image,
			HasCard:       card.HasCard,
		})
		if err != nil {
			fmt.Println(err)
			return log
		}

		// If you don't want to overwhelm services
		// time.Sleep(100 * time.Millisecond)
	}

	return log
}
