package fragrances

import (
	"testing"

	"github.com/slikasp/dbmanfrags/database"
	"github.com/slikasp/dbmanfrags/fragrantica"
)

func TestDbInputMissingValues(t *testing.T) {
	expectedDbParams := database.UpdateFragranceParams{}
	expectedDbParams.FragranticaID = 7702
	expectedDbParams.Name = nullString("jelly-belly-wild-blackberry-peach-cobbler")
	expectedDbParams.Brand = nullString("demeter-fragrance")
	expectedDbParams.Country = nullString("USA")
	expectedDbParams.Gender = nullString("women")
	expectedDbParams.RatingValue = nullString("3.85")
	expectedDbParams.RatingCount = nullInt32(26)
	expectedDbParams.TopNotes = nullString("amalfi lemon")
	expectedDbParams.MiddleNotes = nullString("blackberry")
	expectedDbParams.BaseNotes = nullString("peach")
	expectedDbParams.Perfumer1 = nullString("unknown")
	expectedDbParams.Accord1 = nullString("fruity")
	expectedDbParams.Accord2 = nullString("citrus")
	expectedDbParams.Accord3 = nullString("sweet")

	fragParams := fragrantica.FragranceParams{
		FragranticaID: 7702,
		Name:          "jelly-belly-wild-blackberry-peach-cobbler",
		Brand:         "demeter-fragrance",
		Country:       "USA",
		Gender:        "women",
		RatingValue:   "3.85",
		RatingCount:   26,
		TopNotes:      "amalfi lemon",
		MiddleNotes:   "blackberry",
		BaseNotes:     "peach",
		Perfumer1:     "unknown",
		Accord1:       "fruity",
		Accord2:       "citrus",
		Accord3:       "sweet",
	}

	dbParams := dbInput(fragParams)

	if expectedDbParams != dbParams {
		t.Errorf("Expected: %v\nGot %v", expectedDbParams, dbParams)
	}
}
