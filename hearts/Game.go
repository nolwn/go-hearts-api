package hearts

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// The Two of Clubs is represented by the integer 13
const CardTwoOfClubs Card = 13
const CardJamoke Card = 49

const (
	Nobody = iota - 1
	PlayerOne
	PlayerTwo
	PlayerThree
	PlayerFour
)

// Finished returns true if the game is over and can no longer be played. A game of
// Hearts is considered finished when a player's score has crossed a certain threshhold,
// generally 100 points.
func (h *Hearts) Finished() bool {
	return h.finished
}

// Play in Hearts means one of two things depending on the phase. In the pass
// phase, players pick three cards to pass. In the play phase, players pick one card
// to play into trick.
//
// An error will be returned if it isn't the players turn to play.
func (h *Hearts) Play(player int, cards ...Card) error {
	if h.finished {
		return errors.New("the game is finished")
	}

	playing := h.PlayersTurn()
	canPlay := false

	for _, p := range playing {
		if p == player {
			canPlay = true
		}
	}

	if !canPlay {
		return fmt.Errorf("it is not player %d's turn", player)
	}

	if h.Phase() == PhasePlay {
		return h.playPhase(player, cards...)
	} else {
		return h.passPhase(player, cards...)
	}
}

// Player returns the index of the players who are allowed to take a turn. During the
// pass phase all players who have not yet passed cards are able to play.
//
// During the play phase, the player who can play is either the last player who took a
// trick or, if it's the first round, the player with the two of clubs.
func (h *Hearts) PlayersTurn() []int {
	if h.phase == PhasePass {
		return h.currentlyPassing()
	} else {
		return h.currentlyPlaying()
	}
}

// Setup sets up a new Hearts round. It deals out 13 cards randomly to each player. It
// also clears our each player's Taken slice.
func (h *Hearts) Setup() error {
	cards := 0

	for _, p := range h.Players {
		cards += len(p.Hand)
	}

	if cards != 0 {
		return errors.New("not all cards have been played")
	}

	h.clearTaken()
	h.deal()
	h.lastTaken = -1
	h.finished = false

	for _, p := range h.Players {
		sort(p.Hand, 0, len(p.Hand)-1)
	}

	return nil
}

// Winner returns the winner of the game. In Hearts, the winner is the player who has
// the smallest score at the end of the game. Scores, in this version, are tracked
// started from the score threshhold that determines when the game is finished, and move
// toward zero as the game progresses. Because of this, the winner is actually the player
// with the highest score.
//
// If there is a tie between players, then all players with the winning score are
// returned. If the game is not finished, Winner returns an empty array.
func (h *Hearts) Winner() (winners []int) {
	if !h.finished {
		return
	}

	best := h.Players[PlayerOne].gameScore

	for p, player := range h.Players {
		if player.gameScore > best { // new best score
			best = player.gameScore
			winners = []int{p}
		} else if player.gameScore == best { // tie for the best so far
			winners = append(winners, p)
		}
	}

	return
}

// clearTaken clears out any tricks taken by each of the players
func (h *Hearts) clearTaken() {
	for _, p := range h.Players {
		p.Taken = make([]Card, 0, 13)
	}
}

func (h *Hearts) deal() {
	var n Card = 0 // cards in a deck, starting with 0
	i := PlayerOne // player index
	deck := make([]Card, 0, 52)

	for n < 52 {
		deck = append(deck, n)
		n++
	}

	n-- // the loop above stopped once n was 1 over the max. Need to bump it back down 1

	seed := time.Now().UnixNano()
	r := rand.New(rand.NewSource(seed))

	for n > 0 {
		ri := r.Intn(int(n))
		card := deck[ri]

		// overwrite the selected card with the one that will be removed when we shorten
		// the deck
		deck[ri] = deck[n]
		n--

		h.Players[i].Hand = append(h.Players[i].Hand, card)

		if i == PlayerFour {
			i = PlayerOne
		} else {
			i++
		}
	}

	// the loop ends one card early because we can't generate a random number between 0
	// and 0. We just need to scoop up that last card.
	h.Players[i].Hand = append(h.Players[i].Hand, deck[0])
}

// passAcross returns the target across the table
func (h *Hearts) passAcross(player int, cards []Card) int {
	switch player {
	case PlayerOne:
		return PlayerThree
	case PlayerTwo:
		return PlayerFour
	case PlayerThree:
		return PlayerOne
	default: // case PlayerFour:
		return PlayerTwo
	}
}

