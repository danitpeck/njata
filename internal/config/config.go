package config

import (
    "encoding/json"
    "os"
)

// Config holds runtime configuration loaded from disk.
type Config struct {
    StartRoomVnum        int `json:"start_room_vnum"`
    RespawnDefaultMinutes int `json:"respawn_default_minutes"`
}

// Load reads the config file if it exists. Missing files return defaults.
func Load(path string) (Config, error) {
    var cfg Config
    data, err := os.ReadFile(path)
    if err != nil {
        if os.IsNotExist(err) {
            return cfg, nil
        }
        return cfg, err
    }

    if err := json.Unmarshal(data, &cfg); err != nil {
        return cfg, err
    }

    return cfg, nil
}
