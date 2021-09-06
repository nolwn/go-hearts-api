package hearts

import (
	"fmt"
	"testing"
)

func TestHeartsPlay(t *testing.T) {
	round := New()
	err := round.Setup()

	if err != nil {
		t.Errorf("setup should not have returned error: %s", err)
	}
}

func TestHeartsSetup(t *testing.T) {
	game := setupGame(t)
	cards := make(map[string]bool)

	if len(game.Players) != 4 {
		t.Errorf("not all the players have been set up")
	}

	for i, p := range game.Players {
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

	if game.Finished() {
		t.Error("the game should not be finished")
	}

	if game.Phase() != PhasePass {
		t.Error(("setup should put the game in the pass phase"))
	}

	if game.Round() != 1 {
		t.Errorf("expected round %d to be round 1", game.Round())
	}
}

func TestHeartsPassTurns(t *testing.T) {
	game := setupGame(t)

	firstPlayer := game.Players[0]
	secondPlayer := game.Players[1]
	thirdPlayer := game.Players[2]
	fourthPlayer := game.Players[3]

	// the first phase (0) is the passing phase
	if game.Phase() != 0 {
		t.Errorf("expected phase 0 but the game is in phase %d", game.Phase())
	}

	// check to see which players are able to play
	// expect all players to be able to play
	checkActivePlayers(t, game, []int{0, 1, 2, 3})

	// first player passes three valid cards
	// expect no error
	player1Cards := firstPlayer.Hand
	firstPlayerPassed := []Card{player1Cards[0], player1Cards[1], player1Cards[5]}
	tryToPlassCards(t, game, 0, false, firstPlayerPassed...)

	// check that the passed cards have left the players hand
	checkCardsHaveMoved(t, game.Players[0].Hand, firstPlayerPassed)

	// check to see which players are able to play
	// expect all but the first player to be able to play
	checkActivePlayers(t, game, []int{1, 2, 3})

	// second player tries to pass 1 card they have once, and 1 card they have twice
	// expect an error to be returned
	player2Cards := secondPlayer.Hand
	tryToPlassCards(t, game, 1, true, player2Cards[7], player2Cards[3], player2Cards[3])

	// check to see which players ane able to play
	// expect all but the first player to be able to play
	checkActivePlayers(t, game, []int{1, 2, 3})

	// second player tries to pass a card they don't have
	// expect an error to be returned
	tryToPlassCards(t, game, 1, true, player2Cards[0], player2Cards[6], player1Cards[3])

	// third player plays three cards that they have
	// expect no error
	player3Cards := thirdPlayer.Hand
	thirdPlayerPassed := []Card{player3Cards[2], player3Cards[4], player3Cards[6]}
	tryToPlassCards(t, game, 2, false, thirdPlayerPassed...)

	// check that the passed cards have left the players hand
	checkCardsHaveMoved(t, game.Players[2].Hand, thirdPlayerPassed)

	// check to see which players ane able to play
	// expect the second and fourth players to be able to play
	checkActivePlayers(t, game, []int{1, 3})

	// second player tries to pass too few cards
	// expect an error (why is player 2 so stupid?!)
	tryToPlassCards(t, game, 1, true, player2Cards[2], player2Cards[3])

	// player four tries to pass too many cards
	// expect an error
	player4Cards := fourthPlayer.Hand
	tryToPlassCards(t, game, 3, true, player4Cards...)

	// second player finally passes 3 good cards
	// expect no error
	secondPlayerPassed := []Card{player2Cards[1], player2Cards[3], player2Cards[12]}
	tryToPlassCards(t, game, 1, false, secondPlayerPassed...)

	// check that the passed cards have left the players hand
	checkCardsHaveMoved(t, game.Players[1].Hand, secondPlayerPassed)

	// check to see which players ane able to play
	// expect only the fourth players to be able to play
	checkActivePlayers(t, game, []int{3})

	// fourth player passes three good cards
	// expect no error
	fourthPlayerPassed := []Card{player4Cards[7], player4Cards[4], player4Cards[10]}
	tryToPlassCards(t, game, 3, false, fourthPlayerPassed...)

	// check that the passed cards have left the players hand
	checkCardsHaveMoved(t, game.Players[3].Hand, fourthPlayerPassed)

	// once passing is finished the game should be in the play phase (1)
	if game.Phase() != 1 {
		t.Errorf("expected phase 1 but the game is in phase %d", game.Phase())
	}

	// check to see which players ane able to play
	// expect no players to be able to play
	twoOfClubs := findTwoOfClubs(game)
	checkActivePlayers(t, game, []int{twoOfClubs})

	checkCardsReceived(t, game.Players[0].Hand, secondPlayerPassed)
	checkCardsReceived(t, game.Players[1].Hand, thirdPlayerPassed)
	checkCardsReceived(t, game.Players[2].Hand, fourthPlayer.Taken)
	checkCardsReceived(t, game.Players[3].Hand, firstPlayerPassed)
}

func setupGame(t *testing.T) *Hearts {
	round := New()
	err := round.Setup()

	if err != nil {
		t.Errorf("setup should not have returned error: %s", err)
	}

	return &round
}

// returns true if the hands match
func compareSlices(first []int, second []int) bool {
	if len(first) != len(second) {
		return false
	}

	for i := 0; i < len(first); i++ {
		if first[i] != second[i] {
			return false
		}
	}

	return true
}

func checkCardsHaveMoved(t *testing.T, hand []Card, cards []Card) {
	handMap := map[int]bool{}

	for _, card := range hand {
		handMap[int(card)] = true
	}

	for _, card := range cards {
		if handMap[int(card)] {
			t.Errorf("expected %d not to be in the hand", card)
		}
	}
}

func checkActivePlayers(
	t *testing.T,
	game *Hearts,
	expectedPlayers []int,
) {
	players := game.PlayersTurn()

	if !compareSlices(players, expectedPlayers) {
		t.Errorf("expected player %d but received %d", expectedPlayers, players)
	}
}

func checkCardsReceived(t *testing.T, hand []Card, cards []Card) {
	handMap := make(map[Card]bool)

	for _, c := range hand {
		handMap[c] = true
	}

	for _, c := range cards {
		if !handMap[c] {
			t.Errorf("card %d  get correctly passed", c)
		}
	}
}

func tryToPlassCards(
	t *testing.T,
	game *Hearts,
	player int,
	shouldFail bool,
	cards ...Card,
) {
	err := game.Play(player, cards...)

	if shouldFail {
		if err == nil {
			t.Error("expected an error but did not receive one")
		}

	} else {
		if err != nil {
			t.Errorf("expected no error but received: %s", err)
		}
	}
}

func findTwoOfClubs(game *Hearts) int {
	twoOfClubs := 13

	for i, p := range game.Players {
		for _, c := range p.Hand {
			if c == Card(twoOfClubs) {
				return i
			}
		}
	}

	return -1
}
