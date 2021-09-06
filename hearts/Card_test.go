package hearts

import (
	"fmt"
	"testing"
)

type testCard struct {
	card         Card
	expectedName string
}

var testCards []testCard = []testCard{
	{51, "Ace of Spades"},
	{0, "Two of Diamonds"},
	{12, "Ace of Diamonds"},
	{13, "Two of Clubs"},
	{14, "Three of Clubs"},
	{28, "Four of Hearts"},
}

func TestCardValueSuit(t *testing.T) {
	for _, c := range testCards {
		name := getCardName(c.card)

		if name != c.expectedName {
			t.Errorf("expected %s to have name %s", name, c.expectedName)
		}
	}
}

func TestCardCompare(t *testing.T) {
	var aceOfSpades Card = 51
	var jackOfSpades Card = 48
	var twoOfClubs Card = 13
	var threeOfClubs Card = 14
	var twoOfDiamonds Card = 0

	if aceOfSpades.Compare(jackOfSpades) <= 0 {
		t.Error("Ace of Spades should be greater than Two of Clubs")
	}

	if twoOfDiamonds.Compare(twoOfDiamonds) != 0 {
		t.Error("Two of Diamonds should euqal Two of Diamonds")
	}

	if twoOfClubs.Compare(threeOfClubs) >= 0 {
		t.Error("Two of Clubs should be less than Three of Clubs")
	}
}

func getCardName(c Card) string {
	return fmt.Sprintf("%s of %s", c.Value(), c.Suit())
}
