package game

import (
    "fmt"
    "sort"
    "strings"
    "sync"
    "unicode"
)

type Player struct {
    Name       string
    Output     Output
    Disconnect func(reason string)
}

type Room struct {
    Name        string
    Description string
}

type World struct {
    mu      sync.RWMutex
    room    Room
    players map[string]*Player
}

func CreateDefaultWorld() *World {
    return &World{
        room: Room{
            Name:        "The Crossroads",
            Description: "A simple stone path crosses here, leading to all corners of the land.",
        },
        players: map[string]*Player{},
    }
}

func ValidateName(name string) error {
    if len(name) < 3 || len(name) > 16 {
        return fmt.Errorf("name must be 3-16 characters")
    }

    for _, r := range name {
        if r > 127 || (!unicode.IsLetter(r) && !unicode.IsDigit(r)) {
            return fmt.Errorf("name must be letters or digits")
        }
    }

    return nil
}

func (w *World) RoomSnapshot() Room {
    w.mu.RLock()
    defer w.mu.RUnlock()
    return w.room
}

func (w *World) AddPlayer(player *Player) error {
    if player == nil {
        return fmt.Errorf("player is nil")
    }

    if err := ValidateName(player.Name); err != nil {
        return err
    }

    key := normalizeName(player.Name)
    w.mu.Lock()
    defer w.mu.Unlock()

    if _, exists := w.players[key]; exists {
        return fmt.Errorf("name already in use")
    }

    w.players[key] = player
    return nil
}

func (w *World) RemovePlayer(name string) {
    key := normalizeName(name)
    w.mu.Lock()
    defer w.mu.Unlock()
    delete(w.players, key)
}

func (w *World) PlayersSnapshot() []*Player {
    w.mu.RLock()
    defer w.mu.RUnlock()

    players := make([]*Player, 0, len(w.players))
    for _, player := range w.players {
        players = append(players, player)
    }

    return players
}

func (w *World) ListPlayers() []string {
    players := w.PlayersSnapshot()
    names := make([]string, 0, len(players))
    for _, player := range players {
        names = append(names, player.Name)
    }

    sort.Strings(names)
    return names
}

func (w *World) ListPlayersExcept(name string) []string {
    players := w.PlayersSnapshot()
    names := make([]string, 0, len(players))
    for _, player := range players {
        if !strings.EqualFold(player.Name, name) {
            names = append(names, player.Name)
        }
    }

    sort.Strings(names)
    return names
}

func (w *World) BroadcastSay(speaker string, message string) {
    players := w.PlayersSnapshot()
    for _, player := range players {
        if strings.EqualFold(player.Name, speaker) {
            player.Output.WriteLine(fmt.Sprintf("You say '%s'", message))
            continue
        }
        player.Output.WriteLine(fmt.Sprintf("%s says '%s'", speaker, message))
    }
}

func (w *World) BroadcastSystemExcept(except string, message string) {
    players := w.PlayersSnapshot()
    for _, player := range players {
        if strings.EqualFold(player.Name, except) {
            continue
        }
        player.Output.WriteLine(message)
    }
}

func normalizeName(name string) string {
    return strings.ToLower(strings.TrimSpace(name))
}
