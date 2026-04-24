package fragrances

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"strings"
	"time"

	"github.com/slikasp/dbmanfrags/cards"
	"github.com/slikasp/dbmanfrags/config"
	"github.com/slikasp/dbmanfrags/database"
	"github.com/slikasp/dbmanfrags/fragrantica"
)

// Only use this on a fresh database.
// Not useful otherwise, will only check existing links.
func CheckAllLinks(frags *config.Frags) error {
	ids, err := frags.DB.GetExistingCardIDs(context.Background())
	if err != nil {
		return fmt.Errorf("Failed getting IDs from database: %w", err)
	}

	for _, id := range ids {
		card, err := frags.DB.GetCard(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Could not get card by ID %d from database: %w", id, err)
		}

		urlCard, err := cards.GetLinkFromCard(card.Image)
		if err != nil {
			return fmt.Errorf("Failed parsing QR from image %s: %w", card.Image, err)
		}

		urlFrag, err := frags.DB.GetFragranceLink(context.Background(), id)
		if err != nil {
			// No fragrance with this ID -> add new
			if errors.Is(err, sql.ErrNoRows) {
				err = frags.DB.AddFragranceLink(context.Background(), database.AddFragranceLinkParams{
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
func AddMissingFragrances(frags *config.Frags) error {
	ids, err := frags.DB.GetMissingFragranceIDs(context.Background())
	if err != nil {
		return fmt.Errorf("Failed getting IDs from database: %w", err)
	}

	log.Printf("New cards found: %d", len(ids))
	fragrancesAdded := 0

	for _, id := range ids {
		card, err := frags.DB.GetCard(context.Background(), id)
		if err != nil {
			return fmt.Errorf("Could not get card by ID %d from database: %w", id, err)
		}

		urlCard, err := cards.GetLinkFromCard(card.Image)
		if err != nil {
			// If QR decoding fails, set has_card to false, the card will be redownloaded on the next check
			// This is required because some cards are generated with empty QR codes
			log.Printf("Failed decoding card ID %d", id)
			err = frags.DB.InvalidateCard(context.Background(), id)
			if err != nil {
				return fmt.Errorf("Failed setting has_card for card ID %d to false: %w", id, err)
			}
			continue
		}

		// Should we need a check for existing fragrances?

		err = frags.DB.AddFragranceLink(context.Background(), database.AddFragranceLinkParams{
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

func AddMissingDetails(frags *config.Frags) error {
	// for testing only
	// fragIDs := []int32{920}
	// var err error

	fragIDs, err := frags.DB.GetFragrancesWithoutDetails(context.Background())
	if err != nil {
		return fmt.Errorf("Failed getting IDs from database: %w", err)
	}

	numFrags := len(fragIDs)
	log.Printf("Fragrances to update: %d", numFrags)

	scraper, err := fragrantica.NewScraper()
	if err != nil {
		return fmt.Errorf("Failed creating scraper: %s", err)
	}

	count := 0
	for _, id := range fragIDs {
		count += 1
		_, err = updateFragranceDetails(frags, scraper, id)
		if err != nil {
			return fmt.Errorf("Failed updating fragrance with ID %d: %w", id, err)
		}
		log.Printf("Updated fragrance with ID:%d (%d/%d)", id, count, numFrags)
		// TODO: figure out the delay to not get locked out
		spamDelay(5, 10)
	}
	return nil
}

func updateFragranceDetails(frags *config.Frags, scrapper *fragrantica.Scraper, id int32) (fragrantica.FragranceParams, error) {
	// ID - already in DB
	// URL - already in DB

	// get and parse url for name and brand
	link, err := frags.DB.GetFragranceLink(context.Background(), id)
	if err != nil {
		return fragrantica.FragranceParams{}, fmt.Errorf("Failed getting fragrance link for ID %d: %w", id, err)
	}
	if !link.Valid {
		return fragrantica.FragranceParams{}, fmt.Errorf("Null url for ID %d", id)
	}
	name, brand, err := parseURL(link.String)
	if err != nil {
		return fragrantica.FragranceParams{}, fmt.Errorf("Failed parsing fragrance url '%s': %w", link.String, err)
	}

	// call ParsePage(url) for website parameters
	// this probably needs a delay to not be detected by fragrantica
	params, err := scrapper.ParsePageParams(link.String)
	if err != nil {
		return fragrantica.FragranceParams{}, fmt.Errorf("Failed parsing fragrance parameters '%s': %w", link.String, err)
	}
	// add name and brand which we got from url
	params.Name = name
	params.Brand = brand
	// try getting existing country from database (works on existing entries)
	// TODO: replace this with perfumers table lookup when that is populated
	country, err := frags.DB.GetFragranceCountry(context.Background(), id)
	if err != nil {
		return fragrantica.FragranceParams{}, fmt.Errorf("Failed getting fragrance country for ID %d: %w", id, err)
	}
	if country.Valid {
		params.Country = country.String
	}

	// add ID so sql finds the fragrance to update
	params.FragranticaID = id

	// update frag db
	err = frags.DB.UpdateFragrance(context.Background(), dbInput(params))
	if err != nil {
		return params, fmt.Errorf("Failed updating fragrance for ID %d: %w", id, err)
	}

	return params, nil
}

func parseURL(link string) (name string, brand string, err error) {
	u, err := url.Parse(link)
	if err != nil {
		return name, brand, err
	}

	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	// parts = ["perfume", "Brand", "Name-ID.html"]
	if len(parts) < 3 {
		return name, brand, fmt.Errorf("unexpected URL format: %s", link)
	}

	brand = strings.ToLower(parts[1])

	// Remove .html from last part
	nameWithID := parts[2]
	name = strings.TrimSuffix(nameWithID, ".html")
	// Remove trailing "-ID"
	lastDash := strings.LastIndex(name, "-")
	if lastDash != -1 {
		name = name[:lastDash]
	}
	name = strings.ToLower(name)

	return name, brand, nil
}

func dbInput(params fragrantica.FragranceParams) database.UpdateFragranceParams {
	db := database.UpdateFragranceParams{}
	db.FragranticaID = params.FragranticaID
	db.Name = nullString(params.Name)
	db.Brand = nullString(params.Brand)
	db.Country = nullString(params.Country)
	db.Gender = nullString(params.Gender)
	db.RatingValue = nullString(params.RatingValue)
	db.RatingCount = nullInt32(params.RatingCount)
	db.Year = nullInt32(params.Year)
	db.TopNotes = nullString(params.TopNotes)
	db.MiddleNotes = nullString(params.MiddleNotes)
	db.BaseNotes = nullString(params.BaseNotes)
	db.Perfumer1 = nullString(params.Perfumer1)
	db.Perfumer2 = nullString(params.Perfumer2)
	db.Accord1 = nullString(params.Accord1)
	db.Accord2 = nullString(params.Accord2)
	db.Accord3 = nullString(params.Accord3)
	db.Accord4 = nullString(params.Accord4)
	db.Accord5 = nullString(params.Accord5)
	db.Accord6 = nullString(params.Accord6)
	db.Accord7 = nullString(params.Accord7)
	db.Accord8 = nullString(params.Accord8)
	db.Accord9 = nullString(params.Accord9)
	db.Accord10 = nullString(params.Accord10)

	return db
}

func nullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func nullInt32(n int32) sql.NullInt32 {
	if n == 0 {
		return sql.NullInt32{}
	}
	return sql.NullInt32{
		Int32: n,
		Valid: true,
	}
}

// Delay of n-m seconds
func spamDelay(min, max int) {
	diff := 1
	if min < max {
		diff = max - min
	}
	jitter := rand.Intn(diff)
	time.Sleep(time.Duration(min+jitter) * 1000 * time.Millisecond)
}
