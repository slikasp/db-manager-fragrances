package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Card struct {
	ID    int32
	URL   string
	Path  string
	Found bool
}

func makeCardURL(cardID int32) string {
	url := fmt.Sprintf("https://fimgs.net/mdimg/perfume-social-cards/en-p_c_%d.jpeg", cardID)
	return url
}

func makeFilePath(cardID int32) string {
	path := fmt.Sprintf("cards/%d.jpeg", cardID)
	return path
}

func DownloadCard(cardID int32) (Card, error) {
	card := Card{
		ID: cardID,
	}
	card.URL = makeCardURL(cardID)
	card.Path = makeFilePath(cardID)

	resp, err := http.Get(card.URL)
	if err != nil {
		return card, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return card, errors.New("Unexpected response code")
	}

	card.Found = true

	file, err := os.Create(card.Path)
	if err != nil {
		return card, err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return card, err
	}

	return card, nil
}
