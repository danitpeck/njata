package kits

import "fmt"

// StarterKit represents a starting character configuration
type StarterKit struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	FlavorText  string `json:"flavor_text"`

	// Starting stats
	HP   int `json:"hp"`
	Mana int `json:"mana"`

	// Starting skills (skill_id -> initial proficiency %)
	StartingSkills map[int]int `json:"starting_skills"`

	// Starting equipment item types
	StartingEquipment []string `json:"starting_equipment"`
}

// Predefined starter kits
var (
	ScholarKit = &StarterKit{
		ID:                1,
		Name:              "Scholar's Kit",
		Description:       "A leather-bound libram of magical knowledge, scholar's robes, and study supplies. Begin unlocking magical secrets through research.",
		FlavorText:        "Scholars are seekers of knowledge who unlock magical power by studying\nenchanted items found throughout the world. Your growth is driven by\ncuriosity and exploration.",
		HP:                100,
		Mana:              150,
		StartingSkills:    map[int]int{},
		StartingEquipment: []string{"robes", "spellbook"},
	}

	WarriorKit = &StarterKit{
		ID:          2,
		Name:        "Warrior's Kit",
		Description: "Leather armor, a basic sword, and training manual. Begin your path to combat mastery.",
		FlavorText:  "Warriors are masters of combat who learn devastating maneuvers and\ntechniques. You grow stronger through battle, honing your skills\nagainst increasingly formidable foes.",
		HP:          150,
		Mana:        50,
		StartingSkills: map[int]int{
			2001: 10, // Slash maneuver at 10%
		},
		StartingEquipment: []string{"leather armor", "sword"},
	}

	WandererKit = &StarterKit{
		ID:          3,
		Name:        "Wanderer's Kit",
		Description: "Light armor, a simple weapon, and a tattered spell scroll. Walk your own path between combat and magic through self-teaching.",
		FlavorText:  "Wanderers embrace both blade and spell, forging their own unique path.\nYou begin with knowledge of both disciplines, allowing you to adapt\nto any challenge.",
		HP:          125,
		Mana:        100,
		StartingSkills: map[int]int{
			2001: 5, // Slash maneuver at 5%
		},
		StartingEquipment: []string{"leather armor", "sword", "spell scroll"},
	}
)

var allKits = []*StarterKit{ScholarKit, WarriorKit, WandererKit}

// GetAll returns all available starter kits
func GetAll() []*StarterKit {
	return allKits
}

// GetByID returns a starter kit by its ID
func GetByID(id int) *StarterKit {
	for _, kit := range allKits {
		if kit.ID == id {
			return kit
		}
	}
	return nil
}

// GetByMenuChoice returns a starter kit by menu selection (1-based)
func GetByMenuChoice(choice int) *StarterKit {
	if choice < 1 || choice > len(allKits) {
		return nil
	}
	return allKits[choice-1]
}

// MenuString generates a menu display of all starter kits
func MenuString() string {
	result := ""
	for i, kit := range allKits {
		result += fmt.Sprintf("  %d) %s\n", i+1, kit.Name)
		result += fmt.Sprintf("     %s\n", kit.Description)
	}
	return result
}
