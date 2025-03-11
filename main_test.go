package main

import (
	"testing"
)

func TestSanity(t *testing.T) {
	got := 2
	want := 2

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
