package hearts

import (
	"fmt"
	"testing"
)

const (
	handFull = iota
	handSmall
	handFinal
	handAllHearts
	handAllHeartsSmall
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

	// hands should be sorted at the begining of the game
	sorted := checkHandsAreSorted(game)

	if !sorted {
		t.Error("expected hands to be sorted after passing, but they were not")
	}

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
	play(t, game, 0, false, firstPlayerPassed...)

	// check that the passed cards have left the players hand
	checkCardsHaveMoved(t, game.Players[0].Hand, firstPlayerPassed)

	// check to see which players are able to play
	// expect all but the first player to be able to play
	checkActivePlayers(t, game, []int{1, 2, 3})

	// second player tries to pass 1 card they have once, and 1 card they have twice
	// expect an error to be returned
	player2Cards := secondPlayer.Hand
	play(t, game, 1, true, player2Cards[7], player2Cards[3], player2Cards[3])

	// check to see which players ane able to play
	// expect all but the first player to be able to play
	checkActivePlayers(t, game, []int{1, 2, 3})

	// second player tries to pass a card they don't have
	// expect an error to be returned
	play(t, game, 1, true, player2Cards[0], player2Cards[6], player1Cards[3])

	// third player plays three cards that they have
	// expect no error
	player3Cards := thirdPlayer.Hand
	thirdPlayerPassed := []Card{player3Cards[2], player3Cards[4], player3Cards[6]}
	play(t, game, 2, false, thirdPlayerPassed...)

	// check that the passed cards have left the players hand
	checkCardsHaveMoved(t, game.Players[2].Hand, thirdPlayerPassed)

	// check to see which players ane able to play
	// expect the second and fourth players to be able to play
	checkActivePlayers(t, game, []int{1, 3})

	// second player tries to pass too few cards
	// expect an error (why is player 2 so stupid?!)
	play(t, game, 1, true, player2Cards[2], player2Cards[3])

	// player four tries to pass too many cards
	// expect an error
	player4Cards := fourthPlayer.Hand
	play(t, game, 3, true, player4Cards...)

	// second player finally passes 3 good cards
	// expect no error
	secondPlayerPassed := []Card{player2Cards[1], player2Cards[3], player2Cards[12]}
	play(t, game, 1, false, secondPlayerPassed...)

	// check that the passed cards have left the players hand
	checkCardsHaveMoved(t, game.Players[1].Hand, secondPlayerPassed)

	// check to see which players ane able to play
	// expect only the fourth players to be able to play
	checkActivePlayers(t, game, []int{3})

	// fourth player passes three good cards
	// expect no error
	fourthPlayerPassed := []Card{player4Cards[7], player4Cards[4], player4Cards[10]}
	play(t, game, 3, false, fourthPlayerPassed...)

	// check that the passed cards have left the players hand
	checkCardsHaveMoved(t, game.Players[3].Hand, fourthPlayerPassed)

	// hands should remain sorted at the end of the passing phase
	sorted = checkHandsAreSorted(game)

	if !sorted {
		t.Error("expected hands to be sorted after passing, but they were not")
	}

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

	if hearts.trick != 1 {
		t.Errorf("should start on trick 1 but it's trick %d", hearts.trick)
	}

	// expect player with two of clubs (player two) to be the only active player
	checkActivePlayers(t, &hearts, []int{1})

	// player 3 tries to play a card
	// expect an error
	play(t, &hearts, 2, true, 6)

	// player 2 tries to play too many cards
	// expect an error
	play(t, &hearts, 1, true, 5, 9)

	// player 2 tries to play too few cards
	// expect an error
	play(t, &hearts, 1, true)

	// player 2 tries to play a card they don't hold
	// expect an error
	play(t, &hearts, 1, true, 0)

	// player 2 tries to play a card that isn't the 2 of clubs
	// expect an error
	play(t, &hearts, 1, true, 17)

	// player 2 plays the two of clubs
	// expect no error
	play(t, &hearts, 1, false, 13)

	// expect the next player (player one) to be active
	checkActivePlayers(t, &hearts, []int{0})

	// player 2 tries to play another card
	// expect an error
	play(t, &hearts, 1, true, 1)

	// player 1 tries to play an offsuit card
	// expect an error
	play(t, &hearts, 0, true, 50)

	// player 1 plays a clubs
	// expect no error
	play(t, &hearts, 0, false, 24)

	// expect the next player (player four) to be active
	checkActivePlayers(t, &hearts, []int{3})

	// player 4 plays a club
	// expect no error
	play(t, &hearts, 3, false, 19)

	// expect the next player (player three) to be active
	checkActivePlayers(t, &hearts, []int{2})

	// player 3 plays a club
	// expect no error
	play(t, &hearts, 2, false, 22)

	// expect the player who took the trick (player one) to be active
	checkActivePlayers(t, &hearts, []int{0})

	// expect trick number to increment
	if hearts.trick != 2 {
		t.Errorf("expected trick to increment by 1 but it's trick %d", hearts.trick)
	}
}

