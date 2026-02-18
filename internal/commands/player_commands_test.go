package commands

import "testing"

func TestFormatExits(t *testing.T) {
    exits := []string{"west", "north", "up", "southeast"}
    got := FormatExits(exits)
    want := "Exits: North West Up Southeast."
    if got != want {
        t.Fatalf("expected %q, got %q", want, got)
    }
}

func TestFormatExitsNone(t *testing.T) {
    got := FormatExits(nil)
    want := "Exits: none"
    if got != want {
        t.Fatalf("expected %q, got %q", want, got)
    }
}
