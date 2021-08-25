package hearts

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Player represents a players hand, the tricks they've taken, and the card that was
// played that round. Because different games handle players so differently, it probably
// makes sense for this type to be defined purely at the game individual game level.
type Player struct {

	// Hand represents the player's hand. At the beginning of each round, each player is
	// dealt 13 cards (one quarter of the deck). During the pass phase, players select
	// three cards from their hand to the hand of one of their opponents.
	//
	// During the play phase, each player selects a card from their hand to play into the
	// middle of the table. In this struct, that is represented by removing the card from
	// Hand and setting Played to that card.
	Hand []Card

	// Taken represents the tricks that the player has taken. Each trick consists of four
	// cards which have been played by each player each round. At the end of the game,
	// cards which have a point value will be totaled and added to each players score.
	Taken []Card

	// Played is the card that a player has chosen to play for the round. In a physical
	// game, this card would be played into the middle of the table. It is kept with the
	// player in this struct so that, at the end of the round, each card can be easily
	// connected with the player that played it.
	Played *Card
}

// Hearts is the underlying data of the game. It should be storable in the database with
// few, if any, modifications.
type Hearts struct {

	// These are the four players playing the game. Exactly four players are required in
	// this version of Hearts.
	Players [4]Player

	// finished is a flag that signifies that that the current round has ended. A round
	// is considered ended when every player has played every card in their hands.
	finished bool

	// phase is an int that represents the phase of the game. There are two phases in
	// Hearts, the pass phase (which is 0) and the play phase (which is 1).
	phase int

	// phaseEnd is a flag that signals that the current phase has ended.
	phaseEnd bool

	// round is the round number that is currently being played. round starts with 1.
	round int

	// turn is the index of the player whose turn it is.
	turn int
}

const (

	// PhasePass begins with each player having 13 cards in their Hand, 0 cards in their
	// Taken, and no played cards. Where players pass cards depends on what round it is:
	//
	// Round 1: players pass three cards to the left (their index -1).
	// Round 2: players pass three cards to the right (their index +1).
	// Round 3: players pass three cards across (their index + 2).
	// Round 4: players don't pass any cards, the phase ends immediately.
	//
	// After round 4, the pattern starts again from round 1.
	PhasePass = iota

	// PhasePlay is the phase where players take turns playing cards and then one player
	// picks a trick.
	//
	// The first player to play is the player who last took a trick or, if it is the very
	// first round, the player who holds the two of clubs. Each player must then play
	PhasePlay
)

// New creates a new game of Hearts. It should represent a whole game, not just a round.
func New() Hearts {
	players := [4]Player{}
	round := Hearts{}

	for i := 0; i < 4; i++ {
		players[i] = Player{
			Hand:  make([]Card, 0, 13),
			Taken: make([]Card, 0, 13),
		}
	}

	round.Players = players
	round.round = 1

	return round
}

// Finished will return true when a round (not the full game) is completed.
func (h *Hearts) Finished() bool {
	return h.finished
}

// NextPhase will toggle between the pass phase and the play phase. NextPhase can only be
// called once a phase has ended. An error will be returned if the phase has not ended.
//
// The pass phase is considered ended when each player has picked the three cards they
// want to pass. When NextPhase is called, thoes cards have been moved to the appropriate
// players' hands.
//
// The play phase is considered ended when each player has played one card. When
// NextPhase is called, the Played cards are moved into the Taken slice of the player
// who played the highest on suit card.
func (h *Hearts) NextPhase() (int error) {
	if !h.phaseEnd {
		return errors.New("cannot end phase")
	}

	if h.phase == PhasePlay {
		h.phase = PhasePass
	} else {
		h.phase++
	}

	return nil
}

// Phase returns an int that represents the phase that the game is on. The pass phase is
// is 0, the play phase is 1.
func (h *Hearts) Phase() int {
	return h.phase
}

// Play in Hearts means one of two things depending on the phase. In the pass
// phase, players pick three cards to pass. In the play phase, players pick one card
// to play into trick.
//
// An error will be returned if it isn't the players turn to play.
func (h *Hearts) Play(player int, card ...Card) error {
	if player != h.turn {
		return fmt.Errorf("it is not player %d's turn", player)
	}

	if h.Phase() != PhasePlay {
		return errors.New("it isn't the play phase")
	}

	if len(card) != 1 {
		return errors.New("")
	}

	hand := &h.Players[player].Hand
	played := &h.Players[player].Played

	if !hasCard(*hand, card[0]) {
		return fmt.Errorf("player %d does not have %d", player, card)
	}

	*hand = removeCard(*hand, card[0])
	*played = &card[0]

	return nil
}

// Round returns the round number.
func (h *Hearts) Round() int {
	return h.round
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
	h.finished = false

	return nil
}

// clearTaken clears out any tricks taken by each of the players
func (s *Hearts) clearTaken() {
	for _, p := range s.Players {
		p.Taken = make([]Card, 0, 13)
	}
}

func (s *Hearts) deal() {
	var n Card = 0 // cards in a deck, starting with 0
	i := 0         // player index
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

		s.Players[i].Hand = append(s.Players[i].Hand, card)

		if i >= 3 {
			i = 0
		} else {
			i++
		}
	}

	// the loop ends one card early because we can't generate a random number between 0
	// and 0. We just need to scoop up that last card.
	s.Players[i].Hand = append(s.Players[i].Hand, deck[0])
}

// check the hand for the given card. Return true if the card is in the hand, false if it
// is not.
func hasCard(hand []Card, card Card) bool {
	for _, c := range hand {
		if c.Compare(card) == 0 {
			return true
		}
	}

	return false
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
