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

func LoadRoomsFromDir(path string) (map[int]*game.Room, map[int]*game.Mobile, map[int]*game.Object, int, error) {
    entries, err := os.ReadDir(path)
    if err != nil {
        return nil, nil, nil, 0, err
    }

    rooms := map[int]*game.Room{}
    mobiles := map[int]*game.Mobile{}
    objects := map[int]*game.Object{}

    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        if !strings.HasSuffix(strings.ToLower(entry.Name()), ".json") {
            continue
        }

        filePath := filepath.Join(path, entry.Name())
        parsed, mobs, objs, err := parseAreasFromJSON(filePath)
        if err != nil {
            return nil, nil, nil, 0, fmt.Errorf("%s: %w", entry.Name(), err)
        }

        for vnum, room := range parsed {
            rooms[vnum] = room
        }
        for vnum, mob := range mobs {
            mobiles[vnum] = mob
        }
        for vnum, obj := range objs {
            objects[vnum] = obj
        }
    }

    // Process resets to instantiate mobs/objects in rooms
    for _, room := range rooms {
        for _, reset := range room.MobileResets {
            if proto, ok := mobiles[reset.MobVnum]; ok {
                for i := 0; i < reset.Count; i++ {
                    // Create a copy of the prototype
                    mobCopy := *proto
                    room.Mobiles = append(room.Mobiles, &mobCopy)
                }
            }
        }
        for _, reset := range room.ObjectResets {
            if proto, ok := objects[reset.ObjVnum]; ok {
                for i := 0; i < reset.Count; i++ {
                    // Create a copy of the prototype
                    objCopy := *proto
                    room.Objects = append(room.Objects, &objCopy)
                }
            }
        }
    }

    start := findLowestVnum(rooms)
    return rooms, mobiles, objects, start, nil
}

func parseAreasFromJSON(path string) (map[int]*game.Room, map[int]*game.Mobile, map[int]*game.Object, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        return nil, nil, nil, err
    }

    var areaJSON struct {
        Name          string `json:"name"`
        Author        string `json:"author"`
        ResetMinutes  int    `json:"reset_minutes"`
        Rooms   map[string]struct {
            Vnum          int               `json:"vnum"`
            Name          string            `json:"name"`
            Description   string            `json:"description"`
            Sector        string            `json:"sector"`
            Flags         map[string]bool   `json:"flags"`
            Exits         map[string]int    `json:"exits"`
            ExDescs       map[string]string `json:"exdescs"`
            AreaName      string            `json:"area_name"`
            AreaAuthor    string            `json:"area_author"`
            MobileResets  []struct {
                Vnum  int `json:"vnum"`
                Count int `json:"count"`
                Room  int `json:"room"`
            } `json:"mobile_resets"`
            ObjectResets  []struct {
                Vnum  int `json:"vnum"`
                Count int `json:"count"`
                Room  int `json:"room"`
            } `json:"object_resets"`
        } `json:"rooms"`
        Mobiles map[string]struct {
            Vnum       int    `json:"vnum"`
            Keywords   []string `json:"keywords"`
            Short      string `json:"short"`
            Long       string `json:"long"`
            Race       string `json:"race"`
            Class      string `json:"class"`
            Position   string `json:"position"`
            Gender     string `json:"gender"`
            Level      int    `json:"level"`
            MaxHP      int    `json:"max_hp"`
            HP         int    `json:"hp"`
            Mana       int    `json:"mana"`
            MaxMana    int    `json:"max_mana"`
            Attributes [7]int `json:"attributes"`
        } `json:"mobiles"`
        Objects map[string]struct {
            Vnum     int               `json:"vnum"`
            Keywords []string          `json:"keywords"`
            Type     string            `json:"type"`
            Short    string            `json:"short"`
            Long     string            `json:"long"`
            Weight   int               `json:"weight"`
            Value    int               `json:"value"`
            Flags    map[string]bool   `json:"flags"`
        } `json:"objects"`
    }

    if err := json.Unmarshal(data, &areaJSON); err != nil {
        return nil, nil, nil, err
    }

    rooms := make(map[int]*game.Room)
    for _, roomJSON := range areaJSON.Rooms {
        // Parse resets
        mobResets := make([]game.Reset, len(roomJSON.MobileResets))
        for i, mr := range roomJSON.MobileResets {
            mobResets[i] = game.Reset{
                MobVnum: mr.Vnum,
                Count:   mr.Count,
                Room:    mr.Room,
            }
        }
        objResets := make([]game.Reset, len(roomJSON.ObjectResets))
        for i, or := range roomJSON.ObjectResets {
            objResets[i] = game.Reset{
                ObjVnum: or.Vnum,
                Count:   or.Count,
                Room:    or.Room,
            }
        }

        room := &game.Room{
            Vnum:         roomJSON.Vnum,
            Name:         roomJSON.Name,
            Description:  roomJSON.Description,
            Sector:       roomJSON.Sector,
            Flags:        roomJSON.Flags,
            Exits:        roomJSON.Exits,
            ExDescs:      roomJSON.ExDescs,
            AreaName:     areaJSON.Name,
            AreaAuthor:   areaJSON.Author,
            AreaResetMinutes: areaJSON.ResetMinutes,
            Mobiles:      make([]*game.Mobile, 0),
            Objects:      make([]*game.Object, 0),
            MobileResets: mobResets,
            ObjectResets: objResets,
        }
        if room.Vnum > 0 {
            rooms[room.Vnum] = room
        }
    }

    mobiles := make(map[int]*game.Mobile)
    for _, mobJSON := range areaJSON.Mobiles {
        mob := &game.Mobile{
            Vnum:       mobJSON.Vnum,
            Keywords:   mobJSON.Keywords,
            Short:      mobJSON.Short,
            Long:       mobJSON.Long,
            Race:       mobJSON.Race,
            Class:      mobJSON.Class,
            Position:   mobJSON.Position,
            Gender:     mobJSON.Gender,
            Level:      mobJSON.Level,
            MaxHP:      mobJSON.MaxHP,
            HP:         mobJSON.HP,
            Mana:       mobJSON.Mana,
            MaxMana:    mobJSON.MaxMana,
            Attributes: mobJSON.Attributes,
        }
        if mob.Vnum > 0 {
            mobiles[mob.Vnum] = mob
        }
    }

    objects := make(map[int]*game.Object)
    for _, objJSON := range areaJSON.Objects {
        obj := &game.Object{
            Vnum:     objJSON.Vnum,
            Keywords: objJSON.Keywords,
            Type:     objJSON.Type,
            Short:    objJSON.Short,
            Long:     objJSON.Long,
            Weight:   objJSON.Weight,
            Value:    objJSON.Value,
            Flags:    objJSON.Flags,
        }
        if obj.Vnum > 0 {
            objects[obj.Vnum] = obj
        }
    }

    return rooms, mobiles, objects, nil
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
