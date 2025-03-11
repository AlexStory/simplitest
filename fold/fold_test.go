package fold

import "testing"

func TestNestedSanity(t *testing.T) {
	got := 2
	want := 2

	if got != want {
		t.Errorf("Got %d, want %d", got, want)
	}
}
