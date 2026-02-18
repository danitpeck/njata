package persist

import (
    "testing"

    "njata/internal/game"
)

func TestSaveAndLoadPlayer(t *testing.T) {
    dir := t.TempDir()

    record := PlayerRecord{
        Name:     "Alice",
        Location: 123,
        Hair:     "a wild mane",
        Eyes:     "green eyes",
        Inventory: []game.Object{
            {
                Vnum:  42,
                Short: "a smooth stone",
            },
        },
        Equipment: map[string]game.Object{
            "head": {
                Vnum:  99,
                Short: "a battered helm",
            },
        },
    }
    if err := SavePlayer(dir, record); err != nil {
        t.Fatalf("save player: %v", err)
    }

    loaded, ok, err := LoadPlayer(dir, "Alice")
    if err != nil {
        t.Fatalf("load player: %v", err)
    }
    if !ok || loaded == nil {
        t.Fatalf("expected player record")
    }
    if loaded.Name != "Alice" || loaded.Location != 123 || loaded.Hair != "a wild mane" || loaded.Eyes != "green eyes" {
        t.Fatalf("unexpected record: %+v", loaded)
    }
    if len(loaded.Inventory) != 1 || loaded.Inventory[0].Vnum != 42 {
        t.Fatalf("unexpected inventory: %+v", loaded.Inventory)
    }
    if len(loaded.Equipment) != 1 || loaded.Equipment["head"].Vnum != 99 {
        t.Fatalf("unexpected equipment: %+v", loaded.Equipment)
    }
}

func TestLoadPlayerMissing(t *testing.T) {
    dir := t.TempDir()

    loaded, ok, err := LoadPlayer(dir, "Missing")
    if err != nil {
        t.Fatalf("load missing player: %v", err)
    }
    if ok || loaded != nil {
        t.Fatalf("expected missing player result")
    }
}
