package game

// Card represents a playing card.
type Card interface {
	// Value is the number or value of the card. Value should be written out fully in
	// letters (no digits). Values should be capitalized (e.g Ten, Jack, Ace, etc.).
	// The value returned is for naming purposes and is not suitable for comparing card
	// values.
	Value() string

	// Suit should return a capitalized name for the suit of the card (Spades, Diamonds
	// Hearts, Clubs).
	Suit() string

	// Compare takes a card and returns 0 if the given Card is equal to this Card, a
	// negative number if the given Card is greater than this Card, and a positive number
	// if the given Card is less than this Card.
	Compare(Card) int
}
