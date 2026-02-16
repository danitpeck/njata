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
	Name    string            `json:"name"`
	Author  string            `json:"author"`
	Rooms   map[string]RoomJSON `json:"rooms"` // keyed by vnum as string
}

// RoomJSON is the JSON-serializable version of Room
type RoomJSON struct {
	Vnum       int               `json:"vnum"`
	Name       string            `json:"name"`
	Description string           `json:"description"`
	Sector     string            `json:"sector"`
	Flags      map[string]bool   `json:"flags"`
	Exits      map[string]int    `json:"exits"`
	ExDescs    map[string]string `json:"exdescs"`
	AreaName   string            `json:"area_name"`
	AreaAuthor string            `json:"area_author"`
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
		rooms, areaName, areaAuthor, err := parseAreFile(filePath)
		if err != nil {
			fmt.Printf("Error parsing %s: %v\n", entry.Name(), err)
			continue
		}

		// Convert to JSON format
		roomsJSON := make(map[string]RoomJSON)
		for vnum, room := range rooms {
			roomsJSON[fmt.Sprintf("%d", vnum)] = RoomJSON{
				Vnum:        vnum,
				Name:        room.Name,
				Description: room.Description,
				Sector:      room.Sector,
				Flags:       room.Flags,
				Exits:       room.Exits,
				ExDescs:     room.ExDescs,
				AreaName:    areaName,
				AreaAuthor:  areaAuthor,
			}
		}

		areaJSON := AreaJSON{
			Name:   areaName,
			Author: areaAuthor,
			Rooms:  roomsJSON,
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

		fmt.Printf("Converted: %s -> %s (%d rooms)\n", entry.Name(), filepath.Base(outputFile), len(rooms))
	}

	return nil
}

// Internal .are file parser for conversion
type areRoom struct {
	Name        string
	Description string
	Sector      string
	Flags       map[string]bool
	Exits       map[string]int
	ExDescs     map[string]string
}

func parseAreFile(path string) (map[int]areRoom, string, string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, "", "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	rooms := make(map[int]areRoom)
	var areaName, areaAuthor string
	var currentVnum int
	var current *areRoom

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if line == "#AREADATA" {
			name, author := parseAreaData(scanner)
			areaName = name
			areaAuthor = author
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
		}
	}

	return rooms, areaName, areaAuthor, nil
}

func parseAreaData(scanner *bufio.Scanner) (string, string) {
	var name, author string
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
	}
	return name, author
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
