package game

// Phase contains the methods needed to manage a card game that is played in phases.
// Card games have phases when different parts of a round change in the way the players
// are supposed to act. For instance, in Hearts there is a pass phase where players pci
// three cards to pass to an opponent, and there there is a play phase where players take
// turns picking a card to play into the middle of the table.
type Phase interface {

	// NextPhase sets the phase for a game that has multiple phases. Hearts, for
	// instance, begins with a passing phase, and then moves into a playing phase.
	//
	// NextPhase returns the phase that that the game has been switched to, and an error
	// if there was a problem changing phases.
	NextPhase() (int error)

	// Phase returns the current phase of the game.
	Phase() int
}