// passLeft returns the target to the players left (toward begining of array)
func (h *Hearts) passLeft(player int, cards []Card) int {
	if player == PlayerOne {
		return PlayerFour
	} else {
		return player - 1
	}
}

// passPhase contains the logic for playing a card during the pass phase.
func (h *Hearts) passPhase(player int, cards ...Card) error {
	pass := h.Round() % 4
	playerHand := &h.Players[player].Hand
	var target int

	if len(cards) != 3 {
		return errors.New("player must pass exactly 3 cards")
	}

	for i, c := range cards {
		for j := i + 1; j < len(cards); j++ {
			if c == cards[j] {
				return errors.New("player must play 3 different cards")
			}
		}
	}

	if !hasCard(*playerHand, cards...) {
		return errors.New("player must have the cards to pass them")
	}

	switch pass {
	case 0:
		return errors.New("player cannot pass on the hold round")
	case 1:
		target = h.passLeft(player, cards)
	case 2:
		target = h.passRight(player, cards)
	case 3:
		target = h.passAcross(player, cards)
	}

	*playerHand = removeCard(*playerHand, cards...)
	h.Players[player].hasPassed = true
	h.Players[target].Receiving = cards

	playing := h.PlayersTurn()

	if len(playing) == 0 {
		h.phaseEnd = true
		err := h.NextPhase()

		if err != nil {
			return errors.New("unable to advance the game state")
		}
	}

	return nil
}

// currentlyPassing returns the players who have not yet picked cards to pass
func (h *Hearts) currentlyPassing() (players []int) {
	for i, p := range h.Players {
		if !p.hasPassed {
			players = append(players, i)
		}
	}

	return
}

// passRight returns the target to the players right (toward end of array)
func (h *Hearts) passRight(player int, cards []Card) int {
	if player == PlayerFour {
		return PlayerOne
	} else {
		return player + 1
	}
}

// playPhase contains the logic for playing a card during the play phase.
func (h *Hearts) playPhase(p int, cards ...Card) error {
	if len(cards) != 1 {
		return errors.New("player must play exactly one card")
	}

	hand := &h.Players[p].Hand
	played := &h.Players[p].Played
	keepPlaying := false

	// check that the player has the card
	if !hasCard(*hand, cards[0]) {
		return fmt.Errorf("player %d does not have %d", p, cards)
	}

	// if the player has the two of clubs, they MUST play it
	if hasTwoOfClubs(*hand) {
		if cards[0] != CardTwoOfClubs {
			return fmt.Errorf(
				"player has the two of clubs, but is trying to play %d",
				cards[0],
			)
		}
	}

	// if a suit was led, and the player MUST follow suit, UNLESS they don't have any
	// cards in that suit
	if h.suit != "" && h.suit != cards[0].Suit() {
		if hasSuit(*hand, h.suit) {
			return fmt.Errorf(
				"must follow suit: %s, but player played %s",
				h.suit,
				cards[0].Suit(),
			)
		} else if h.trick == 1 && cards[0].Suit() == SuitHearts {
			if !onlyHasHearts(*hand) {
				return errors.New("cannot play a heart on the first trick")
			}
		}
	} else if !h.brokenHearted && cards[0].Suit() == SuitHearts { // leading with a heart
		if !onlyHasHearts(*hand) {
			return errors.New("cannot lead with a heart until hearts are broken")
		}
	}

	*hand = removeCard(*hand, cards[0])
	*played = &cards[0]

	h.lastPlayed = p

	// look to see if any player have not yet played
	for _, player := range h.Players {
		if player.Played == nil {
			keepPlaying = true // and if not set flag so we can continue
		}
	}

	// if no suit had been led before, then this card must be the new leading suit
	if h.suit == "" {
		h.suit = cards[0].Suit()
	}

	// if no player was found who hasn't played, then the trick is over
	if !keepPlaying {
		if len(h.Players[PlayerOne].Hand) == 0 {
			h.nextRound()
		} else {
			h.nextTrick()
		}
	}

	return nil
}

// currentlyPlaying returns either the player who has the two of clubs, or the last player
// to take a trick
func (h *Hearts) currentlyPlaying() (players []int) {
	if h.lastPlayed != Nobody {
		players = []int{nextPlayer(h.lastPlayed)}

	} else if h.lastTaken != Nobody {
		players = []int{h.lastTaken}

	} else {

		// if no one took the last trick, return the player with two of clubs
		for i, p := range h.Players { // look through players...
			for _, c := range p.Hand { // look through players' hands...
				if c == Card(CardTwoOfClubs) {
					players = []int{i}
					return
				}
			}
		}
	}

	return
}

