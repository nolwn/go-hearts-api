package game

// Scorable should be used when a game or round is scored. If the game is played in
// rounds, the score can either represent the round score or the score of the overall
// game, depending on what makes most sense for the game. If Score cannot be called
// during a game or round, then it should either return the previous scores or zeroed out
// scores.
type Scorable interface {

	// Score returns the score of the game or round. The score is represented by a map
	// with an int representing a player as a key and an int representing that player's
	// score.
	Score() map[int]int
}
