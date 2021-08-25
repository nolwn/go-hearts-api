package game

// Round represents a turn-based card game. Its methods are generic so that it can be
// reused for different card games.
//
// When the game state makes a method call impossible, the method may throw an error.
// Otherwise, methods should be reliable and should not throw errors.
type Round interface {

	// Finished returns true if the game or round has finished and false if it is still
	// in progress. A game should not be marked finished unless is over, or needs to be
	// setup again for another round. Finished should not return true if the game is just
	// moving to another phase.
	//
	// Most often, when a game is finished either scores are updated or a winnder is
	// declared.
	Finished() bool

	// NextPhase sets the phase for a game that has multiple phases. Hearts, for
	// instance, begins with a passing phase, and then moves into a playing phase. Your
	// game should mark a new phase whenever there is a change in the way the
	// players are supposed to act.
	//
	// NextPhase returns the phase that that the game has been switched to, and an error
	// if there was a problem changing phases.
	NextPhase() (int error)

	// Phase returns the current phase of the game.
	Phase() int

	// Play plays a card. What that means differs from game to game, and phase to phase.
	// It might mean that a card is placed face up infront of a player, or it might mean
	// that it is passed to another player, or it might mean that it is traded in for
	// another card.
	//
	// Play takes a player and a card which are integers. If that player cannot play,
	// or that card cannot be played, then an error should be returned.
	Play(player int, card ...Card) error

	// Player returns an int that represents the player whose turn it is.
	Player() int

	// Round returns the round number. The first round returns 1.
	Round() int

	// Setup sets up a table for a new game or round. If your game has different phases
	// or parts that might affect the way setup happens, you should first set that with
	// the Phase method.
	Setup() error

	// State returns the current game state. It should return a data structure that has
	// enough information that the programmer can understand exactly what is happening
	// in the game, and the game can be completely recreated by the data returned.
	State() (state interface{})
}
