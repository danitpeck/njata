package game

import (
    "strings"
    "sync"
    "testing"
)

type bufferOutput struct {
    mu    sync.Mutex
    lines []string
}

func (b *bufferOutput) Write(text string) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.lines = append(b.lines, text)
}

func (b *bufferOutput) WriteLine(text string) {
    b.Write(text)
}

func (b *bufferOutput) Contains(substring string) bool {
    b.mu.Lock()
    defer b.mu.Unlock()
    for _, line := range b.lines {
        if strings.Contains(line, substring) {
            return true
        }
    }
    return false
}

func TestBroadcastSay(t *testing.T) {
    world := CreateDefaultWorld()

    aliceOut := &bufferOutput{}
    bobOut := &bufferOutput{}

    if err := world.AddPlayer(&Player{Name: "Alice", Output: aliceOut}); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if err := world.AddPlayer(&Player{Name: "Bob", Output: bobOut}); err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    world.BroadcastSay("Alice", "hello")

    if !aliceOut.Contains("You say 'hello'") {
        t.Fatalf("expected sender message not found")
    }
    if !bobOut.Contains("Alice says 'hello'") {
        t.Fatalf("expected receiver message not found")
    }
}
