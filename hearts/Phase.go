package hearts

import "errors"

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
		h.round++ // each passing phase signifies the start of a new round

		// every fourth round skips the passing phase
		if !(h.round%4 == 0) {
			h.phase = PhasePass
		}

	} else {
		// put all receiving cards into players hands before the passing phase
		for i := range h.Players {
			h.Players[i].Hand = mergeCards(h.Players[i].Hand, h.Players[i].Recieving)
			h.Players[i].Recieving = []Card{}
		}

		h.phase++
	}

	h.finished = false

	return nil
}

// Phase returns an int that represents the phase that the game is on. The pass phase is
// is 0, the play phase is 1.
func (h *Hearts) Phase() int {
	return h.phase
}

// mergeCards takes a hand and some cards and returns a sorted slice which contains both.
// mergeCards assumes that the hand is already sorted, but that the cards are not.
func mergeCards(hand []Card, cards []Card) []Card {
	// the new hand will be the size of the old hand, plus the cards it's receiving
	capacity := len(hand) + len(cards)
	newHand := make([]Card, 0, capacity)

	// hand should already be sorted, but the cards might not be
	sort(cards, 0, len(cards)-1)

	h := 0 // hand index
	c := 0 // cards index

	for i := 0; i < capacity; i++ {
		if h >= len(hand) { // all hand cards are already added...
			newHand = append(newHand, cards[c])
		} else if c >= len(cards) { // ...all cards are already added..
			newHand = append(newHand, hand[h])
		} else if hand[h] < cards[c] { // ... the next hand card is smaller...
			newHand = append(newHand, hand[h])
			h++
		} else { // ...otherwise, the next card is smaller.
			newHand = append(newHand, cards[c])
			c++
		}
	}

	return newHand
}
