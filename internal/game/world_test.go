package game

import "testing"

func TestBroadcastSay(t *testing.T) {
    world := CreateDefaultWorld()

    aliceOut := &bufferOutput{}
    bobOut := &bufferOutput{}

    speaker := &Player{Name: "Alice", Output: aliceOut}
    if err := world.AddPlayer(speaker); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if err := world.AddPlayer(&Player{Name: "Bob", Output: bobOut}); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    world.BroadcastSay(speaker, "hello")

    if !aliceOut.Contains("You say 'hello'") {
        t.Fatalf("expected sender message not found")
    }
    if !bobOut.Contains("Alice says 'hello'") {
        t.Fatalf("expected receiver message not found")
    }
}
