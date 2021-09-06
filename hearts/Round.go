package hearts

// Finished will return true when a round (not the full game) is completed.
func (h *Hearts) Finished() bool {
	return h.finished
}

// Round returns the round number.
func (h *Hearts) Round() int {
	return h.round
}
