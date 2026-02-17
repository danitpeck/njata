package classes

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

var (
	classesByID   map[int]*ClassJSON
	classesByName map[string]*ClassJSON
	classList    []*ClassJSON // ordered list for menu display
)

// Load reads all class JSON files from the classes directory and indexes them
func Load(classesDir string) error {
	classesByID = make(map[int]*ClassJSON)
	classesByName = make(map[string]*ClassJSON)
	classList = make([]*ClassJSON, 0)

	entries, err := os.ReadDir(classesDir)
	if err != nil {
		return fmt.Errorf("failed to read classes directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(classesDir, entry.Name())
		data, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("failed to read class file %s: %w", entry.Name(), err)
		}

		var class ClassJSON
		if err := json.Unmarshal(data, &class); err != nil {
			return fmt.Errorf("failed to parse class file %s: %w", entry.Name(), err)
		}

		classesByID[class.ClassID] = &class
		classesByName[class.Name] = &class
		classList = append(classList, &class)
	}

	// Sort by class ID for consistent menu ordering
	sort.Slice(classList, func(i, j int) bool {
		return classList[i].ClassID < classList[j].ClassID
	})

	return nil
}

// GetByID returns a class by its ID
func GetByID(id int) *ClassJSON {
	return classesByID[id]
}

// GetByName returns a class by name
func GetByName(name string) *ClassJSON {
	return classesByName[name]
}

// List returns all classes in order
func List() []*ClassJSON {
	return classList
}

// Count returns the number of loaded classes
func Count() int {
	return len(classList)
}

// MenuString returns a formatted list of classes for the character creation menu
func MenuString() string {
	var result string
	for i, class := range classList {
		result += fmt.Sprintf("  %2d) %s\n", i+1, class.Name)
	}
	return result
}

// GetByMenuChoice returns the class at the given menu position (1-indexed)
func GetByMenuChoice(choice int) *ClassJSON {
	if choice < 1 || choice > len(classList) {
		return nil
	}
	return classList[choice-1]
}

// MVPClasses returns only the two MVP classes: Scholar and Warrior
func MVPClasses() []*ClassJSON {
	var mvp []*ClassJSON
	for _, class := range classList {
		if class.Name == "Scholar" || class.Name == "Warrior" {
			mvp = append(mvp, class)
		}
	}
	return mvp
}

// MVPMenuString returns a formatted menu for MVP classes only (Scholar and Warrior)
func MVPMenuString() string {
	var result string
	for i, class := range MVPClasses() {
		result += fmt.Sprintf("  %d) %s\n", i+1, class.Name)
	}
	return result
}

// GetMVPByMenuChoice returns the MVP class at the given menu position (1-indexed)
func GetMVPByMenuChoice(choice int) *ClassJSON {
	mvp := MVPClasses()
	if choice < 1 || choice > len(mvp) {
		return nil
	}
	return mvp[choice-1]
}
