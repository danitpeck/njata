package races

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// RaceJSON represents a race definition in JSON format
type RaceJSON struct {
	Name        string   `json:"name"`
	RaceID      int      `json:"race_id"`
	FlavorText  string   `json:"flavor_text"`
	Classes     int      `json:"classes"`
	StrPlus     int      `json:"str_plus"`
	DexPlus     int      `json:"dex_plus"`
	WisPlus     int      `json:"wis_plus"`
	IntPlus     int      `json:"int_plus"`
	ConPlus     int      `json:"con_plus"`
	ChaPlus     int      `json:"cha_plus"`
	LckPlus     int      `json:"lck_plus"`
	Hit         int      `json:"hit"`
	Mana        int      `json:"mana"`
	Affected    int      `json:"affected"`
	Resist      int      `json:"resist"`
	Suscept     int      `json:"suscept"`
	Language    int      `json:"language"`
	Align       int      `json:"align"`
	MinAlign    int      `json:"min_align"`
	MaxAlign    int      `json:"max_align"`
	ACPlus      int      `json:"ac_plus"`
	ExpMult     int      `json:"exp_mult"`
	Attacks     int      `json:"attacks"`
	Defenses    int      `json:"defenses"`
	Height      int      `json:"height"`
	Weight      int      `json:"weight"`
	HungerMod   int      `json:"hunger_mod"`
	ThirstMod   int      `json:"thirst_mod"`
	ManaRegen   int      `json:"mana_regen"`
	HPRegen     int      `json:"hp_regen"`
	MoveRegen   int      `json:"move_regen"`
	RaceRecall  int      `json:"race_recall"`
	WhereNames  []string `json:"where_names"`
}

// ConvertRacesToJSON converts all .race files from a directory to JSON format
func ConvertRacesToJSON(raceDir, outputDir string) error {
	entries, err := os.ReadDir(raceDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".race") {
			continue
		}

		filePath := filepath.Join(raceDir, entry.Name())
		race, err := parseRaceFile(filePath)
		if err != nil {
			fmt.Printf("Error parsing %s: %v\n", entry.Name(), err)
			continue
		}

		// Write JSON file
		outputFile := filepath.Join(outputDir, strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))+".json")
		data, err := json.MarshalIndent(race, "", "  ")
		if err != nil {
			fmt.Printf("Error marshaling JSON for %s: %v\n", entry.Name(), err)
			continue
		}

		if err := os.WriteFile(outputFile, data, 0644); err != nil {
			fmt.Printf("Error writing JSON file %s: %v\n", outputFile, err)
			continue
		}

		fmt.Printf("Converted: %s -> %s\n", entry.Name(), filepath.Base(outputFile))
	}

	return nil
}

func parseRaceFile(path string) (*RaceJSON, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	race := &RaceJSON{
		MinAlign:  -1000,
		MaxAlign:  1000,
		WhereNames: make([]string, 0),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		value := strings.Join(parts[1:], " ")
		value = strings.TrimSuffix(value, "~")
		value = strings.TrimSpace(value)

		switch key {
		case "Name":
			race.Name = value
		case "Race":
			if v, err := strconv.Atoi(value); err == nil {
				race.RaceID = v
			}
		case "Classes":
			if v, err := strconv.Atoi(value); err == nil {
				race.Classes = v
			}
		case "Str_Plus":
			if v, err := strconv.Atoi(value); err == nil {
				race.StrPlus = v
			}
		case "Dex_Plus":
			if v, err := strconv.Atoi(value); err == nil {
				race.DexPlus = v
			}
		case "Wis_Plus":
			if v, err := strconv.Atoi(value); err == nil {
				race.WisPlus = v
			}
		case "Int_Plus":
			if v, err := strconv.Atoi(value); err == nil {
				race.IntPlus = v
			}
		case "Con_Plus":
			if v, err := strconv.Atoi(value); err == nil {
				race.ConPlus = v
			}
		case "Cha_Plus":
			if v, err := strconv.Atoi(value); err == nil {
				race.ChaPlus = v
			}
		case "Lck_Plus":
			if v, err := strconv.Atoi(value); err == nil {
				race.LckPlus = v
			}
		case "Hit":
			if v, err := strconv.Atoi(value); err == nil {
				race.Hit = v
			}
		case "Mana":
			if v, err := strconv.Atoi(value); err == nil {
				race.Mana = v
			}
		case "Affected":
			if v, err := strconv.Atoi(value); err == nil {
				race.Affected = v
			}
		case "Resist":
			if v, err := strconv.Atoi(value); err == nil {
				race.Resist = v
			}
		case "Suscept":
			if v, err := strconv.Atoi(value); err == nil {
				race.Suscept = v
			}
		case "Language":
			if v, err := strconv.Atoi(value); err == nil {
				race.Language = v
			}
		case "Align":
			if v, err := strconv.Atoi(value); err == nil {
				race.Align = v
			}
		case "Min_Align":
			if v, err := strconv.Atoi(value); err == nil {
				race.MinAlign = v
			}
		case "Max_Align":
			if v, err := strconv.Atoi(value); err == nil {
				race.MaxAlign = v
			}
		case "AC_Plus":
			if v, err := strconv.Atoi(value); err == nil {
				race.ACPlus = v
			}
		case "Exp_Mult":
			if v, err := strconv.Atoi(value); err == nil {
				race.ExpMult = v
			}
		case "Attacks":
			if v, err := strconv.Atoi(value); err == nil {
				race.Attacks = v
			}
		case "Defenses":
			if v, err := strconv.Atoi(value); err == nil {
				race.Defenses = v
			}
		case "Height":
			if v, err := strconv.Atoi(value); err == nil {
				race.Height = v
			}
		case "Weight":
			if v, err := strconv.Atoi(value); err == nil {
				race.Weight = v
			}
		case "Hunger_Mod":
			if v, err := strconv.Atoi(value); err == nil {
				race.HungerMod = v
			}
		case "Thirst_mod":
			if v, err := strconv.Atoi(value); err == nil {
				race.ThirstMod = v
			}
		case "Mana_Regen":
			if v, err := strconv.Atoi(value); err == nil {
				race.ManaRegen = v
			}
		case "HP_Regen":
			if v, err := strconv.Atoi(value); err == nil {
				race.HPRegen = v
			}
		case "Move_Regen":
			if v, err := strconv.Atoi(value); err == nil {
				race.MoveRegen = v
			}
		case "Race_Recall":
			if v, err := strconv.Atoi(value); err == nil {
				race.RaceRecall = v
			}
		case "WhereName":
			race.WhereNames = append(race.WhereNames, value)
		}
	}

	return race, nil
}
