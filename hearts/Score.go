package hearts

// Score returns each plays distance from losing. In a standard game of Hearts, players
// play to 100. There are variations, however, where players play to smaller number, for
// instance 75. In order to accommodate any possible target score, players scores start at
// the target score and go down from there. So, in a standard game, each player would
// start at 100 points and their score would fall round after round until someone hits 0.
func (h *Hearts) Score() map[int]int {
	scores := make(map[int]int)

	for p, player := range h.Players {
		scores[p] = player.gameScore
	}

	return scores
}
