package area

import (
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "strings"

    "njata/internal/game"
)

func LoadRoomsFromDir(path string) (map[int]*game.Room, int, error) {
    entries, err := os.ReadDir(path)
    if err != nil {
        return nil, 0, err
    }

    rooms := map[int]*game.Room{}
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        if !strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
            continue
        }

        filePath := filepath.Join(path, entry.Name())
        parsed, err := parseRoomsFromJSON(filePath)
        if err != nil {
            return nil, 0, fmt.Errorf("%s: %w", entry.Name(), err)
        }

        for vnum, room := range parsed {
            rooms[vnum] = room
        }
    }

    start := findLowestVnum(rooms)
    return rooms, start, nil
}

func parseRoomsFromJSON(path string) (map[int]*game.Room, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }

    var areaJSON struct {
        Name   string `json:"name"`
        Author string `json:"author"`
        Rooms  map[string]struct {
            Vnum        int               `json:"vnum"`
            Name        string            `json:"name"`
            Description string            `json:"description"`
            Sector      string            `json:"sector"`
            Flags       map[string]bool   `json:"flags"`
            Exits       map[string]int    `json:"exits"`
            ExDescs     map[string]string `json:"exdescs"`
            AreaName    string            `json:"area_name"`
            AreaAuthor  string            `json:"area_author"`
        } `json:"rooms"`
    }

    if err := json.Unmarshal(data, &areaJSON); err != nil {
        return nil, err
    }

    rooms := make(map[int]*game.Room)
    for _, roomJSON := range areaJSON.Rooms {
        room := &game.Room{
            Vnum:        roomJSON.Vnum,
            Name:        roomJSON.Name,
            Description: roomJSON.Description,
            Sector:      roomJSON.Sector,
            Flags:       roomJSON.Flags,
            Exits:       roomJSON.Exits,
            ExDescs:     roomJSON.ExDescs,
            AreaName:    areaJSON.Name, // Use area-level name
            AreaAuthor:  areaJSON.Author, // Use area-level author
        }
        if room.Vnum > 0 {
            rooms[room.Vnum] = room
        }
    }

    return rooms, nil
}

func findLowestVnum(rooms map[int]*game.Room) int {
    if len(rooms) == 0 {
        return 0
    }

    vnums := make([]int, 0, len(rooms))
    for vnum := range rooms {
        vnums = append(vnums, vnum)
    }
    sort.Ints(vnums)
    return vnums[0]
}
