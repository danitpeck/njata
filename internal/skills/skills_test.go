package skills

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

var skillsPath string

func init() {
	// Build the path to skills.json relative to this test file's location
	// This file is in internal/skills, so we go up 2 levels to project root
	_, thisFile, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(thisFile)))
	skillsPath = filepath.Join(projectRoot, "skills", "skills.json")
}

func TestLoadSpells(t *testing.T) {
	// Load the skills.json file
	err := Load(skillsPath)
	if err != nil {
		t.Fatalf("Failed to load spells: %v", err)
	}

	// Verify we have spells loaded
	spells := AllSpells()
	if len(spells) == 0 {
		t.Fatal("Expected spells to be loaded, got none")
	}

	// We expect 8 MVP spells + 1 Slash ability
	if len(spells) != 9 {
		t.Errorf("Expected 9 spells, got %d", len(spells))
	}
}

func TestGetSpellByID(t *testing.T) {
	err := Load(skillsPath)
	if err != nil {
		t.Fatalf("Failed to load spells: %v", err)
	}

	tests := []struct {
		id           int
		expectedName string
		shouldExist  bool
	}{
		{1001, "Arcane Bolt", true},
		{1002, "Leviathan's Fire", true},
		{1003, "Mend", true},
		{1004, "Shadow Veil", true},
		{1005, "Ephemeral Step", true},
		{1006, "Path Shift", true},
		{1007, "Winter's Whisper", true},
		{1008, "Knowing", true},
		{9999, "", false},
	}

	for _, tt := range tests {
		spell := GetSpell(tt.id)
		if tt.shouldExist {
			if spell == nil {
				t.Errorf("GetSpell(%d): expected spell to exist, got nil", tt.id)
			} else if spell.Name != tt.expectedName {
				t.Errorf("GetSpell(%d): expected name %q, got %q", tt.id, tt.expectedName, spell.Name)
			}
		} else {
			if spell != nil {
				t.Errorf("GetSpell(%d): expected nil, got spell %q", tt.id, spell.Name)
			}
		}
	}
}

func TestGetSpellByName(t *testing.T) {
	err := Load(skillsPath)
	if err != nil {
		t.Fatalf("Failed to load spells: %v", err)
	}

	tests := []struct {
		name        string
		expectedID  int
		shouldExist bool
	}{
		{"arcane bolt", 1001, true},
		{"Arcane Bolt", 1001, true}, // Case insensitive
		{"ARCANE BOLT", 1001, true}, // Case insensitive
		{"leviathan's fire", 1002, true},
		{"mend", 1003, true},
		{"nonexistent spell", 0, false},
	}

	for _, tt := range tests {
		spell := GetSpellByName(tt.name)
		if tt.shouldExist {
			if spell == nil {
				t.Errorf("GetSpellByName(%q): expected spell to exist, got nil", tt.name)
			} else if spell.ID != tt.expectedID {
				t.Errorf("GetSpellByName(%q): expected ID %d, got %d", tt.name, tt.expectedID, spell.ID)
			}
		} else {
			if spell != nil {
				t.Errorf("GetSpellByName(%q): expected nil, got spell ID %d", tt.name, spell.ID)
			}
		}
	}
}

func TestSpellStructure(t *testing.T) {
	err := Load(skillsPath)
	if err != nil {
		t.Fatalf("Failed to load spells: %v", err)
	}

	spell := GetSpell(1001)
	if spell == nil {
		t.Fatal("Expected to find spell 1001")
	}

	// Verify spell has required fields
	if spell.ID != 1001 {
		t.Errorf("Spell ID: expected 1001, got %d", spell.ID)
	}
	if spell.Name == "" {
		t.Error("Spell name is empty")
	}
	if spell.ManaCost <= 0 {
		t.Errorf("Spell mana cost should be positive, got %d", spell.ManaCost)
	}
	if spell.CooldownSeconds < 0 {
		t.Errorf("Spell cooldown should be non-negative, got %d", spell.CooldownSeconds)
	}
	if spell.LevelRequired < 1 {
		t.Errorf("Spell level required should be at least 1, got %d", spell.LevelRequired)
	}
}

