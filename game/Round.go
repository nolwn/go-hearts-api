package game

// Round contains the methods needed to manage a game that is played in rounds.
// Rounds occur in games that can be divided into repeating segments.
type Round interface {

	// Round returns the round number. The first round returns 1.
	Round() int
}
