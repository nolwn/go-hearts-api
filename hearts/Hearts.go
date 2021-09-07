package hearts

// Hearts is the underlying data of the game. It should be storable in the database with
// few, if any, modifications.
type Hearts struct {

	// These are the four players playing the game. Exactly four players are required in
	// this version of Hearts.
	Players [4]Player

	// finished is a flag that signifies that that the current round has ended. A round
	// is considered ended when every player has played every card in their hands.
	finished bool

	//lastPlayed is the index of the last player who played a card
	lastPlayed int

	// lastTrick is the index of the last player who took a trick
	lastTrick int

	// phase is an int that represents the phase of the game. There are two phases in
	// Hearts, the pass phase (which is 0) and the play phase (which is 1).
	phase int

	// phaseEnd is a flag that signals that the current phase has ended.
	phaseEnd bool

	// round is the round number that is currently being played. round starts with 1.
	round int

	// suit is the suit of the first card played into the trick. It is the suit that must
	// be followed
	suit string
}

// Player represents a players hand, the tricks they've taken, and the card that was
// played that round. Because different games handle players so differently, it probably
// makes sense for this type to be defined purely at the game individual game level.
type Player struct {

	// Hand represents the player's hand. At the beginning of each round, each player is
	// dealt 13 cards (one quarter of the deck). During the pass phase, players select
	// three cards from their hand to the hand of one of their opponents.
	//
	// During the play phase, each player selects a card from their hand to play into the
	// middle of the table. In this struct, that is represented by removing the card from
	// Hand and setting Played to that card.
	Hand []Card

	// Taken represents the tricks that the player has taken. Each trick consists of four
	// cards which have been played by each player each round. At the end of the game,
	// cards which have a point value will be totaled and added to each players score.
	Taken []Card

	// Played is the card that a player has chosen to play for the round. In a physical
	// game, this card would be played into the middle of the table. It is kept with the
	// player in this struct so that, at the end of the round, each card can be easily
	// connected with the player that played it.
	Played *Card

	// Recieving represents the cards that are being passed to this player during the
	// passing phase. After a player has picked three cards to pass, they get moved
	// into the correct player's Recieving slice.
	Recieving []Card

	// hasPassed is a flag that signals that a player has chosen three cards to pass.
	hasPassed bool

	// roundScore keeps track of a players score as the round goes on
	roundScore int
}

// New creates a new game of Hearts. It should represent a whole game, not just a round.
func New() Hearts {
	players := [4]Player{}
	hearts := Hearts{}

	for i := 0; i < 4; i++ {
		players[i] = Player{
			Hand:  make([]Card, 0, 13),
			Taken: make([]Card, 0, 13),
		}
	}

	hearts.Players = players
	hearts.round = 1
	hearts.lastPlayed = -1
	hearts.lastTrick = -1

	return hearts
}
