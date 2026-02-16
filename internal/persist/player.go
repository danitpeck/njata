package persist

import (
    "encoding/json"
    "errors"
    "os"
    "path/filepath"
    "strings"
)

type PlayerRecord struct {
    Name     string `json:"name"`
    Location int    `json:"location"`
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