func TestCardPassDirection(t *testing.T) {
	// round 1 should pass left
	hearts := setupCannedHands(handFull)
	hearts.round = 1

	cards := passCards(&hearts) // pass cards

	has := checkPassedCards(hearts, cards, getLeftIndex) // make sure cards went left

	if !has {
		t.Error("expected cards to be passed left on round 1, but they were not")
	}

	// round 2 should pass right
	hearts = setupCannedHands(handFull)
	hearts.round = 2

	cards = passCards(&hearts)

	has = checkPassedCards(hearts, cards, getRightIndex)

	if !has {
		t.Error("expected cards to be passed right on round 2, but they were not")
	}

	// round 3 should pass across
	hearts = setupCannedHands(handFull)
	hearts.round = 3

	cards = passCards(&hearts)

	has = checkPassedCards(hearts, cards, getAcrossIndex)

	if !has {
		t.Error("expected cards to be passed across on round 3, but they were not")
	}

	// round 4 should not have a passing phase
	hearts = setupCannedHands(handFull)

	// we start on round three so we can advance to round 4 and make sure it starts on
	// the play phase. If we didn't do it this way, we could force an incorrect game
	// state.
	hearts.round = 3
	hearts.phase = PhasePlay
	hearts.phaseEnd = true
	hearts.NextPhase()

	if hearts.phase != PhasePlay {
		t.Error("expected round 4 to skip the passing phase, but it did not")
	}

	// round 5 should restart the pattern
	hearts = setupCannedHands(handFull)
	hearts.round = 5

	cards = passCards(&hearts) // pass cards

	has = checkPassedCards(hearts, cards, getLeftIndex) // make sure cards went left

	if !has {
		t.Error("expected cards to be passed left on round 5, but they were not")
	}
}

func TestPointValuesHearts(t *testing.T) {
	h := setupCannedHands(handSmall)
	h.phase = PhasePlay
	h.lastTrick = 2
	h.trick = 7

	// player 3 plays a club
	// player 2 plays their highest club
	// player 1 sluffs a heart
	// player 4 plays a club
	play(t, &h, 2, false, card(h.Players[2].Hand, SuitClubs))
	play(t, &h, 1, false, card(h.Players[1].Hand, SuitDiamonds))
	play(t, &h, 0, false, card(h.Players[0].Hand, SuitHearts))
	play(t, &h, 3, false, card(h.Players[3].Hand, SuitClubs))

	// The trick has now been taken. Figure out who took it and count up their points.
	took := h.lastTrick
	takenScore := h.Players[took].roundScore

	// the number of points taken should be 1 for the one heart that was sluffed by player 1
	if takenScore != 1 {
		t.Errorf(
			"player %d should have taken 1 point worth hearts, but instead they took %d",
			took,
			takenScore,
		)
	}

	// setup a new game
	h = setupCannedHands(handSmall)
	h.phase = PhasePlay
	h.lastTrick = 2
	h.trick = 7

	// player 3 plays their highest club
	// player 2 leads with their highest club
	// player 1 sluffs a heart
	// player 4 plays their highest club
	play(t, &h, 2, false, card(h.Players[2].Hand, SuitClubs))
	play(t, &h, 1, false, card(h.Players[1].Hand, SuitHearts))
	play(t, &h, 0, false, card(h.Players[0].Hand, SuitHearts))
	play(t, &h, 3, false, card(h.Players[3].Hand, SuitClubs))

	took = h.lastTrick
	takenScore = h.Players[took].roundScore

	if takenScore != 2 {
		t.Errorf(
			"player %d should have taken 2 point worth hearts, but instead they took %d",
			took,
			takenScore,
		)
	}
}

