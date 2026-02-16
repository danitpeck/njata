package classes

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// ClassJSON represents a class definition in JSON format
type ClassJSON struct {
	Name           string            `json:"name"`
	ClassID        int               `json:"class_id"`
	FlavorText     string            `json:"flavor_text"`
	Races          int               `json:"races"`
	AttrPrime      int               `json:"attr_prime"`
	AttrSecond     int               `json:"attr_second"`
	AttrDeficient  int               `json:"attr_deficient"`
	Weapon         int               `json:"weapon"`
	Guild          int               `json:"guild"`
	Thac0          int               `json:"thac0"`
	Thac32         int               `json:"thac32"`
	Hpmin          int               `json:"hpmin"`
	Hpmax          int               `json:"hpmax"`
	Mana           int               `json:"mana"`
	Expbase        int               `json:"expbase"`
	Affected       int               `json:"affected"`
	Resist         int               `json:"resist"`
	Suscept        int               `json:"suscept"`
	Skills         map[string]SkillDef `json:"skills"`
}

// SkillDef represents a skill definition for a class
type SkillDef struct {
	Learned int `json:"learned"`
	Max     int `json:"max"`
}

// ConvertClassesToJSON converts all .class files from a directory to JSON format
func ConvertClassesToJSON(classDir, outputDir string) error {
	entries, err := os.ReadDir(classDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if !strings.HasSuffix(strings.ToLower(entry.Name()), ".class") {
			continue
		}

		filePath := filepath.Join(classDir, entry.Name())
		class, err := parseClassFile(filePath)
		if err != nil {
			fmt.Printf("Error parsing %s: %v\n", entry.Name(), err)
			continue
		}

		// Write JSON file
		outputFile := filepath.Join(outputDir, strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))+".json")
		data, err := json.MarshalIndent(class, "", "  ")
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

func parseClassFile(path string) (*ClassJSON, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	class := &ClassJSON{
		Skills: make(map[string]SkillDef),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Handle Skill lines specially
		if strings.HasPrefix(line, "Skill") {
			parts := strings.Split(line, "'")
			if len(parts) >= 3 {
				skillName := parts[1]
				skillVals := strings.Fields(parts[2])
				if len(skillVals) >= 3 {
					learned, _ := strconv.Atoi(skillVals[1])
					maxVal, _ := strconv.Atoi(skillVals[2])
					class.Skills[skillName] = SkillDef{Learned: learned, Max: maxVal}
				}
			}
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
			class.Name = value
		case "Class":
			if v, err := strconv.Atoi(value); err == nil {
				class.ClassID = v
			}
		case "Races":
			if v, err := strconv.Atoi(value); err == nil {
				class.Races = v
			}
		case "AttrPrime":
			if v, err := strconv.Atoi(value); err == nil {
				class.AttrPrime = v
			}
		case "AttrSecond":
			if v, err := strconv.Atoi(value); err == nil {
				class.AttrSecond = v
			}
		case "AttrDeficient":
			if v, err := strconv.Atoi(value); err == nil {
				class.AttrDeficient = v
			}
		case "Weapon":
			if v, err := strconv.Atoi(value); err == nil {
				class.Weapon = v
			}
		case "Guild":
			if v, err := strconv.Atoi(value); err == nil {
				class.Guild = v
			}
		case "Thac0":
			if v, err := strconv.Atoi(value); err == nil {
				class.Thac0 = v
			}
		case "Thac32":
			if v, err := strconv.Atoi(value); err == nil {
				class.Thac32 = v
			}
		case "Hpmin":
			if v, err := strconv.Atoi(value); err == nil {
				class.Hpmin = v
			}
		case "Hpmax":
			if v, err := strconv.Atoi(value); err == nil {
				class.Hpmax = v
			}
		case "Mana":
			if v, err := strconv.Atoi(value); err == nil {
				class.Mana = v
			}
		case "Expbase":
			if v, err := strconv.Atoi(value); err == nil {
				class.Expbase = v
			}
		case "Affected":
			if v, err := strconv.Atoi(value); err == nil {
				class.Affected = v
			}
		case "Resist":
			if v, err := strconv.Atoi(value); err == nil {
				class.Resist = v
			}
		case "Suscept":
			if v, err := strconv.Atoi(value); err == nil {
				class.Suscept = v
			}
		}
	}

	return class, nil
}
