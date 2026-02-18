package game

import "testing"

func TestMovePlayer(t *testing.T) {
    rooms := map[int]*Room{
        1: {Vnum: 1, Name: "One", Description: "Room one", Exits: map[string]int{"east": 2}},
        2: {Vnum: 2, Name: "Two", Description: "Room two", Exits: map[string]int{"west": 1}},
    }
    world := CreateWorldFromRooms(rooms, 1)

    player := &Player{Name: "Alice", Output: &bufferOutput{}}
    if err := world.AddPlayer(player); err != nil {
        t.Fatalf("add player: %v", err)
    }

    view, err := world.MovePlayer(player, "east")
    if err != nil {
        t.Fatalf("move player: %v", err)
    }
    if player.Location != 2 {
        t.Fatalf("expected location 2, got %d", player.Location)
    }
    if view.Name != "Two" {
        t.Fatalf("expected room Two, got %s", view.Name)
    }
}