func TestPlayerSkillProgress(t *testing.T) {
	progress := &PlayerSkillProgress{
		SpellID:       1001,
		Proficiency:   50,
		Learned:       true,
		LifetimeCasts: 5,
		LastCastTime:  1000000,
	}

	if progress.SpellID != 1001 {
		t.Errorf("Expected spell ID 1001, got %d", progress.SpellID)
	}
	if progress.Proficiency != 50 {
		t.Errorf("Expected proficiency 50, got %d", progress.Proficiency)
	}
	if !progress.Learned {
		t.Error("Expected spell to be marked as learned")
	}
	if progress.LifetimeCasts != 5 {
		t.Errorf("Expected 5 lifetime casts, got %d", progress.LifetimeCasts)
	}
}

func TestColisthaSpellNames(t *testing.T) {
	err := Load(skillsPath)
	if err != nil {
		t.Fatalf("Failed to load spells: %v", err)
	}

	// Verify that spells have Colista-themed names, not generic ones
	expectedNames := map[int]string{
		1001: "Arcane Bolt",
		1002: "Leviathan's Fire",
		1003: "Mend",
		1004: "Shadow Veil",
		1005: "Ephemeral Step",
		1006: "Path Shift",
		1007: "Winter's Whisper",
		1008: "Knowing",
	}

	for id, expectedName := range expectedNames {
		spell := GetSpell(id)
		if spell == nil {
			t.Errorf("Spell %d not found", id)
			continue
		}
		if spell.Name != expectedName {
			t.Errorf("Spell %d: expected name %q, got %q", id, expectedName, spell.Name)
		}
	}
}

func TestSpellMessages(t *testing.T) {
	err := Load(skillsPath)
	if err != nil {
		t.Fatalf("Failed to load spells: %v", err)
	}

	spell := GetSpell(1001)
	if spell == nil {
		t.Fatal("Expected to find spell 1001")
	}

	// Verify spell has messages
	if spell.Messages.Cast == "" {
		t.Error("Spell 1001 missing cast message")
	}
	if spell.Messages.Hit == "" {
		t.Error("Spell 1001 missing hit message")
	}
	if spell.Messages.Miss == "" {
		t.Error("Spell 1001 missing miss message")
	}
	if spell.Messages.CastRoom == "" {
		t.Error("Spell 1001 missing cast_room message")
	}
}

func TestSpellTargeting(t *testing.T) {
	err := Load(skillsPath)
	if err != nil {
		t.Fatalf("Failed to load spells: %v", err)
	}

	spell := GetSpell(1001)
	if spell == nil {
		t.Fatal("Expected to find spell 1001")
	}

	// Verify targeting information
	if spell.Targeting.Mode == "" {
		t.Error("Spell 1001 missing targeting mode")
	}
	if spell.Targeting.Range < 0 {
		t.Errorf("Spell 1001 range should be non-negative, got %d", spell.Targeting.Range)
	}
	if spell.Targeting.Radius < 0 {
		t.Errorf("Spell 1001 radius should be non-negative, got %d", spell.Targeting.Radius)
	}
}

func TestSpellEffects(t *testing.T) {
	err := Load(skillsPath)
	if err != nil {
		t.Fatalf("Failed to load spells: %v", err)
	}

	// Test damage spell (Arcane Bolt)
	spell := GetSpell(1001)
	if spell == nil {
		t.Fatal("Expected to find spell 1001")
	}
	if spell.Effects.Damage == "" {
		t.Error("Arcane Bolt should have damage formula")
	}
	if spell.Effects.DamageType == "" {
		t.Error("Arcane Bolt should have damage type")
	}

	// Test healing spell (Mend)
	spell = GetSpell(1003)
	if spell == nil {
		t.Fatal("Expected to find spell 1003")
	}
	if spell.Effects.Healing == "" {
		t.Error("Mend should have healing formula")
	}
}

func TestAllSpells(t *testing.T) {
	err := Load(skillsPath)
	if err != nil {
		t.Fatalf("Failed to load spells: %v", err)
	}

	spells := AllSpells()
	if len(spells) != 9 {
		t.Errorf("Expected 9 spells, got %d", len(spells))
	}

	// Verify we can iterate through all spells
	spellCount := 0
	for id, spell := range spells {
		if spell == nil {
			t.Errorf("Spell map contains nil for ID %d", id)
		}
		if spell.ID != id {
			t.Errorf("Spell map mismatch: key %d, spell.ID %d", id, spell.ID)
		}
		spellCount++
	}

	if spellCount != 9 {
		t.Errorf("Expected 9 spells in iteration, got %d", spellCount)
	}
}

// Helper to reset spells between tests if needed
func resetSpells() {
	spellRegistry = make(map[int]*Spell)
	spellsByName = make(map[string]int)
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
