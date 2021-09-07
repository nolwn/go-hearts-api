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
	tryToPlayCards(t, game, 0, false, firstPlayerPassed...)

	// check that the passed cards have left the players hand
	checkCardsHaveMoved(t, game.Players[0].Hand, firstPlayerPassed)

	// check to see which players are able to play
	// expect all but the first player to be able to play
	checkActivePlayers(t, game, []int{1, 2, 3})

	// second player tries to pass 1 card they have once, and 1 card they have twice
	// expect an error to be returned
	player2Cards := secondPlayer.Hand
	tryToPlayCards(t, game, 1, true, player2Cards[7], player2Cards[3], player2Cards[3])

	// check to see which players ane able to play
	// expect all but the first player to be able to play
	checkActivePlayers(t, game, []int{1, 2, 3})

	// second player tries to pass a card they don't have
	// expect an error to be returned
	tryToPlayCards(t, game, 1, true, player2Cards[0], player2Cards[6], player1Cards[3])

	// third player plays three cards that they have
	// expect no error
	player3Cards := thirdPlayer.Hand
	thirdPlayerPassed := []Card{player3Cards[2], player3Cards[4], player3Cards[6]}
	tryToPlayCards(t, game, 2, false, thirdPlayerPassed...)

	// check that the passed cards have left the players hand
	checkCardsHaveMoved(t, game.Players[2].Hand, thirdPlayerPassed)

	// check to see which players ane able to play
	// expect the second and fourth players to be able to play
	checkActivePlayers(t, game, []int{1, 3})

	// second player tries to pass too few cards
	// expect an error (why is player 2 so stupid?!)
	tryToPlayCards(t, game, 1, true, player2Cards[2], player2Cards[3])

	// player four tries to pass too many cards
	// expect an error
	player4Cards := fourthPlayer.Hand
	tryToPlayCards(t, game, 3, true, player4Cards...)

	// second player finally passes 3 good cards
	// expect no error
	secondPlayerPassed := []Card{player2Cards[1], player2Cards[3], player2Cards[12]}
	tryToPlayCards(t, game, 1, false, secondPlayerPassed...)

	// check that the passed cards have left the players hand
	checkCardsHaveMoved(t, game.Players[1].Hand, secondPlayerPassed)

	// check to see which players ane able to play
	// expect only the fourth players to be able to play
	checkActivePlayers(t, game, []int{3})

	// fourth player passes three good cards
	// expect no error
	fourthPlayerPassed := []Card{player4Cards[7], player4Cards[4], player4Cards[10]}
	tryToPlayCards(t, game, 3, false, fourthPlayerPassed...)

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

func TestPlayPhasePlay(t *testing.T) {
	playerCards1 := []Card{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48}
	playerCards2 := []Card{1, 5, 9, 13, 17, 21, 25, 29, 33, 37, 41, 45, 49}
	playerCards3 := []Card{2, 6, 10, 14, 18, 22, 26, 30, 34, 38, 42, 46, 50}
	playerCards4 := []Card{3, 7, 11, 15, 19, 23, 27, 31, 35, 39, 43, 47, 51}

	hearts := New()

	hearts.Setup()
	hearts.phase = 1
	hearts.Players[0].Hand = playerCards1
	hearts.Players[1].Hand = playerCards2
	hearts.Players[2].Hand = playerCards3
	hearts.Players[3].Hand = playerCards4

	// expect player with two of clubs (player two) to be the only active player
	checkActivePlayers(t, &hearts, []int{1})

	// player 3 tries to play a card
	// expect an error
	tryToPlayCards(t, &hearts, 2, true, 6)

	// player 2 tries to play too many cards
	// expect an error
	tryToPlayCards(t, &hearts, 1, true, 5, 9)

	// player 2 tries to play too few cards
	// expect an error
	tryToPlayCards(t, &hearts, 1, true)

	// player 2 tries to play a card they don't hold
	// expect an error
	tryToPlayCards(t, &hearts, 1, true, 0)

	// player 2 tries to play a card that isn't the 2 of clubs
	// expect an error
	tryToPlayCards(t, &hearts, 1, true, 17)

	// player 2 plays the two of clubs
	// expect no error
	tryToPlayCards(t, &hearts, 1, false, 13)

	// expect the next player (player three) to be active
	checkActivePlayers(t, &hearts, []int{2})

	// player 2 tries to play another card
	// expect an error
	tryToPlayCards(t, &hearts, 1, true, 1)

	// player 3 tries to play an offsuit card
	// expect an error
	tryToPlayCards(t, &hearts, 2, true, 50)

	// player 3 plays a clubs
	// expect no error
	tryToPlayCards(t, &hearts, 2, false, 22)

	// expect the next player (player four) to be active
	checkActivePlayers(t, &hearts, []int{3})

	// player 4 plays a club
	// expect no error
	tryToPlayCards(t, &hearts, 3, false, 19)

	// expect the next player (player one) to be active
	checkActivePlayers(t, &hearts, []int{0})

	// player 1 plays a club
	// expect no error
	tryToPlayCards(t, &hearts, 0, false, 20)

	// expect the player who took the trick (player three) to be active
	checkActivePlayers(t, &hearts, []int{2})
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

func tryToPlayCards(
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
