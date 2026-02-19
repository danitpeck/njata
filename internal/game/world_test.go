package game

import (
	"strings"
	"testing"
)

func TestBroadcastSay(t *testing.T) {
	world := CreateDefaultWorld()

	aliceOut := &bufferOutput{}
	bobOut := &bufferOutput{}

	speaker := &Player{Name: "Alice", Output: aliceOut}
	if err := world.AddPlayer(speaker); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := world.AddPlayer(&Player{Name: "Bob", Output: bobOut}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	world.BroadcastSay(speaker, "hello")

	if !aliceOut.Contains("You say 'hello'") {
		t.Fatalf("expected sender message not found")
	}
	if !bobOut.Contains("Alice says 'hello'") {
		t.Fatalf("expected receiver message not found")
	}
}

func TestRespawnTick(t *testing.T) {
	// Create world with test rooms and mobs
	world := CreateDefaultWorld()

	testMob := &Mobile{
		Vnum:     101,
		Keywords: []string{"goblin"},
		Short:    "A goblin",
		Long:     "A mean goblin",
		Level:    3,
		MaxHP:    20,
		HP:       20,
	}

	testObj := &Object{
		Vnum:     201,
		Keywords: []string{"sword"},
		Short:    "A steel sword",
		Type:     "weapon",
	}

	// Create a test room
	testRoom := &Room{
		Vnum:             1001,
		Name:             "Test Room",
		Description:      "A test room",
		AreaName:         "Test Area",
		AreaResetMinutes: 1, // 1 minute reset
		Mobiles:          []*Mobile{},
		Objects:          []*Object{},
		MobileResets:     []Reset{{MobVnum: 101, Count: 1}},
		ObjectResets:     []Reset{{ObjVnum: 201, Count: 1}},
	}

	world.rooms[1001] = testRoom
	world.SetPrototypes(
		map[int]*Mobile{101: testMob},
		map[int]*Object{201: testObj},
	)

	// Verify initial spawn is present
	if len(testRoom.Mobiles) == 0 || len(testRoom.Objects) == 0 {
		// Initial spawn from resets happens in loader, not in CreateWorld
		// So we manually instantiate to simulate that
		testRoom.Mobiles = []*Mobile{testMob}
		testRoom.Objects = []*Object{testObj}
	}

	initialMobCount := len(testRoom.Mobiles)
	initialObjCount := len(testRoom.Objects)

	// Capture logs
	var logs []string
	logger := func(msg string) {
		logs = append(logs, msg)
	}

	// First tick - should mark area as respawned
	world.RespawnTick(60, logger)

	if len(logs) == 0 {
		t.Fatalf("expected log output from first RespawnTick")
	}
	if !strings.Contains(logs[0], "initial") {
		t.Fatalf("expected 'initial' in first tick log, got: %s", logs[0])
	}

	// Clear mobs/objects to simulate passage of time
	testRoom.Mobiles = []*Mobile{}
	testRoom.Objects = []*Object{}

	if len(testRoom.Mobiles) != 0 || len(testRoom.Objects) != 0 {
		t.Fatalf("failed to clear mobs/objects")
	}

	// Manually advance the last respawn time backward to simulate time passing
	world.mu.Lock()
	world.areaLastRespawn["Test Area"] = world.areaLastRespawn["Test Area"].Add(-2 * 60000000000) // Subtract 2 minutes
	world.mu.Unlock()

	logs = []string{}

	// Second tick - should respawn
	world.RespawnTick(60, logger)

	if len(testRoom.Mobiles) != initialMobCount {
		t.Fatalf("expected %d mobs after respawn, got %d", initialMobCount, len(testRoom.Mobiles))
	}
	if len(testRoom.Objects) != initialObjCount {
		t.Fatalf("expected %d objects after respawn, got %d", initialObjCount, len(testRoom.Objects))
	}

	if len(logs) == 0 {
		t.Fatalf("expected log output from respawn tick")
	}
	if !strings.Contains(logs[0], "areas respawned") {
		t.Fatalf("expected respawn message, got: %s", logs[0])
	}
	if !strings.Contains(logs[0], "Test Area") {
		t.Fatalf("expected Test Area in respawn message, got: %s", logs[0])
	}
}
