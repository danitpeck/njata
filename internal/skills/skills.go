package skills

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

type Targeting struct {
	Mode   string `json:"mode"`   // hostile_single, hostile_area, ally_single, self, any_object
	Range  int    `json:"range"`  // Maximum range in squares
	Radius int    `json:"radius"` // Area effect radius (0 for single target)
}

type Effects struct {
	Damage     string  `json:"damage"`      // Damage formula: "4d8 + I"
	DamageType string  `json:"damage_type"` // fire, cold, magic, none
	SaveType   string  `json:"save_type"`   // reflex, will, fortitude, none
	SaveDC     int     `json:"save_dc"`     // Difficulty class for save
	Healing    string  `json:"healing"`     // Healing formula (if applicable)
	Affect     *Affect `json:"affect"`      // Optional: buff/debuff effect
}

type Affect struct {
	Name        string         `json:"name"`
	Duration    string         `json:"duration"` // "120 + W*4" for formula-based
	Description string         `json:"description"`
	ACPenalty   int            `json:"ac_penalty"` // Positive = worse AC
	StatMods    map[string]int `json:"stat_mods"`
}

type Messages struct {
	Cast     string `json:"cast"`
	Hit      string `json:"hit"`
	Miss     string `json:"miss"`
	Save     string `json:"save"`
	CastRoom string `json:"cast_room"`
}

type Spell struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ManaCost        int       `json:"mana_cost"`
	CooldownSeconds int       `json:"cooldown_seconds"`
	LevelRequired   int       `json:"level_required"`
	Targeting       Targeting `json:"targeting"`
	Effects         Effects   `json:"effects"`
	Messages        Messages  `json:"messages"`
}

var (
	spellRegistry map[int]*Spell
	spellsByName  map[string]int
	mu            sync.RWMutex
)

func init() {
	spellRegistry = make(map[int]*Spell)
	spellsByName = make(map[string]int)
}

// Load loads all spells from a JSON file
func Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read spells file: %w", err)
	}

	var spells []Spell
	if err := json.Unmarshal(data, &spells); err != nil {
		return fmt.Errorf("failed to parse spells JSON: %w", err)
	}

	mu.Lock()
	defer mu.Unlock()

	spellRegistry = make(map[int]*Spell)
	spellsByName = make(map[string]int)

	for i := range spells {
		spell := &spells[i]
		spellRegistry[spell.ID] = spell
		spellsByName[strings.ToLower(spell.Name)] = spell.ID
	}

	return nil
}

// GetSpell retrieves a spell by ID
func GetSpell(id int) *Spell {
	mu.RLock()
	defer mu.RUnlock()
	return spellRegistry[id]
}

// GetSpellByName retrieves a spell by name (case-insensitive)
func GetSpellByName(name string) *Spell {
	mu.RLock()
	defer mu.RUnlock()
	id, ok := spellsByName[strings.ToLower(name)]
	if !ok {
		return nil
	}
	return spellRegistry[id]
}

// FindSpellByPartial finds a spell by partial name match (case-insensitive)
// Returns the spell and the number of matches found
func FindSpellByPartial(query string) (*Spell, int) {
	mu.RLock()
	defer mu.RUnlock()

	query = strings.ToLower(query)
	var matches []*Spell

	for _, spell := range spellRegistry {
		if strings.Contains(strings.ToLower(spell.Name), query) {
			matches = append(matches, spell)
		}
	}

	if len(matches) == 1 {
		return matches[0], 1
	} else if len(matches) > 1 {
		return nil, len(matches) // Ambiguous
	}
	return nil, 0 // Not found
}

// AllSpells returns all loaded spells
func AllSpells() map[int]*Spell {
	mu.RLock()
	defer mu.RUnlock()
	return spellRegistry
}

// PlayerSkillProgress tracks a player's proficiency with a spell
type PlayerSkillProgress struct {
	SpellID       int   `json:"spell_id"`
	Proficiency   int   `json:"proficiency"` // 0-100
	Learned       bool  `json:"learned"`
	LifetimeCasts int   `json:"lifetime_casts"`
	LastCastTime  int64 `json:"last_cast_time"` // Unix nanoseconds for cooldown
}

// CanCast checks if player can cast this spell (has mana, off cooldown, learned, etc)
func (psp *PlayerSkillProgress) CanCast(spell *Spell, currentMana int, currentTime int64) (bool, string) {
	if spell == nil {
		return false, "Skill not found"
	}

	if currentMana < spell.ManaCost {
		return false, fmt.Sprintf("Not enough mana (need %d, have %d)", spell.ManaCost, currentMana)
	}

	// Check cooldown using nanosecond timestamps
	timeSinceCastSecs := int64((currentTime - psp.LastCastTime) / 1_000_000_000)
	cooldownSecs := int64(spell.CooldownSeconds)

	if timeSinceCastSecs < cooldownSecs {
		remaining := cooldownSecs - timeSinceCastSecs
		return false, fmt.Sprintf("Skill on cooldown (%d seconds remaining)", remaining)
	}

	return true, ""
}

// IsReady checks if a spell is ready to cast (off cooldown)
func (psp *PlayerSkillProgress) IsReady(spell *Spell, currentTime int64) bool {
	if spell == nil {
		return false
	}
	timeSinceCastSecs := int64((currentTime - psp.LastCastTime) / 1_000_000_000)
	cooldownSecs := int64(spell.CooldownSeconds)
	return timeSinceCastSecs >= cooldownSecs
}

// GetCooldownRemaining returns seconds remaining on cooldown
func (psp *PlayerSkillProgress) GetCooldownRemaining(spell *Spell, currentTime int64) int64 {
	if spell == nil {
		return 0
	}
	timeSinceCastSecs := int64((currentTime - psp.LastCastTime) / 1_000_000_000)
	cooldownSecs := int64(spell.CooldownSeconds)
	if timeSinceCastSecs >= cooldownSecs {
		return 0
	}
	return cooldownSecs - timeSinceCastSecs
}

// UpdateCooldown sets the last cast time to now
func (psp *PlayerSkillProgress) UpdateCooldown() {
	psp.LastCastTime = time.Now().UnixNano()
}

// UpdateProficiency increases proficiency (for learning through use)
func (psp *PlayerSkillProgress) UpdateProficiency(gain int) {
	psp.Proficiency += gain
	if psp.Proficiency > 100 {
		psp.Proficiency = 100
	}
	psp.LifetimeCasts++
}
