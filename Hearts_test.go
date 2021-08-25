package hearts

import (
	"fmt"
	"testing"

	"github.com/nolwn/go-hearts/hearts"
)

type testCard struct {
	card         hearts.Card
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

func TestHeartsSetup(t *testing.T) {
	round := hearts.New()
	err := round.Setup()
	cards := make(map[string]bool)

	if err != nil {
		t.Errorf("setup should not have returned error: %s", err)
	}

	if len(round.Players) != 4 {
		t.Errorf("not all the players have been set up")
	}

	for i, p := range round.Players {
		if len(p.Hand) != 13 {
			t.Errorf("player %d is missing a hand", i)
		}

		for _, c := range p.Hand {
			card := fmt.Sprintf("%s of %s", c.Value(), c.Suit())
			if cards[card] == true {
				t.Errorf("%s was dealt more than once", card)
			}

			cards[card] = true
		}
	}

	if round.Finished() {
		t.Error("the game should not be finished")
	}

	if round.Phase() != hearts.PhasePass {
		t.Error(("setup should put the game in the pass phase"))
	}

	if round.Round() != 1 {
		t.Errorf("expected round %d to be round 1", round.Round())
	}
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
	var aceOfSpades hearts.Card = 51
	var jackOfSpades hearts.Card = 48
	var twoOfClubs hearts.Card = 13
	var threeOfClubs hearts.Card = 14
	var twoOfDiamonds hearts.Card = 0

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

func getCardName(c hearts.Card) string {
	return fmt.Sprintf("%s of %s", c.Value(), c.Suit())
}
