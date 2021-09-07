package hearts

const (
	SuitDiamonds = "Diamonds"
	SuitClubs    = "Clubs"
	SuitHearts   = "Hearts"
	SuitSpades   = "Spades"
)

type Card int

// Compare this card against a given card. If the given card is bigger, it will return
// a negative number, if the given card is the same it will return 0 and if it's smaller
// it will return a positive number.
//
// It doesn't matter, in Hearts, what the difference in value is between two cards of
// different suits. Because of that, there are no guarantees about what will happen if
// you try to compare cards of different suits.
func (c Card) Compare(other Card) int {
	return int(c - other)
}

// Suit returns the cards suit
func (c Card) Suit() string {
	if c < 13 {
		return SuitDiamonds
	} else if c < 26 {
		return SuitClubs
	} else if c < 39 {
		return SuitHearts
	} else {
		return SuitSpades
	}
}

// Value returns the value of the card
func (c Card) Value() string {
	values := []string{
		"Ace", // Ace looks like the smallest but it's the largest
		"Two",
		"Three",
		"Four",
		"Five",
		"Six",
		"Seven",
		"Eight",
		"Nine",
		"Ten",
		"Jack",
		"Queen",
		"King",
	}
	idx := (c + 1) % 13 // 13 % 13 == 0 which is why Ace is actually the biggest

	return values[idx]
}
