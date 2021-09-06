package hearts

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

// Play in Hearts means one of two things depending on the phase. In the pass
// phase, players pick three cards to pass. In the play phase, players pick one card
// to play into trick.
//
// An error will be returned if it isn't the players turn to play.
func (h *Hearts) Play(player int, cards ...Card) error {
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
		return h.passPlayers()
	} else {
		return h.playPlayers()
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
	h.lastTrick = -1
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

// passAcross returns the target across the table
func (h *Hearts) passAcross(player int, cards []Card) int {
	switch player {
	case 0:
		return 2
	case 1:
		return 3
	case 2:
		return 0
	default: // case 3:
		return 1
	}
}

// passLeft returns the target to the players left (toward begining of array)
func (h *Hearts) passLeft(player int, cards []Card) int {
	if player == 0 {
		return 3
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
				return errors.New("player must play 3 differnt cards")
			}
		}
	}

	if !hasCard(*playerHand, cards...) {
		return errors.New("player must have the cards to pass them")
	}

	switch pass {
	case 0:
		return errors.New("plater cannot pass on the hold round")
	case 1:
		target = h.passLeft(player, cards)
	case 2:
		target = h.passRight(player, cards)
	case 3:
		target = h.passAcross(player, cards)
	}

	*playerHand = removeCard(*playerHand, cards...)
	h.Players[player].hasPassed = true
	h.Players[target].Recieving = cards

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

// passPlayers returns the players who have not yet picked cards to pass
func (h *Hearts) passPlayers() (players []int) {
	for i, p := range h.Players {
		if !p.hasPassed {
			players = append(players, i)
		}
	}

	return
}

// passRight returns the target to the players right (toward end of array)
func (h *Hearts) passRight(player int, cards []Card) int {
	if player == 3 {
		return 0
	} else {
		return player + 1
	}
}

// playPhase contains the logic for playing a card during the play phase.
func (h *Hearts) playPhase(player int, cards ...Card) error {
	if len(cards) != 1 {
		return errors.New("player must play exactly one card")
	}

	hand := &h.Players[player].Hand
	played := &h.Players[player].Played

	if !hasCard(*hand, cards[0]) {
		return fmt.Errorf("player %d does not have %d", player, cards)
	}

	*hand = removeCard(*hand, cards[0])
	*played = &cards[0]

	return nil
}

// playPlayers returns either the player who has the two of clubs, or the last player
// to take a trick
func (h *Hearts) playPlayers() (players []int) {
	twoOfClubs := 13

	if h.lastTrick != -1 {
		players = []int{h.lastTrick}
	} else {

		// if no one took the last trick, return the player with two of clubs
		for i, p := range h.Players { // look through players...
			for _, c := range p.Hand { // look through players' hands...
				if c == Card(twoOfClubs) {
					players = []int{i}
					return
				}
			}
		}
	}

	return
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
