package skills

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Skill struct {
	ID               int      `json:"id"`
	Name             string   `json:"name"`
	Type             string   `json:"type"` // offensive, healing, defensive, etc
	ManaCost         int      `json:"mana_cost"`
	CooldownSeconds  int      `json:"cooldown_seconds"`
	DamageFormula    string   `json:"damage_formula"`    // e.g., "int_bonus * 1.2 + 10"
	HealFormula      string   `json:"heal_formula"`
	ACBonus          int      `json:"ac_bonus"`
	DurationSeconds  int      `json:"duration_seconds"`
	Description      string   `json:"description"`
	Classes          []int    `json:"classes"` // Which classes can learn this
}

type SkillsData struct {
	Skills map[string]*Skill `json:"skills"`
}

var (
	skillsDB map[int]*Skill
	mu       sync.RWMutex
)

func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var sd SkillsData
	if err := json.Unmarshal(data, &sd); err != nil {
		return err
	}

	mu.Lock()
	defer mu.Unlock()

	skillsDB = make(map[int]*Skill)
	for _, skill := range sd.Skills {
		if skill.ID > 0 {
			skillsDB[skill.ID] = skill
		}
	}

	return nil
}

func GetSkill(id int) (*Skill, bool) {
	mu.RLock()
	defer mu.RUnlock()
	skill, ok := skillsDB[id]
	return skill, ok
}

func GetSkillsByClass(classID int) []*Skill {
	mu.RLock()
	defer mu.RUnlock()

	var result []*Skill
	for _, skill := range skillsDB {
		for _, c := range skill.Classes {
			if c == classID {
				result = append(result, skill)
				break
			}
		}
	}
	return result
}

func ListAll() []*Skill {
	mu.RLock()
	defer mu.RUnlock()

	result := make([]*Skill, 0, len(skillsDB))
	for _, skill := range skillsDB {
		result = append(result, skill)
	}
	return result
}

// PlayerSkill tracks a player's proficiency with a skill
type PlayerSkill struct {
	SkillID          int   `json:"skill_id"`
	Proficiency      int   `json:"proficiency"` // 0-100 or experience points
	LastCastUnixTime int64 `json:"last_cast_unix_time"` // For cooldown tracking
}

// CanCast checks if player can cast this skill (has mana, off cooldown, etc)
func (ps *PlayerSkill) CanCast(skill *Skill, currentMana int, currentTime int64) (bool, string) {
	if skill == nil {
		return false, "Skill not found"
	}

	if currentMana < skill.ManaCost {
		return false, fmt.Sprintf("Not enough mana (need %d, have %d)", skill.ManaCost, currentMana)
	}

	timeSinceCast := currentTime - ps.LastCastUnixTime
	cooldown := int64(skill.CooldownSeconds)

	if timeSinceCast < cooldown {
		return false, fmt.Sprintf("Skill on cooldown (%d seconds remaining)", cooldown-timeSinceCast)
	}

	return true, ""
}

// StringifySkills returns a formatted list of skills
func StringifySkills(skills []*Skill) string {
	if len(skills) == 0 {
		return "No skills available"
	}

	result := "Available Skills:\n"
	for _, s := range skills {
		result += fmt.Sprintf("  %s - %s (Mana: %d, Cooldown: %ds)\n", s.Name, s.Description, s.ManaCost, s.CooldownSeconds)
	}
	return result
}
