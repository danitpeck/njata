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
    Name       string `json:"name"`
    Location   int    `json:"location"`
    Class      int    `json:"class"`
    Race       int    `json:"race"`
    Sex        int    `json:"sex"`
    Level      int    `json:"level"`
    HP         int    `json:"hp"`
    MaxHP      int    `json:"max_hp"`
    Mana       int    `json:"mana"`
    MaxMana    int    `json:"max_mana"`
    Move       int    `json:"move"`
    MaxMove    int    `json:"max_move"`
    Gold       int    `json:"gold"`
    Experience int    `json:"experience"`
    Attributes [7]int `json:"attributes"`
    Alignment  int    `json:"alignment"`
    Hitroll    int    `json:"hitroll"`
    Damroll    int    `json:"damroll"`
    Armor      int    `json:"armor"`
    Skills     map[int]*skills.PlayerSkillProgress `json:"skills"`
    IsKeeper   bool   `json:"is_keeper"`
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

    return PlayerRecord{
        Name:        p.Name,
        Location:    p.Location,
        Class:       p.Class,
        Race:        p.Race,
        Sex:         p.Sex,
        Level:       p.Level,
        HP:          p.HP,
        MaxHP:       p.MaxHP,
        Mana:        p.Mana,
        MaxMana:     p.MaxMana,
        Move:        p.Move,
        MaxMove:     p.MaxMove,
        Gold:        p.Gold,
        Experience:  p.Experience,
        Attributes:  p.Attributes,
        Alignment:   p.Alignment,
        Hitroll:     p.Hitroll,
        Damroll:     p.Damroll,
        Armor:       p.Armor,
        Skills:      skillsCopy,
        IsKeeper:    p.IsKeeper,
    }
}

// RecordToPlayer applies a PlayerRecord's data to a game.Player
func RecordToPlayer(p *game.Player, r *PlayerRecord) {
    p.Location = r.Location
    p.Class = r.Class
    p.Race = r.Race
    p.Sex = r.Sex
    p.Level = r.Level
    p.HP = r.HP
    p.MaxHP = r.MaxHP
    p.Mana = r.Mana
    p.MaxMana = r.MaxMana
    p.Move = r.Move
    p.MaxMove = r.MaxMove
    p.Gold = r.Gold
    p.Experience = r.Experience
    p.Attributes = r.Attributes
    p.Alignment = r.Alignment
    p.Hitroll = r.Hitroll
    p.Damroll = r.Damroll
    p.Armor = r.Armor
    p.Skills = r.Skills
    p.IsKeeper = r.IsKeeper
}
