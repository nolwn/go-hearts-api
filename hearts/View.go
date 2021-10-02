package hearts

import "encoding/json"

type JSONCard struct {
	Suit  string `json:"suit"`
	Value string `json:"value"`
}

type Perspective struct {

	// Broken is set to true if hearts have been sloughed. The Jamoke does not
	// count as a heart.
	Broken bool `json:"brokenHearted"`

	// Finished keeps track of whether the game has ended or not
	Finished bool `json:"finished"`

	// Hand is the hand of the player being viewed.
	Hand []JSONCard `json:"hand"`

	// HasPassed is a slice of player ids representing the players who have passed their
	// cards during the passing phase.
	HasPassed []int `json:"hasPassed,omitempty"`

	// LastTrick are the cards played in the last trick
	LastTrick []JSONCard `json:"lastTrick,omitempty"`

	// PassTo is a string which can either be `left`, `right`, `across` or `hold`.
	PassTo string `json:"passTo,omitempty"`

	// Phase is an int that represents the Phase of the game. There are two phases in
	// Hearts, the pass Phase (which is 0) and the play Phase (which is 1).
	Phase string `json:"phase"`

	// Round is the Round number that is currently being played. Round starts with 1.
	Round int `json:"round"`

	// suit is the suit of the first card played into the trick. It is the suit that must
	// be followed.
	Suit string `json:"suit,omitempty"`

	// ThisTrick is the cards that have been played into the trick so far
	ThisTrick []JSONCard `json:"thisTrick,omitempty"`

	// Turn is the id of the player whose turn it is.
	Turn int `json:"turn,omitempty"`

	// Took is the id of the last player who took a trick. It starts at 1.
	Took int `json:"took,omitempty"`

	// Winner is the player who won the game if the game has finished.
	Winner []int `json:"winner,omitempty"`
}

func (h *Hearts) From(player int) ([]byte, error) {
	per := Perspective{
		Broken:    h.brokenHearted,
		Finished:  h.finished,
		Hand:      cardsToJSONCards(h.Players[player].Hand...),
		HasPassed: playersToHasPassed(h.Players),
		LastTrick: getLastTrick(h.trick, h.lastTrick),
		PassTo:    roundToPassDirection(h.round),
		Phase:     phaseToJSONPhase(h.phase),
		Round:     h.round,
		Suit:      h.suit,
		ThisTrick: playersToThisTrick(h.Players),
		Turn:      getToTurn(h.Phase(), h.PlayersTurn()),
		Took:      h.lastTaken + 1,
		Winner:    playerIndicesToIDs(h.Winner()),
	}

	b, err := json.Marshal(per)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func cardsToJSONCards(cards ...Card) []JSONCard {
	JSONCards := make([]JSONCard, 0, 13)

	for _, card := range cards {
		JSONCards = append(JSONCards, JSONCard{Suit: card.Suit(), Value: card.Value()})
	}

	return JSONCards
}

func getLastTrick(trick int, lastTrick [4]Card) []JSONCard {
	if trick > 1 {
		return cardsToJSONCards(lastTrick[:]...)
	} else {
		return []JSONCard{}
	}
}

func getToTurn(phase int, playersTurn []int) int {
	if phase == PhasePass {
		return 0
	} else {
		return playersTurn[0] + 1
	}
}

func phaseToJSONPhase(phase int) string {
	if phase == PhasePass {
		return "pass"
	} else {
		return "play"
	}
}

func playersToHasPassed(players [4]Player) []int {
	hasPassed := make([]int, 0, 4)

	for p, player := range players {
		if player.hasPassed {
			// players are marshalled by id (starting at 1) not index (starting at 0)
			hasPassed = append(hasPassed, p+1)
		}
	}

	return hasPassed
}

func playersToThisTrick(players [4]Player) []JSONCard {
	trick := make([]Card, 0, 3)

	for _, player := range players {
		if player.Played != nil {
			trick = append(trick, *player.Played)
		}
	}

	return cardsToJSONCards(trick...)
}

func roundToPassDirection(round int) string {
	switch round % 4 {
	case 1:
		return "left"
	case 2:
		return "right"
	case 3:
		return "across"
	default: // case 0:
		return "hold"
	}
}

func playerIndicesToIDs(indices []int) []int {
	IDs := make([]int, 0, 4)

	for _, index := range indices {
		IDs = append(IDs, index+1)
	}

	return IDs
}
