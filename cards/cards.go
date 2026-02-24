package cards

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/slikasp/dbmanfrags/config"
	"github.com/slikasp/dbmanfrags/database"
)

type Log struct {
	Checked  int32
	Found    []int32
	NotFound []int32
}

func (l *Log) card(id int32, found bool) {
	l.Checked = id
	if found {
		l.Found = append(l.Found, id)
	} else {
		l.NotFound = append(l.NotFound, id)
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
		err = errors.New("Unexpected response code")
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

	return card, nil
}

func GetAllCards(state config.State) Log {
	log := Log{}

	startCardID := state.CurrentID
	endCardID := state.LastID

	for cardID := startCardID; cardID <= endCardID; cardID++ {
		card, err := downloadCard(cardID)
		log.card(cardID, card.HasCard)
		// If card was found, but error still returned stop execution
		if err != nil && card.HasCard {
			return log
		}
		log.card(cardID, true)
		// Add card details to the database
		state.DB.AddCard(context.Background(), card)

		// Don't want to overwhelm services
		time.Sleep(1 * time.Second)
	}

	return log
}
