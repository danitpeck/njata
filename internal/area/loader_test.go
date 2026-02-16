package area

import (
    "os"
    "path/filepath"
    "testing"
)

func TestLoadRoomsFromDir(t *testing.T) {
    dir := t.TempDir()
    content := `#FUSSAREA
#ROOM
Vnum     100
Name     Test Room~
Desc     First line
Second line
~
#EXDESC
ExDescKey    sign plaque~
ExDesc       A test sign.
~
#ENDEXDESC
#EXIT
Direction north~
ToRoom    101
#ENDEXIT
#ENDROOM

#ROOM
Vnum     101
Name     Second Room~
Desc     Second desc~
#ENDROOM
`

    filePath := filepath.Join(dir, "sample.are")
    if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
        t.Fatalf("write file: %v", err)
    }

    rooms, start, err := LoadRoomsFromDir(dir)
    if err != nil {
        t.Fatalf("load rooms: %v", err)
    }

    if len(rooms) != 2 {
        t.Fatalf("expected 2 rooms, got %d", len(rooms))
    }
    if start != 100 {
        t.Fatalf("expected start vnum 100, got %d", start)
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
    if room.Exits["north"] != 101 {
        t.Fatalf("expected north exit to 101")
    }
    if room.ExDescs["sign"] != "A test sign." {
        t.Fatalf("expected exdesc for sign")
    }
}
