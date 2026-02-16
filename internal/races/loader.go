package races

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

var (
	racesByID   map[int]*RaceJSON
	racesByName map[string]*RaceJSON
	raceList    []*RaceJSON // ordered list for menu display
)

// Load reads all race JSON files from the races directory and indexes them
func Load(racesDir string) error {
	racesByID = make(map[int]*RaceJSON)
	racesByName = make(map[string]*RaceJSON)
	raceList = make([]*RaceJSON, 0)

	entries, err := os.ReadDir(racesDir)
	if err != nil {
		return fmt.Errorf("failed to read races directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(racesDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read race file %s: %w", entry.Name(), err)
		}

		var race RaceJSON
		if err := json.Unmarshal(data, &race); err != nil {
			return fmt.Errorf("failed to parse race file %s: %w", entry.Name(), err)
		}

		racesByID[race.RaceID] = &race
		racesByName[race.Name] = &race
		raceList = append(raceList, &race)
	}

	// Sort by race ID for consistent menu ordering
	sort.Slice(raceList, func(i, j int) bool {
		return raceList[i].RaceID < raceList[j].RaceID
	})

	return nil
}

// GetByID returns a race by its ID
func GetByID(id int) *RaceJSON {
	return racesByID[id]
}

// GetByName returns a race by name
func GetByName(name string) *RaceJSON {
	return racesByName[name]
}

// List returns all races in order
func List() []*RaceJSON {
	return raceList
}

// Count returns the number of loaded races
func Count() int {
	return len(raceList)
}

// MenuString returns a formatted list of races for the character creation menu
func MenuString() string {
	var result string
	for i, race := range raceList {
		result += fmt.Sprintf("  %2d) %s\n", i+1, race.Name)
	}
	return result
}

// GetByMenuChoice returns the race at the given menu position (1-indexed)
func GetByMenuChoice(choice int) *RaceJSON {
	if choice < 1 || choice > len(raceList) {
		return nil
	}
	return raceList[choice-1]
}