func TestPointValuesJamoke(t *testing.T) {
	h := setupCannedHands(handSmall)
	h.phase = PhasePlay
	h.lastTrick = 2
	h.trick = 7

	// player 3 plays a club
	// player 2 players their highest club
	// player 1 sluffs a heart
	// player 4 plays a club
	play(t, &h, 2, false, card(h.Players[2].Hand, SuitClubs))
	play(t, &h, 1, false, CardJamoke)
	play(t, &h, 0, false, card(h.Players[0].Hand, SuitSpades))
	play(t, &h, 3, false, card(h.Players[3].Hand, SuitClubs))

	took := h.lastTrick
	takenScore := h.Players[took].roundScore

	if takenScore != 13 {
		t.Errorf(
			"player %d should have taken 13 point worth hearts, but instead they took %d",
			took,
			takenScore,
		)
	}

	h = setupCannedHands(handSmall)
	h.phase = PhasePlay
	h.lastTrick = 2
	h.trick = 7

	// player 3 plays a club
	// player 2 plays The Jamoke
	// player 1 sluffs a heart
	// player 4 plays a club
	play(t, &h, 2, false, card(h.Players[2].Hand, SuitClubs))
	play(t, &h, 1, false, CardJamoke)
	play(t, &h, 0, false, card(h.Players[0].Hand, SuitHearts))
	play(t, &h, 3, false, card(h.Players[3].Hand, SuitClubs))

	took = h.lastTrick
	takenScore = h.Players[took].roundScore

	if takenScore != 14 {
		t.Errorf(
			"player %d should have taken 14 point worth hearts, but instead they took %d",
			took,
			takenScore,
		)
	}
}

func TestNoPointsOnFirstTrick(t *testing.T) {
	h := setupCannedHands(handSmall)
	h.phase = PhasePlay
	h.lastTrick = PlayerFour
	h.trick = 1 // Although the two of clubs is gone, we are pretending this is round 1

	// player 4 leads with clubs
	// player 3 follows suit
	// player 2 has not clubs and tries to play a heart which should fail.
	play(t, &h, PlayerFour, false, card(h.Players[PlayerFour].Hand, SuitClubs))
	play(t, &h, PlayerThree, false, card(h.Players[PlayerThree].Hand, SuitClubs))
	play(t, &h, PlayerTwo, true, card(h.Players[PlayerTwo].Hand, SuitHearts))

	// however, if it weren't round 1...
	h.trick = 2

	// then it should succeed
	play(t, &h, PlayerTwo, false, card(h.Players[PlayerTwo].Hand, SuitHearts))

	// An exception is made if the player ONLY has Hearts (a rare but possible situation)
	h = setupCannedHands(handAllHearts)
	h.phase = PhasePlay
	h.lastTrick = PlayerThree
	h.trick = 1

	// player 4 leads with clubs
	// player 3 follows suit
	// player 2 has not clubs and tries to play a heart which should fail.
	play(t, &h, PlayerThree, false, CardTwoOfClubs)
	play(t, &h, PlayerTwo, false, card(h.Players[PlayerTwo].Hand, SuitClubs))
	play(t, &h, PlayerOne, false, card(h.Players[PlayerOne].Hand, SuitHearts))
}

func TestNoLeadingHeartsBeforeBroken(t *testing.T) {
	h := setupCannedHands(handSmall)
	h.phase = PhasePlay
	h.lastTrick = PlayerFour
	h.trick = 7
	h.brokenHearted = false

	// player 4 leads with hearts
	play(t, &h, PlayerFour, true, card(h.Players[PlayerFour].Hand, SuitHearts))

	// however, if hearts were broken...
	h.brokenHearted = true

	play(t, &h, PlayerFour, false, card(h.Players[PlayerFour].Hand, SuitHearts))

	h = setupCannedHands(handAllHeartsSmall)
	h.phase = PhasePlay
	h.lastTrick = PlayerOne
	h.trick = 7
	h.brokenHearted = false

	// An exception is mad if the palyer ONLY has Hearts (a not so rare situation)
	play(t, &h, PlayerOne, false, card(h.Players[PlayerOne].Hand, SuitHearts))
}

