package game

import (
    "fmt"
    "sort"
    "strings"
    "sync"
    "time"
    "unicode"

    "njata/internal/skills"
)

type Player struct {
    Name       string
    Output     Output
    Disconnect func(reason string)
    Location   int
    AutoExits  bool
    
    // Character attributes
    Class      int // Index into legacy/classes
    Race       int // Index into legacy/races
    Sex        int // 0=neuter, 1=male, 2=female
    Age        int // 0=child, 1=youth, 2=adult, 3=middle-aged, 4=elderly
    Level      int
    
    // Vital stats
    HP         int
    MaxHP      int
    Mana       int
    MaxMana    int
    Move       int
    MaxMove    int
    Gold       int
    Experience int
    
    // Attribute scores (STR, INT, WIS, DEX, CON, LCK, CHM)
    Attributes [7]int
    
    // Combat stats
    Alignment  int
    Hitroll    int
    Damroll    int
    Armor      int
    
    // Skills tracking
    Skills     map[int]*skills.PlayerSkillProgress // spell_id -> proficiency progress
    
    // Keeper flag - player who maintains the world
    IsKeeper   bool
}

type Mobile struct {
    Vnum        int
    Keywords    []string
    Short       string
    Long        string
    Race        string
    Class       string
    Position    string
    Gender      string
    Level       int
    MaxHP       int
    HP          int
    Mana        int
    MaxMana     int
    Attributes  [7]int // STR, INT, WIS, DEX, CON, LCK, CHA
}

type Object struct {
    Vnum        int
    Keywords    []string
    Type        string
    Short       string
    Long        string
    Weight      int
    Value       [4]int // [0]=quantity [1]=unused [2]=unused [3]=spell_id for magical items
    Flags       map[string]bool
}

type Room struct {
    Vnum        int
    Name        string
    Description string
    Sector      string
    Flags       map[string]bool
    Exits       map[string]int
    ExDescs     map[string]string
    AreaName    string
    AreaAuthor  string
    AreaResetMinutes int
    Mobiles     []*Mobile
    Objects     []*Object
    MobileResets  []Reset
    ObjectResets  []Reset
}

type Reset struct {
    MobVnum    int // for Mobile resets
    ObjVnum    int // for Object resets
    Count      int // how many to load
    Room       int // which room to load into
}

type RoomView struct {
    Name        string
    Description string
    Exits       []string
    Others      []string
    Mobiles     []string // NPC descriptions
    Objects     []string // Object descriptions
    AreaName    string
    AreaAuthor  string
}

type World struct {
    mu               sync.RWMutex
    rooms            map[int]*Room
    start            int
    players          map[string]*Player
    mobiles          map[int]*Mobile  // Prototypes for respawning
    objects          map[int]*Object  // Prototypes for respawning
    areaLastRespawn  map[string]time.Time  // Track when each area last respawned
}

func CreateDefaultWorld() *World {
    defaultRoom := &Room{
        Vnum:        1,
        Name:        "The Crossroads",
        Description: "A simple stone path crosses here, leading to all corners of the land.",
        Sector:      "",
        Flags:       map[string]bool{},
        Exits:       map[string]int{},
        ExDescs:     map[string]string{},
    }

    return &World{
        rooms:           map[int]*Room{defaultRoom.Vnum: defaultRoom},
        start:           defaultRoom.Vnum,
        players:         map[string]*Player{},
        mobiles:         map[int]*Mobile{},
        objects:         map[int]*Object{},
        areaLastRespawn: map[string]time.Time{},
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
        rooms:           rooms,
        start:           start,
        players:         map[string]*Player{},
        mobiles:         map[int]*Mobile{},
        objects:         map[int]*Object{},
        areaLastRespawn: map[string]time.Time{},
    }
}

