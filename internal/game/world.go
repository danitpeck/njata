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
    Location   int
}

type Room struct {
    Vnum        int
    Name        string
    Description string
    Exits       map[string]int
    ExDescs     map[string]string
}

type RoomView struct {
    Name        string
    Description string
    Exits       []string
    Others      []string
}

type World struct {
    mu      sync.RWMutex
    rooms   map[int]*Room
    start   int
    players map[string]*Player
}

func CreateDefaultWorld() *World {
    defaultRoom := &Room{
        Vnum:        1,
        Name:        "The Crossroads",
        Description: "A simple stone path crosses here, leading to all corners of the land.",
        Exits:       map[string]int{},
        ExDescs:     map[string]string{},
    }

    return &World{
        rooms: map[int]*Room{defaultRoom.Vnum: defaultRoom},
        start: defaultRoom.Vnum,
        players: map[string]*Player{},
    }
}

func CreateWorldFromRooms(rooms map[int]*Room, start int) *World {
    if len(rooms) == 0 {
        return CreateDefaultWorld()
    }

    if start == 0 {
        for vnum := range rooms {
            if start == 0 || vnum < start {
                start = vnum
            }
        }
    }

    return &World{
        rooms:   rooms,
        start:   start,
        players: map[string]*Player{},
    }
}

func (w *World) StartRoom() int {
    w.mu.RLock()
    defer w.mu.RUnlock()
    return w.start
}

func (w *World) HasRoom(vnum int) bool {
    w.mu.RLock()
    defer w.mu.RUnlock()
    _, ok := w.rooms[vnum]
    return ok
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
    if room, ok := w.rooms[w.start]; ok {
        return *room
    }
    return Room{}
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

    if player.Location == 0 {
        player.Location = w.start
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

func (w *World) DescribeRoom(player *Player) (RoomView, error) {
    w.mu.RLock()
    defer w.mu.RUnlock()

    room, ok := w.rooms[player.Location]
    if !ok {
        return RoomView{}, fmt.Errorf("room not found")
    }

    others := make([]string, 0, len(w.players))
    for _, other := range w.players {
        if other.Location == room.Vnum && !strings.EqualFold(other.Name, player.Name) {
            others = append(others, other.Name)
        }
    }
    sort.Strings(others)

    exits := make([]string, 0, len(room.Exits))
    for exit := range room.Exits {
        exits = append(exits, exit)
    }
    sort.Strings(exits)

    return RoomView{
        Name:        room.Name,
        Description: room.Description,
        Exits:       exits,
        Others:      others,
    }, nil
}

func (w *World) FindRoomExDesc(player *Player, keyword string) (string, bool) {
    w.mu.RLock()
    defer w.mu.RUnlock()

    room, ok := w.rooms[player.Location]
    if !ok {
        return "", false
    }

    if room.ExDescs == nil {
        return "", false
    }

    key := strings.ToLower(strings.TrimSpace(keyword))
    if key == "" {
        return "", false
    }

    value, ok := room.ExDescs[key]
    return value, ok
}

func (w *World) MovePlayer(player *Player, direction string) (RoomView, error) {
    w.mu.Lock()
    room, ok := w.rooms[player.Location]
    if !ok {
        w.mu.Unlock()
        return RoomView{}, fmt.Errorf("room not found")
    }

    targetVnum, ok := room.Exits[direction]
    if !ok {
        w.mu.Unlock()
        return RoomView{}, fmt.Errorf("no exit")
    }

    targetRoom, ok := w.rooms[targetVnum]
    if !ok {
        w.mu.Unlock()
        return RoomView{}, fmt.Errorf("exit leads nowhere")
    }

    player.Location = targetRoom.Vnum
    w.mu.Unlock()

    return w.DescribeRoom(player)
}

func (w *World) BroadcastSay(speaker *Player, message string) {
    w.mu.RLock()
    location := speaker.Location
    players := make([]*Player, 0, len(w.players))
    for _, player := range w.players {
        if player.Location == location {
            players = append(players, player)
        }
    }
    w.mu.RUnlock()

    for _, player := range players {
        if strings.EqualFold(player.Name, speaker.Name) {
            player.Output.WriteLine(fmt.Sprintf("You say '%s'", message))
            continue
        }
        player.Output.WriteLine(fmt.Sprintf("%s says '%s'", speaker.Name, message))
    }
}

func (w *World) BroadcastSystemToRoomExcept(except *Player, message string) {
    w.mu.RLock()
    location := except.Location
    players := make([]*Player, 0, len(w.players))
    for _, player := range w.players {
        if player.Location == location && !strings.EqualFold(player.Name, except.Name) {
            players = append(players, player)
        }
    }
    w.mu.RUnlock()

    for _, player := range players {
        player.Output.WriteLine(message)
    }
}

func normalizeName(name string) string {
    return strings.ToLower(strings.TrimSpace(name))
}
