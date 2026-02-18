package game

import (
    "strings"
    "sync"
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
