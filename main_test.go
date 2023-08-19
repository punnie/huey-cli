package main

import (
	"testing"
)

func TestGreet(t *testing.T) {
	want := "Hello, Pedro!"
	if got := Greet("Pedro"); got != want {
		t.Errorf("Greet() = %q, want %q", got, want)
	}
}
