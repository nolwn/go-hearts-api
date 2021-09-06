package game

// Round represents a turn-based card game. Its methods are generic so that it can be
// reused for different card games.
//
// When the game state makes a method call impossible, the method may throw an error.
// Otherwise, methods should be reliable and should not throw errors.
type CardGame interface {

	// Play plays a card. What that means differs from game to game, and phase to phase.
	// It might mean that a card is placed face up infront of a player, or it might mean
	// that it is passed to another player, or it might mean that it is traded in for
	// another card.
	//
	// Play takes a player and a card which are integers. If that player cannot play,
	// or that card cannot be played, then an error should be returned.
	Play(player int, card ...Card) error

	// Players returns an int that represents the player whose turn it is.
	PlayersTurn() []int

	// Setup sets up a table for a new game or round. If your game has different phases
	// or parts that might affect the way setup happens, you should first set that with
	// the Phase method.
	Setup() error

	// State returns the current game state. It should return a data structure that has
	// enough information that the programmer can understand exactly what is happening
	// in the game, and the game can be completely recreated by the data returned.
	State() (state interface{})
}
