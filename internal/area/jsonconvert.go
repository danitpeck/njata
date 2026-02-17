package area

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// AreaJSON represents an area file in JSON format
type AreaJSON struct {
	Name          string               `json:"name"`
	Author        string               `json:"author"`
	ResetMinutes  int                  `json:"reset_minutes"`
	Rooms         map[string]RoomJSON  `json:"rooms"`   // keyed by vnum as string
	Mobiles       map[string]MobileJSON `json:"mobiles"` // keyed by vnum as string
	Objects       map[string]ObjectJSON `json:"objects"` // keyed by vnum as string
}

// RoomJSON is the JSON-serializable version of Room
type RoomJSON struct {
	Vnum           int               `json:"vnum"`
	Name           string            `json:"name"`
	Description    string            `json:"description"`
	Sector         string            `json:"sector"`
	Flags          map[string]bool   `json:"flags"`
	Exits          map[string]int    `json:"exits"`
	ExDescs        map[string]string `json:"exdescs"`
	AreaName       string            `json:"area_name"`
	AreaAuthor     string            `json:"area_author"`
	MobileResets   []ResetJSON       `json:"mobile_resets"`
	ObjectResets   []ResetJSON       `json:"object_resets"`
}

// ResetJSON represents a Reset command
type ResetJSON struct {
	Vnum  int `json:"vnum"`
	Count int `json:"count"`
	Room  int `json:"room"`
}

// MobileJSON is the JSON-serializable version of Mobile
type MobileJSON struct {
	Vnum       int      `json:"vnum"`
	Keywords   []string `json:"keywords"`
	Short      string   `json:"short"`
	Long       string   `json:"long"`
	Race       string   `json:"race"`
	Class      string   `json:"class"`
	Position   string   `json:"position"`
	Gender     string   `json:"gender"`
	Level      int      `json:"level"`
	MaxHP      int      `json:"max_hp"`
	HP         int      `json:"hp"`
	Mana       int      `json:"mana"`
	MaxMana    int      `json:"max_mana"`
	Attributes [7]int   `json:"attributes"`
	RoomVnum   int      `json:"room_vnum"` // which room this mobile is in
}

// ObjectJSON is the JSON-serializable version of Object
type ObjectJSON struct {
	Vnum     int               `json:"vnum"`
	Keywords []string          `json:"keywords"`
	Type     string            `json:"type"`
	Short    string            `json:"short"`
	Long     string            `json:"long"`
	Weight   int               `json:"weight"`
	Value    int               `json:"value"`
	Flags    map[string]bool   `json:"flags"`
	RoomVnum int               `json:"room_vnum"` // which room this object is in
}

// ConvertAreToJSON converts all .are files in a directory to JSON format
func ConvertAreToJSON(areDir, outputDir string) error {
	entries, err := os.ReadDir(areDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".are") {
			continue
		}

		filePath := filepath.Join(areDir, entry.Name())
		rooms, mobiles, objects, areaName, areaAuthor, resetMinutes, err := parseAreFile(filePath)
		if err != nil {
			fmt.Printf("Error parsing %s: %v\n", entry.Name(), err)
			continue
		}

		// Convert to JSON format
		roomsJSON := make(map[string]RoomJSON)
		for vnum, room := range rooms {
			mobileResetsJSON := make([]ResetJSON, len(room.MobileResets))
			for i, mr := range room.MobileResets {
				mobileResetsJSON[i] = ResetJSON{Vnum: mr.Vnum, Count: mr.Count, Room: mr.Room}
			}
			objectResetsJSON := make([]ResetJSON, len(room.ObjectResets))
			for i, or := range room.ObjectResets {
				objectResetsJSON[i] = ResetJSON{Vnum: or.Vnum, Count: or.Count, Room: or.Room}
			}
			
			roomsJSON[fmt.Sprintf("%d", vnum)] = RoomJSON{
				Vnum:         vnum,
				Name:         room.Name,
				Description:  room.Description,
				Sector:       room.Sector,
				Flags:        room.Flags,
				Exits:        room.Exits,
				ExDescs:      room.ExDescs,
				AreaName:     areaName,
				AreaAuthor:   areaAuthor,
				MobileResets: mobileResetsJSON,
				ObjectResets: objectResetsJSON,
			}
		}

		mobilesJSON := make(map[string]MobileJSON)
		for vnum, mob := range mobiles {
			mobilesJSON[fmt.Sprintf("%d", vnum)] = mob
		}

		objectsJSON := make(map[string]ObjectJSON)
		for vnum, obj := range objects {
			objectsJSON[fmt.Sprintf("%d", vnum)] = obj
		}

		areaJSON := AreaJSON{
			Name:         areaName,
			Author:       areaAuthor,
			ResetMinutes: resetMinutes,
			Rooms:        roomsJSON,
			Mobiles:      mobilesJSON,
			Objects:      objectsJSON,
		}

		// Write JSON file
		outputFile := filepath.Join(outputDir, strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))+".json")
		data, err := json.MarshalIndent(areaJSON, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling JSON for %s: %v\n", entry.Name(), err)
			continue
		}

		if err := os.WriteFile(outputFile, data, 0644); err != nil {
			fmt.Printf("Error writing JSON file %s: %v\n", outputFile, err)
			continue
		}

		fmt.Printf("Converted: %s -> %s (%d rooms, %d mobs, %d objs)\n", entry.Name(), filepath.Base(outputFile), len(rooms), len(mobiles), len(objects))
	}

	return nil
}