func TestRoundEnd(t *testing.T) {
	h := setupCannedHands(handFinal)
	h.phase = PhasePlay
	h.lastTrick = 1
	h.trick = 12

	startingScore := h.Score()

	playerOneCard := card(h.Players[0].Hand, SuitHearts)
	playerTwoCard := card(h.Players[1].Hand, SuitDiamonds)
	playerThreeCard := card(h.Players[2].Hand, SuitHearts)
	playerFourCard := card(h.Players[3].Hand, SuitDiamonds)

	// player 2 players their last diamond
	// player 1 sluffs a heart
	// player 4 plays their last diamond
	// player 3 sluffs a heart
	play(t, &h, 1, false, playerTwoCard)
	play(t, &h, 0, false, playerOneCard)
	play(t, &h, 3, false, playerFourCard)
	play(t, &h, 2, false, playerThreeCard)

	var took int

	if playerTwoCard > playerFourCard {
		took = PlayerTwo
	} else {
		took = PlayerFour
	}

	if h.Phase() != PhasePass {
		t.Error("the round should have ended, but we are still in the pass phase")
	}

	finalScore := h.Score()

	for _, player := range h.Players {
		if len(player.Hand) != 13 {
			t.Error("expected cards to be dealt, but they were not")
		}
	}

	for p, start := range startingScore {
		if p == took {
			if start-finalScore[p] != 2 {
				t.Errorf(
					"player %d should have lost 2 points, but they lost %d",
					p,
					start-finalScore[p],
				)
			}
		}
	}
}

func TestMoonShot(t *testing.T) {
	h := setupCannedHands(handFinal)
	h.phase = PhasePlay
	h.lastTrick = 1
	h.trick = 12

	startingScore := h.Score()

	playerOneCard := card(h.Players[0].Hand, SuitHearts)
	playerTwoCard := card(h.Players[1].Hand, SuitDiamonds)
	playerThreeCard := card(h.Players[2].Hand, SuitHearts)
	playerFourCard := card(h.Players[3].Hand, SuitDiamonds)

	var takes int

	if playerTwoCard > playerFourCard {
		takes = PlayerTwo
	} else {
		takes = PlayerFour
	}

	h.Players[takes].roundScore = 24 // 2 hearts from a moonshot

	// player 2 players their last diamond
	// player 1 sluffs a heart
	// player 4 plays their last diamond
	// player 3 sluffs a heart
	play(t, &h, 1, false, playerTwoCard)
	play(t, &h, 0, false, playerOneCard)
	play(t, &h, 3, false, playerFourCard)
	play(t, &h, 2, false, playerThreeCard)

	if h.Phase() != PhasePass {
		t.Error("the round should have ended, but we are still in the pass phase")
	}

	finalScore := h.Score()

	for p, start := range startingScore {
		if p == takes {
			if start-finalScore[p] != 0 {
				t.Errorf(
					"player %d should have lost 0 points, but they lost %d",
					p,
					start-finalScore[p],
				)
			}
		} else {
			if start-finalScore[p] != 26 {
				t.Errorf(
					"expected other players to take 26 points, but player %d took %d",
					p,
					start-finalScore[p],
				)
			}
		}
	}
}

func TestGameEnd(t *testing.T) {
	h := setupCannedHands(handFinal)
	h.phase = PhasePlay
	h.lastTrick = PlayerTwo
	h.trick = 12

	for p := range h.Players {
		player := &h.Players[p]
		player.gameScore = 1
	}

	h.Players[PlayerThree].gameScore = 100

	// player 2 players their last diamond
	// player 1 sluffs a heart
	// player 4 plays their last diamond
	// player 3 sluffs a heart
	play(t, &h, PlayerTwo, false, card(h.Players[PlayerTwo].Hand, SuitDiamonds))
	play(t, &h, PlayerOne, false, card(h.Players[PlayerOne].Hand, SuitHearts))
	play(t, &h, PlayerFour, false, card(h.Players[PlayerFour].Hand, SuitDiamonds))
	play(t, &h, PlayerThree, false, card(h.Players[PlayerThree].Hand, SuitHearts))

	if !h.Finished() {
		t.Error("the game should be finished but it is not")
	}

	if !compareSlices(h.Winner(), []int{PlayerThree}) {
		t.Errorf("expected %d to win, but recieved %v instead", PlayerThree, h.Winner())
	}
}

