package persist

import "testing"

func TestSaveAndLoadPlayer(t *testing.T) {
    dir := t.TempDir()

    record := PlayerRecord{Name: "Alice", Location: 123}
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
    if loaded.Name != "Alice" || loaded.Location != 123 {
        t.Fatalf("unexpected record: %+v", loaded)
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