// Internal .are file parser for conversion
type areRoom struct {
	Name          string
	Description   string
	Sector        string
	Flags         map[string]bool
	Exits         map[string]int
	ExDescs       map[string]string
	MobileResets  []struct {
		Vnum  int
		Count int
		Room  int
	}
	ObjectResets  []struct {
		Vnum  int
		Count int
		Room  int
	}
}

func parseAreFile(path string) (map[int]areRoom, map[int]MobileJSON, map[int]ObjectJSON, string, string, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, nil, "", "", 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	rooms := make(map[int]areRoom)
	mobiles := make(map[int]MobileJSON)
	objects := make(map[int]ObjectJSON)
	var areaName, areaAuthor string
	var resetMinutes int
	var currentVnum int
	var current *areRoom

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if line == "#AREADATA" {
			name, author, resetFreq := parseAreaData(scanner)
			areaName = name
			areaAuthor = author
			resetMinutes = resetFreq
			continue
		}

		if line == "#MOBILE" {
			mob := parseMobileSection(scanner)
			if mob.Vnum > 0 {
				mobiles[mob.Vnum] = mob
			}
			continue
		}

		if line == "#OBJECT" {
			obj := parseObjectSection(scanner)
			if obj.Vnum > 0 {
				objects[obj.Vnum] = obj
			}
			continue
		}

		if line == "#ROOM" {
			current = &areRoom{
				Flags:   make(map[string]bool),
				Exits:   make(map[string]int),
				ExDescs: make(map[string]string),
			}
			continue
		}

		if current == nil {
			continue
		}

		switch {
		case line == "#ENDROOM":
			if currentVnum > 0 {
				rooms[currentVnum] = *current
			}
			current = nil
		case strings.HasPrefix(line, "Vnum"):
			vnum, _ := parseIntField(line)
			currentVnum = vnum
		case strings.HasPrefix(line, "Name"):
			name, _ := readTildeTextAre(fieldRemainder(line, "Name"), scanner)
			current.Name = name
		case strings.HasPrefix(line, "Sector"):
			sector, _ := readTildeTextAre(fieldRemainder(line, "Sector"), scanner)
			current.Sector = strings.ToLower(strings.TrimSpace(sector))
		case strings.HasPrefix(line, "Flags"):
			flags, _ := readTildeTextAre(fieldRemainder(line, "Flags"), scanner)
			for _, flag := range strings.Fields(strings.ToLower(flags)) {
				if flag != "" {
					current.Flags[flag] = true
				}
			}
		case strings.HasPrefix(line, "Desc"):
			desc, _ := readTildeTextAre(fieldRemainder(line, "Desc"), scanner)
			current.Description = desc
		case line == "#EXIT":
			direction, toRoom := readExitAre(scanner)
			if direction != "" && toRoom > 0 {
				current.Exits[direction] = toRoom
			}
		case line == "#EXDESC":
			keys, desc := readExDescAre(scanner)
			for _, key := range keys {
				current.ExDescs[key] = desc
			}
		case strings.HasPrefix(line, "Reset"):
			parts := strings.Fields(line)
			if len(parts) >= 6 {
				resetType := parts[1] // "M" or "O"
				vnum, _ := strconv.Atoi(parts[3])
				count, _ := strconv.Atoi(parts[4])
				room, _ := strconv.Atoi(parts[5])
				
				if resetType == "M" {
					current.MobileResets = append(current.MobileResets, struct {
						Vnum  int
						Count int
						Room  int
					}{vnum, count, room})
				} else if resetType == "O" {
					current.ObjectResets = append(current.ObjectResets, struct {
						Vnum  int
						Count int
						Room  int
					}{vnum, count, room})
				}
			}
		}
	}

	return rooms, mobiles, objects, areaName, areaAuthor, resetMinutes, nil
}

