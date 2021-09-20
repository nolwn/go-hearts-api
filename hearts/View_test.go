package hearts

import (
	"fmt"
	"testing"
)

func TestFrom(t *testing.T) {
	hearts := New()
	hearts.deal()

	b, err := hearts.From(0)

	if err != nil {
		t.Errorf("expected no error but received: %s", err)
	}

	s := string(b)

	fmt.Println(s)

	if b == nil {
		t.Error("expected a byte array, but received nil")
	}
}