func hasCards(hand []Card, cards ...Card) bool {
	handMap := map[Card]bool{}

	for _, c := range hand {
		handMap[c] = true
	}

	for _, c := range cards {
		if !handMap[c] {
			return false
		}
	}

	return true
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

// checkHandsAreSorted takes a game state and checks if all the players' hands are sorted
// in ascending order. It returns true if they are, otherwise it returns false.
func checkHandsAreSorted(hearts *Hearts) bool {
	for _, p := range hearts.Players {
		hand := p.Hand

		// we need to skip the last index, since we will be comparing each index with
		// the next index
		for i := 0; i < len(hand)-1; i++ {
			if hand[i] > hand[i+1] {
				return false
			}
		}
	}

	return true
}

func checkPassedCards(hearts Hearts, cards [][]Card, getIdx func(int) int) bool {
	for i, c := range cards {
		passedTo := getIdx(i)
		if !hasCards(hearts.Players[passedTo].Hand, c...) {
			return false
		}
	}

	return true
}

func getAcrossIndex(i int) int {
	return getLeftIndex(getLeftIndex(i)) // pass 2 to the left
}

func getLeftIndex(i int) int {
	if i == 0 {
		return 3
	} else {
		return i - 1
	}
}

func getRightIndex(i int) int {
	if i == 3 {
		return 0
	} else {
		return i + 1
	}
}

func passCards(hearts *Hearts) [][]Card {
	cards := [][]Card{
		{0, 4, 12},
		{1, 5, 9},
		{2, 6, 10},
		{3, 7, 11},
	}

	hearts.Play(0, cards[0]...)
	hearts.Play(1, cards[1]...)
	hearts.Play(2, cards[2]...)
	hearts.Play(3, cards[3]...)

	return cards
}

// card returns the highest card of a given suit in the given hand. If no cards
// in the given suit are found, returns -1.
func card(hand []Card, suit string) Card {
	for i := len(hand) - 1; i >= 0; i-- {
		if hand[i].Suit() == suit {
			return hand[i]
		}
	}

	return -1
}

func setupCannedHands(hand int) Hearts {
	// As a helpful reminder:
	// Diamonds 0–12
	// Clubs 	13–25
	// Hearts 	26–38
	// Spades 	39–51

	hearts := New()
	hearts.Setup()

	var playerCards1 []Card
	var playerCards2 []Card
	var playerCards3 []Card
	var playerCards4 []Card

	switch hand {

	// Everyone has full hands and the cards are perfectly evenly divided. Player two
	// has the Two of Clubs.
	case handFull:
		playerCards1 = []Card{0, 4, 8, 12, 16, 20, 24, 28, 32, 36, 40, 44, 48}
		playerCards2 = []Card{1, 5, 9, 13, 17, 21, 25, 29, 33, 37, 41, 45, 49}
		playerCards3 = []Card{2, 6, 10, 14, 18, 22, 26, 30, 34, 38, 42, 46, 50}
		playerCards4 = []Card{3, 7, 11, 15, 19, 23, 27, 31, 35, 39, 43, 47, 51}

	// Player one and two are both fully out of clubs. Both have hearts they can sluff, and player
	// two has the queen of spades.
	case handSmall:
		playerCards1 = []Card{4, 8, 12, 28, 32, 36, 40, 44, 48}
		playerCards2 = []Card{1, 5, 9, 29, 33, 37, 41, 45, 49}
		playerCards3 = []Card{2, 6, 10, 18, 22, 30, 34, 42, 50}
		playerCards4 = []Card{3, 11, 19, 27, 35, 39, 43, 47, 51}

	// Player one and three have hearts. Player two and four have diamonds.
	case handFinal:
		playerCards1 = []Card{28}
		playerCards2 = []Card{1}
		playerCards3 = []Card{30}
		playerCards4 = []Card{3}

	// Everyone has full hands, but player one has all the hearts. Player three has the
	// Two of Clubs.
	case handAllHearts:
		playerCards1 = []Card{26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38}
		playerCards2 = []Card{0, 3, 6, 9, 12, 15, 18, 21, 24, 39, 42, 45, 48}
		playerCards3 = []Card{1, 4, 7, 10, 13, 16, 19, 22, 25, 40, 43, 46, 49}
		playerCards4 = []Card{2, 5, 8, 11, 14, 17, 20, 23, 41, 44, 47, 50, 51}

	// Everyone has some cards, and the two of clubs is gone, but player one has all the
	// hearts.
	case handAllHeartsSmall:
		playerCards1 = []Card{26, 27, 28, 29, 30, 31}
		playerCards2 = []Card{0, 3, 6, 9, 12, 15}
		playerCards3 = []Card{1, 4, 7, 10, 16, 19}
		playerCards4 = []Card{2, 5, 8, 11, 14, 17}
	}

	hearts.Players[0].Hand = playerCards1
	hearts.Players[1].Hand = playerCards2
	hearts.Players[2].Hand = playerCards3
	hearts.Players[3].Hand = playerCards4

	return hearts
}

func play(
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

// findTwoOfClubs returns the index of the player who has the 2 of clubs. If no player
// is holding it, it returns -1.
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