func parseAreaData(scanner *bufio.Scanner) (string, string, int) {
	var name, author string
	var resetFreq int
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "#ENDAREADATA" {
			break
		}
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "Name") {
			name = strings.TrimSuffix(strings.TrimSpace(fieldRemainder(line, "Name")), "~")
		}
		if strings.HasPrefix(line, "Author") {
			author = strings.TrimSuffix(strings.TrimSpace(fieldRemainder(line, "Author")), "~")
		}
		if strings.HasPrefix(line, "ResetFreq") {
			fields := strings.Fields(line)
			if len(fields) > 1 {
				if value, err := strconv.Atoi(fields[1]); err == nil {
					resetFreq = value
				}
			}
		}
	}
	return name, author, resetFreq
}

func readExDescAre(scanner *bufio.Scanner) ([]string, string) {
	var keys []string
	var desc string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "#ENDEXDESC" {
			return keys, desc
		}
		if strings.HasPrefix(line, "ExDescKey") {
			value, _ := readTildeTextAre(fieldRemainder(line, "ExDescKey"), scanner)
			for _, field := range strings.Fields(strings.ToLower(value)) {
				if field != "" {
					keys = append(keys, field)
				}
			}
		}
		if strings.HasPrefix(line, "ExDesc") {
			value, _ := readTildeTextAre(fieldRemainder(line, "ExDesc"), scanner)
			desc = value
		}
	}
	return keys, desc
}

func readExitAre(scanner *bufio.Scanner) (string, int) {
	var direction string
	var toRoom int
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "#ENDEXIT" {
			return direction, toRoom
		}
		if strings.HasPrefix(line, "Direction") {
			value, _ := readTildeTextAre(fieldRemainder(line, "Direction"), scanner)
			direction = strings.ToLower(strings.TrimSpace(value))
		}
		if strings.HasPrefix(line, "ToRoom") {
			if parsed, err := strconv.Atoi(strings.Fields(line)[1]); err == nil {
				toRoom = parsed
			}
		}
	}
	return direction, toRoom
}

func parseIntField(line string) (int, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return 0, fmt.Errorf("invalid numeric field: %s", line)
	}
	return strconv.Atoi(fields[1])
}

func fieldRemainder(line, prefix string) string {
	return strings.TrimSpace(strings.TrimPrefix(line, prefix))
}

