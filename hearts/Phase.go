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
		h.phase = PhasePass
	} else {
		for i := range h.Players {
			h.Players[i].Hand = append(h.Players[i].Hand, h.Players[i].Recieving...)
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
