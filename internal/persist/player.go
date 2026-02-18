package persist

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"njata/internal/game"
	"njata/internal/skills"
)

type PlayerRecord struct {
	Name         string                              `json:"name"`
	Location     int                                 `json:"location"`
	Race         int                                 `json:"race"`
	Sex          int                                 `json:"sex"`
	Hair         string                              `json:"hair"`
	Eyes         string                              `json:"eyes"`
	HP           int                                 `json:"hp"`
	MaxHP        int                                 `json:"max_hp"`
	Mana         int                                 `json:"mana"`
	MaxMana      int                                 `json:"max_mana"`
	Gold         int                                 `json:"gold"`
	Strength     int                                 `json:"strength"`
	Dexterity    int                                 `json:"dexterity"`
	Constitution int                                 `json:"constitution"`
	Intelligence int                                 `json:"intelligence"`
	Wisdom       int                                 `json:"wisdom"`
	Charisma     int                                 `json:"charisma"`
	Luck         int                                 `json:"luck"`
	Armor        int                                 `json:"armor"`
	Skills       map[int]*skills.PlayerSkillProgress `json:"skills"`
	IsKeeper     bool                                `json:"is_keeper"`
	Inventory    []game.Object                       `json:"inventory"`
	Equipment    map[string]game.Object              `json:"equipment"`
}

func LoadPlayer(dir string, name string) (*PlayerRecord, bool, error) {
	path := playerPath(dir, name)
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, false, nil
		}
		return nil, false, err
	}

	var record PlayerRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, false, err
	}

	return &record, true, nil
}

func SavePlayer(dir string, record PlayerRecord) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return err
	}

	path := playerPath(dir, record.Name)
	return os.WriteFile(path, data, 0644)
}

func playerPath(dir string, name string) string {
	normalized := strings.ToLower(name)
	normalized = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			return r
		}
		return '-'
	}, normalized)
	return filepath.Join(dir, normalized+".json")
}

// PlayerToRecord converts a game.Player to a PlayerRecord for saving
func PlayerToRecord(p *game.Player) PlayerRecord {
	// Deep copy the Skills map to ensure proper serialization
	skillsCopy := make(map[int]*skills.PlayerSkillProgress)
	if p.Skills != nil {
		for k, v := range p.Skills {
			skillsCopy[k] = v
		}
	}

	inventoryCopy := make([]game.Object, 0, len(p.Inventory))
	for _, item := range p.Inventory {
		if item != nil {
			inventoryCopy = append(inventoryCopy, *item)
		}
	}

	equipmentCopy := map[string]game.Object{}
	for slot, item := range p.Equipment {
		if item != nil {
			equipmentCopy[slot] = *item
		}
	}

	return PlayerRecord{
		Name:         p.Name,
		Location:     p.Location,
		Race:         p.Race,
		Sex:          p.Sex,
		Hair:         p.Hair,
		Eyes:         p.Eyes,
		HP:           p.HP,
		MaxHP:        p.MaxHP,
		Mana:         p.Mana,
		MaxMana:      p.MaxMana,
		Gold:         p.Gold,
		Strength:     p.Strength,
		Dexterity:    p.Dexterity,
		Constitution: p.Constitution,
		Intelligence: p.Intelligence,
		Wisdom:       p.Wisdom,
		Charisma:     p.Charisma,
		Luck:         p.Luck,
		Armor:        p.Armor,
		Skills:       skillsCopy,
		IsKeeper:     p.IsKeeper,
		Inventory:    inventoryCopy,
		Equipment:    equipmentCopy,
	}
}

// RecordToPlayer applies a PlayerRecord's data to a game.Player
func RecordToPlayer(p *game.Player, r *PlayerRecord) {
	p.Location = r.Location
	p.Race = r.Race
	p.Sex = r.Sex
	p.Hair = r.Hair
	p.Eyes = r.Eyes
	p.HP = r.HP
	p.MaxHP = r.MaxHP
	p.Mana = r.Mana
	p.MaxMana = r.MaxMana
	p.Gold = r.Gold
	p.Strength = r.Strength
	p.Dexterity = r.Dexterity
	p.Constitution = r.Constitution
	p.Intelligence = r.Intelligence
	p.Wisdom = r.Wisdom
	p.Charisma = r.Charisma
	p.Luck = r.Luck
	p.Armor = r.Armor
	p.Skills = r.Skills
	p.IsKeeper = r.IsKeeper
	if len(r.Inventory) > 0 {
		p.Inventory = make([]*game.Object, 0, len(r.Inventory))
		for _, item := range r.Inventory {
			obj := item
			p.Inventory = append(p.Inventory, &obj)
		}
	}
	if len(r.Equipment) > 0 {
		p.Equipment = make(map[string]*game.Object)
		for slot, item := range r.Equipment {
			obj := item
			p.Equipment[slot] = &obj
		}
	}
}