func readTildeTextAre(initial string, scanner *bufio.Scanner) (string, error) {
	if idx := strings.Index(initial, "~"); idx >= 0 {
		return strings.TrimRight(strings.TrimRight(initial[:idx], "\r"), " "), nil
	}

	var builder strings.Builder
	trimmed := strings.TrimRight(strings.TrimRight(initial, "\r"), " ")
	if trimmed != "" {
		if builder.Len() > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(trimmed)
	}

	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")
		if idx := strings.Index(line, "~"); idx >= 0 {
			segment := strings.TrimRight(line[:idx], " ")
			if segment != "" {
				if builder.Len() > 0 {
					builder.WriteString("\n")
				}
				builder.WriteString(segment)
			}
			return builder.String(), nil
		}
		if builder.Len() > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(strings.TrimRight(line, " "))
	}

	return builder.String(), nil
}
func parseMobileSection(scanner *bufio.Scanner) MobileJSON {
	mob := MobileJSON{
		Attributes: [7]int{},
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Only break on explicit mobile end markers or new sections
		if line == "#ENDMOBILE" || strings.HasPrefix(line, "#") {
			break
		}
		// Skip standalone tildes (they terminate multi-line fields, not the mobile)
		if line == "~" {
			continue
		}
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "Vnum"):
			vnum, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Vnum")))
			mob.Vnum = vnum
		case strings.HasPrefix(line, "Keywords"):
			keywords := strings.TrimSpace(strings.TrimPrefix(line, "Keywords"))
			keywords = strings.TrimSuffix(keywords, "~")
			mob.Keywords = strings.Fields(keywords)
		case strings.HasPrefix(line, "ShortDesc") || strings.HasPrefix(line, "Short"):
			desc, _ := readTildeTextAre(strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "ShortDesc"), "Short")), scanner)
			mob.Short = desc
		case strings.HasPrefix(line, "LongDesc") || strings.HasPrefix(line, "Long"):
			desc, _ := readTildeTextAre(strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "LongDesc"), "Long")), scanner)
			mob.Long = desc
		case strings.HasPrefix(line, "Race"):
			race := strings.TrimSuffix(strings.TrimSpace(strings.TrimPrefix(line, "Race")), "~")
			mob.Race = race
		case strings.HasPrefix(line, "Class"):
			class := strings.TrimSuffix(strings.TrimSpace(strings.TrimPrefix(line, "Class")), "~")
			mob.Class = class
		case strings.HasPrefix(line, "Gender"):
			gender := strings.TrimSuffix(strings.TrimSpace(strings.TrimPrefix(line, "Gender")), "~")
			mob.Gender = gender
		case strings.HasPrefix(line, "Level"):
			level, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Level")))
			mob.Level = level
		case strings.HasPrefix(line, "Hitpoints"):
			hp, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Hitpoints")))
			mob.HP = hp
			mob.MaxHP = hp
		case strings.HasPrefix(line, "Mana"):
			mana, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Mana")))
			mob.Mana = mana
			mob.MaxMana = mana
		case strings.HasPrefix(line, "Strength"):
			str, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Strength")))
			mob.Attributes[0] = str
		case strings.HasPrefix(line, "Dexterity"):
			dex, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Dexterity")))
			mob.Attributes[1] = dex
		case strings.HasPrefix(line, "Constitution"):
			con, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Constitution")))
			mob.Attributes[2] = con
		case strings.HasPrefix(line, "Intelligence"):
			intel, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Intelligence")))
			mob.Attributes[3] = intel
		case strings.HasPrefix(line, "Wisdom"):
			wis, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Wisdom")))
			mob.Attributes[4] = wis
		case strings.HasPrefix(line, "Charisma"):
			cha, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Charisma")))
			mob.Attributes[5] = cha
		case strings.HasPrefix(line, "Luck"):
			luck, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Luck")))
			mob.Attributes[6] = luck
		case strings.HasPrefix(line, "Attribs"):
			attribs := strings.TrimSpace(strings.TrimPrefix(line, "Attribs"))
			attribs = strings.TrimSuffix(attribs, "~")
			values := strings.Fields(attribs)
			for i := 0; i < len(values) && i < 7; i++ {
				val, _ := strconv.Atoi(values[i])
				mob.Attributes[i] = val
			}
		case strings.HasPrefix(line, "Stats2"):
			stats := strings.TrimSpace(strings.TrimPrefix(line, "Stats2"))
			stats = strings.TrimSuffix(stats, "~")
			values := strings.Fields(stats)
			if len(values) > 0 {
				maxhp, _ := strconv.Atoi(values[0])
				mob.MaxHP = maxhp
				mob.HP = maxhp
			}
			if len(values) > 2 {
				mana, _ := strconv.Atoi(values[2])
				mob.Mana = mana
				mob.MaxMana = mana
			}
		case strings.HasPrefix(line, "Position"):
			pos := strings.TrimSuffix(strings.TrimSpace(strings.TrimPrefix(line, "Position")), "~")
			mob.Position = pos
		}
	}

	return mob
}

func parseObjectSection(scanner *bufio.Scanner) ObjectJSON {
	obj := ObjectJSON{
		Flags: make(map[string]bool),
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Only break on explicit object end markers or new sections
		if line == "#ENDOBJECT" || strings.HasPrefix(line, "#") {
			break
		}
		// Skip standalone tildes (they terminate multi-line fields, not the object)
		if line == "~" {
			continue
		}
		if line == "" {
			continue
		}

		switch {
		case strings.HasPrefix(line, "Vnum"):
			vnum, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Vnum")))
			obj.Vnum = vnum
		case strings.HasPrefix(line, "Keywords"):
			keywords := strings.TrimSpace(strings.TrimPrefix(line, "Keywords"))
			keywords = strings.TrimSuffix(keywords, "~")
			obj.Keywords = strings.Fields(keywords)
		case strings.HasPrefix(line, "ShortDesc") || strings.HasPrefix(line, "Short"):
			desc, _ := readTildeTextAre(strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "ShortDesc"), "Short")), scanner)
			obj.Short = desc
		case strings.HasPrefix(line, "LongDesc") || strings.HasPrefix(line, "Long"):
			desc, _ := readTildeTextAre(strings.TrimSpace(strings.TrimPrefix(strings.TrimPrefix(line, "LongDesc"), "Long")), scanner)
			obj.Long = desc
		case strings.HasPrefix(line, "Type"):
			objType := strings.TrimSuffix(strings.TrimSpace(strings.TrimPrefix(line, "Type")), "~")
			obj.Type = objType
		case strings.HasPrefix(line, "Weight"):
			weight, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Weight")))
			obj.Weight = weight
		case strings.HasPrefix(line, "Value"):
			value, _ := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "Value")))
			obj.Value = value
		}
	}

	return obj
}