package game

type View interface {

	// From takes a player position returns a JSON representation of what that player can
	// see. It should not reveal things that should not be visible to that player, for
	// instance hidden cards in other players' hands.
	//
	// It should not be incumbent on the client to have to figure out too much about the
	// game, so things like flags that would be helpful to the client should also be
	// returned. Remember that if a client resets, they might lose their "memory" of the
	// game so be generous about what non-private information you share.
	From() string
}
