package hearts

import "errors"

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
