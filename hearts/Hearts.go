package hearts

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

// Hearts is the underlying data of the game. It should be storable in the database with
// few, if any, modifications.
type Hearts struct {

	// These are the four players playing the game. Exactly four players are required in
	// this version of Hearts.
	Players [4]Player

	// finished is a flag that signifies that that the current round has ended. A round
	// is considered ended when every player has played every card in their hands.
	finished bool

	// phase is an int that represents the phase of the game. There are two phases in
	// Hearts, the pass phase (which is 0) and the play phase (which is 1).
	phase int

	// phaseEnd is a flag that signals that the current phase has ended.
	phaseEnd bool

	// round is the round number that is currently being played. round starts with 1.
	round int

	// lastTrick is the index of the last player who took a trick
	lastTrick int
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

	// hasPassed is a flag that signals that a player has chosen three cards to pass
	hasPassed bool
}

// New creates a new game of Hearts. It should represent a whole game, not just a round.
func New() Hearts {
	players := [4]Player{}
	round := Hearts{}

	for i := 0; i < 4; i++ {
		players[i] = Player{
			Hand:  make([]Card, 0, 13),
			Taken: make([]Card, 0, 13),
		}
	}

	round.Players = players
	round.round = 1

	return round
}