// CreateWorldFromRoomsWithPrototypes creates a world with prototype data for respawning
func (w *World) SetPrototypes(mobiles map[int]*Mobile, objects map[int]*Object) {
    w.mu.Lock()
    defer w.mu.Unlock()
    w.mobiles = mobiles
    w.objects = objects
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

// FindPlayer retrieves a player by name (case-insensitive)
func (w *World) FindPlayer(name string) (*Player, bool) {
    w.mu.RLock()
    defer w.mu.RUnlock()

    for playerName, player := range w.players {
        if strings.EqualFold(playerName, name) {
            return player, true
        }
    }

    return nil, false
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

    // Format mobiles
    mobiles := make([]string, 0, len(room.Mobiles))
    for _, mob := range room.Mobiles {
        if mob.Short != "" {
            position := strings.TrimSpace(mob.Position)
            if position == "" {
                position = "standing"
            }
            mobiles = append(mobiles, mob.Short+" is "+position+".")
        }
    }
    sort.Strings(mobiles)

    // Format objects
    objects := make([]string, 0, len(room.Objects))
    for _, obj := range room.Objects {
        if obj.Short != "" {
            objects = append(objects, obj.Short+" is here.")
        }
    }
    sort.Strings(objects)

    return RoomView{
        Name:        room.Name,
        Description: room.Description,
        Exits:       exits,
        Others:      others,
        Mobiles:     mobiles,
        Objects:     objects,
        AreaName:    room.AreaName,
        AreaAuthor:  room.AreaAuthor,
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

// RespawnTick checks each area for respawn eligibility and respawns as needed
func (w *World) RespawnTick(defaultMinutes int, logger func(string)) {
    w.mu.Lock()
    defer w.mu.Unlock()

    now := time.Now()
    areasRespawned := 0
    var respawnLog []string

    // Group rooms by area
    areaRooms := make(map[string][]*Room)
    for _, room := range w.rooms {
        areaName := room.AreaName
        if areaName == "" {
            areaName = "Unknown"
        }
        areaRooms[areaName] = append(areaRooms[areaName], room)
    }

    // Check each area for respawn
    for areaName, rooms := range areaRooms {
        if len(rooms) == 0 {
            continue
        }

        // Get reset minutes for this area (from first room)
        resetMinutes := rooms[0].AreaResetMinutes
        if resetMinutes <= 0 {
            resetMinutes = defaultMinutes
        }

        // Check if area is due for respawn
        lastRespawn, hasLastRespawn := w.areaLastRespawn[areaName]
        if !hasLastRespawn {
            // First respawn - mark as just respawned
            w.areaLastRespawn[areaName] = now
            areasRespawned++
            respawnLog = append(respawnLog, fmt.Sprintf("%s (initial, next in %dm)", areaName, resetMinutes))
            continue
        }

        timeSinceRespawn := now.Sub(lastRespawn)
        respawnDue := timeSinceRespawn >= time.Duration(resetMinutes)*time.Minute

        if respawnDue {
            // Respawn this area
            for _, room := range rooms {
                // Clear existing mobs and objects
                room.Mobiles = make([]*Mobile, 0)
                room.Objects = make([]*Object, 0)

                // Re-instantiate from resets
                for _, reset := range room.MobileResets {
                    if proto, ok := w.mobiles[reset.MobVnum]; ok {
                        for i := 0; i < reset.Count; i++ {
                            mobCopy := *proto
                            room.Mobiles = append(room.Mobiles, &mobCopy)
                        }
                    }
                }
                for _, reset := range room.ObjectResets {
                    if proto, ok := w.objects[reset.ObjVnum]; ok {
                        for i := 0; i < reset.Count; i++ {
                            objCopy := *proto
                            room.Objects = append(room.Objects, &objCopy)
                        }
                    }
                }
            }

            w.areaLastRespawn[areaName] = now
            areasRespawned++
            respawnLog = append(respawnLog, fmt.Sprintf("%s (next in %dm)", areaName, resetMinutes))
        }
    }

    if logger == nil {
        return
    }

    if areasRespawned == 0 {
        logger(fmt.Sprintf("Respawn tick: no areas ready (%d monitored)", len(areaRooms)))
        return
    }

    logger(fmt.Sprintf("Respawn tick: %d/%d areas respawned - %s", areasRespawned, len(areaRooms), strings.Join(respawnLog, ", ")))
}
// FindMobInRoom searches for a mobile in the player's current room by keyword
func (w *World) FindMobInRoom(player *Player, keyword string) (*Mobile, bool) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	room, ok := w.rooms[player.Location]
	if !ok {
		return nil, false
	}

	keyword = strings.ToLower(keyword)
	for _, mob := range room.Mobiles {
		// Check if keyword matches any of the mob's keywords
		for _, mobKeyword := range mob.Keywords {
			if strings.ToLower(mobKeyword) == keyword {
				return mob, true
			}
		}
		// Also check if keyword appears in mob's short description
		if strings.Contains(strings.ToLower(mob.Short), keyword) {
			return mob, true
		}
	}

	return nil, false
}

// DamageMob deals damage to a mobile and handles death
func (w *World) DamageMob(player *Player, mob *Mobile, damage int) (died bool) {
	w.mu.Lock()
	defer w.mu.Unlock()

	mob.HP -= damage
	
	if mob.HP <= 0 {
		mob.HP = 0
		// Remove mob from room
		room, ok := w.rooms[player.Location]
		if ok {
			newMobiles := make([]*Mobile, 0, len(room.Mobiles)-1)
			for _, m := range room.Mobiles {
				if m != mob {
					newMobiles = append(newMobiles, m)
				}
			}
			room.Mobiles = newMobiles
		}
		return true
	}

	return false
}

// BroadcastCombatMessage sends a combat message to the player's room
func (w *World) BroadcastCombatMessage(player *Player, message string) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	for _, other := range w.players {
		if other.Location == player.Location && !strings.EqualFold(other.Name, player.Name) {
			other.Output.WriteLine(message)
		}
	}
}