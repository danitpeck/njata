package area

import (
    "os"
    "path/filepath"
    "testing"
)

func TestLoadRoomsFromDir(t *testing.T) {
    dir := t.TempDir()
    content := `{
  "name": "Test Area",
  "author": "Test Author",
  "rooms": {
    "100": {
      "vnum": 100,
      "name": "Test Room",
      "description": "First line\nSecond line",
      "sector": "city",
      "flags": {
        "nomob": true,
        "indoors": true
      },
      "exits": {
        "north": 101
      },
      "exdescs": {
        "sign": "A test sign.",
        "plaque": "A test sign."
      },
      "area_name": "Test Area",
      "area_author": "Test Author"
    },
    "101": {
      "vnum": 101,
      "name": "Second Room",
      "description": "Second desc",
      "sector": "city",
      "flags": {},
      "exits": {},
      "exdescs": {},
      "area_name": "Test Area",
      "area_author": "Test Author"
    }
  }
}`

    filePath := filepath.Join(dir, "sample.json")
    if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
        t.Fatalf("write file: %v", err)
    }

    rooms, mobiles, objects, start, err := LoadRoomsFromDir(dir)
    if err != nil {
        t.Fatalf("load rooms: %v", err)
    }

    if len(rooms) != 2 {
        t.Fatalf("expected 2 rooms, got %d", len(rooms))
    }
    if start != 100 {
        t.Fatalf("expected start vnum 100, got %d", start)
    }

    // Verify prototypes are returned even if empty for this test
    if mobiles == nil || objects == nil {
        t.Fatalf("expected non-nil prototypes")
    }

    room := rooms[100]
    if room == nil {
        t.Fatalf("room 100 missing")
    }
    if room.Name != "Test Room" {
        t.Fatalf("unexpected room name: %s", room.Name)
    }
    if room.Description != "First line\nSecond line" {
        t.Fatalf("unexpected room desc: %s", room.Description)
    }
    if room.Sector != "city" {
        t.Fatalf("expected sector city, got %s", room.Sector)
    }
    if !room.Flags["nomob"] || !room.Flags["indoors"] {
        t.Fatalf("expected nomob and indoors flags")
    }
    if room.Exits["north"] != 101 {
        t.Fatalf("expected north exit to 101")
    }
    if room.ExDescs["sign"] != "A test sign." {
        t.Fatalf("expected exdesc for sign")
    }
}