// nextRound advances to the next phase and increments the round number.
func (h *Hearts) nextRound() {
	h.nextTrick()
	shot := Nobody

	for i := 0; i < len(h.Players); i++ {
		player := &h.Players[i]

		if shot == Nobody { // everyone is taking their round score
			if player.roundScore == 26 { // discovered that someone shot the moon
				shot = i // mark that person as having shot
				i = -1   // reset the loop with someone shooting; -1 so we don't skip 0

			} else { // otherwise, assume no one shot and give everyone their round score
				player.gameScore -= player.roundScore
				player.roundScore = 0
			}

		} else if i != shot { // else another player shot the moon!
			player.gameScore -= 26 // suck 26 points, loser!
		}

		// detect player has crossed the threshhold and ended that game
		if player.gameScore <= 0 {
			h.finished = true
		}
	}

	h.lastTaken = Nobody
	h.phaseEnd = true
	h.round++
	h.NextPhase()
}

// nextTrick cleans up, adds up the points taken in the trick, figures out who takes them
// and sets up for the next trick.
func (h *Hearts) nextTrick() {
	var highestCard Card = Nobody
	highestPlayer := Nobody
	trick := [4]Card{}

	for p, player := range h.Players {
		card := player.Played
		trick[p] = *card

		if card.Suit() == h.suit { // only card that are on suit can take the trick
			if *card > highestCard {
				highestCard = *card // Set the highest card...
				highestPlayer = p   // ...and the player who took it.
			}
		}
	}

	h.lastTaken = highestPlayer
	h.lastTrick = trick
	h.trick += 1
	h.lastPlayed = Nobody
	trickTotal := sumTrickPoints(trick)
	h.Players[highestPlayer].roundScore += trickTotal
}

// check the hand for the given cards. Return true if the cards are in the hand, false if
// they are not.
func hasCard(hand []Card, cards ...Card) bool {
	for _, card := range cards {
		has := false
		for _, heldCard := range hand {
			if card.Compare(heldCard) == 0 {
				has = true
			}
		}

		if !has {
			return false
		}
	}

	return true
}

// hasSuit returns true if the given suit appears in the given hand
func hasSuit(hand []Card, suit string) bool {
	for _, c := range hand {
		if c.Suit() == suit {
			return true
		}
	}

	return false
}

// hasTwoOfClubs returns true if the two of clubs is found in the given hand
func hasTwoOfClubs(hand []Card) bool {
	for _, c := range hand {
		if c == CardTwoOfClubs {
			return true
		}
	}

	return false
}

// nextPlayer returns the index to the "left" of a given player. Left, for our purposes,
// is toward the begining of the slice. If the index is 0, it wraps around to the last
// index.
func nextPlayer(lastPlayer int) int {
	player := lastPlayer - 1

	if player < PlayerOne {
		player = PlayerFour
	}

	return player
}

func onlyHasHearts(hand []Card) bool {
	for _, card := range hand {
		if card.Suit() != SuitHearts {
			return false
		}
	}

	return true
}

// removeCard returns a new array without the given cards in it.
func removeCard(hand []Card, cards ...Card) []Card {
	new := make([]Card, 0, len(hand))
	set := map[Card]bool{}

	for _, c := range cards {
		set[c] = true
	}

	for _, c := range hand {
		if !set[c] {
			new = append(new, c)
		}
	}

	return new
}

// sort sorts a hand using quicksort.
// low is the smallest index to be sorted (probably 0)
// high is the largest index to be sorted (probably len(hand) - 1)
func sort(hand []Card, low int, high int) {
	partition := low
	mid := low

	if low >= high {
		return
	}

	for i := low + 1; i <= high; i++ {
		if hand[i] < hand[partition] {
			mid++
			swap(hand, mid, i)
		}
	}

	swap(hand, mid, partition)

	sort(hand, mid+1, high)
	sort(hand, low, mid-1)
}

// sumTrickPoints takes a trick and returns the sum of all the points contained in it.
func sumTrickPoints(trick [4]Card) int {
	total := 0

	for _, card := range trick {
		if card.Suit() == SuitHearts {
			total += 1
		}

		if card == CardJamoke {
			total += 13
		}
	}

	return total
}

// swap takes a hand and two indices and swaps the values at those indices
func swap(hand []Card, first int, second int) {
	tmp := hand[first]
	hand[first] = hand[second]
	hand[second] = tmp
}
