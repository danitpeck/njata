package commands

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"njata/internal/game"
	"njata/internal/persist"
	"njata/internal/races"
	"njata/internal/skills"
)

func RegisterBuiltins(registry *Registry) {
	registry.Register("look", cmdLook)
	registry.Register("say", cmdSay)
	registry.Register("who", cmdWho)
	registry.Register("stats", cmdStats)
	registry.Register("exits", cmdExits)
	registry.Register("autoexits", cmdAutoexits)
	registry.Register("astat", cmdAstat)
	registry.Register("abilities", cmdAbilities)
	registry.Register("cast", cmdCast)
	registry.Register("slash", cmdSlash)
	registry.Register("study", cmdStudy)
	registry.Register("save", cmdSave)
	registry.Register("makekeeper", cmdMakeKeeper)
	registry.Register("removekeeper", cmdRemoveKeeper)
	registry.Register("spawn", cmdSpawn)
	registry.Register("teleport", cmdTeleport)
	registerMovement(registry)
	registry.Register("help", func(ctx Context, args string) {
		commands := registry.List()
		ctx.Output.WriteLine("Commands: " + strings.Join(commands, ", "))
	})
	registry.Register("quit", cmdQuit)
}

// DisplayRoomView is a shared function to display a room view consistently
func DisplayRoomView(output game.Output, view game.RoomView, autoExits bool) {
	output.WriteLine(view.Name)
	if view.Description != "" {
		output.WriteLine(view.Description)
	}

	if autoExits {
		output.WriteLine(FormatExits(view.Exits))
	}

	// Display mobiles
	if len(view.Mobiles) > 0 {
		for _, mob := range view.Mobiles {
			output.WriteLine(mob)
		}
	}

	// Display objects
	if len(view.Objects) > 0 {
		for _, obj := range view.Objects {
			output.WriteLine(obj)
		}
	}

	// Display other players
	if len(view.Others) > 0 {
		output.WriteLine("Also here: " + strings.Join(view.Others, ", "))
	} else if len(view.Mobiles) == 0 && len(view.Objects) == 0 {
		output.WriteLine("You are alone here.")
	}
}

func cmdLook(ctx Context, args string) {
	trimmed := strings.TrimSpace(args)
	if trimmed != "" {
		keyword := trimmed
		if strings.HasPrefix(strings.ToLower(trimmed), "in ") {
			keyword = strings.TrimSpace(trimmed[3:])
		}

		if desc, ok := ctx.World.FindRoomExDesc(ctx.Player, keyword); ok {
			ctx.Output.WriteLine(desc)
			return
		}

		ctx.Output.WriteLine("You see nothing special.")
		return
	}

	view, err := ctx.World.DescribeRoom(ctx.Player)
	if err != nil {
		ctx.Output.WriteLine("You are nowhere.")
		return
	}

	DisplayRoomView(ctx.Output, view, ctx.Player.AutoExits)
}

func cmdSay(ctx Context, args string) {
	if strings.TrimSpace(args) == "" {
		ctx.Output.WriteLine("Say what?")
		return
	}

	ctx.World.BroadcastSay(ctx.Player, args)
}

func cmdWho(ctx Context, args string) {
	players := ctx.World.ListPlayers()
	ctx.Output.WriteLine(fmt.Sprintf("Players online (%d): %s", len(players), strings.Join(players, ", ")))
}

func cmdStats(ctx Context, args string) {
	p := ctx.Player
	ctx.Output.WriteLine(fmt.Sprintf("=== %s (Level %d) ===", p.Name, p.Level))

	raceName := "Unknown"
	if r := races.GetByID(p.Race); r != nil {
		raceName = r.Name
	}

	sexNames := []string{"neuter", "male", "female"}
	sexName := "unknown"
	if p.Sex >= 0 && p.Sex < len(sexNames) {
		sexName = sexNames[p.Sex]
	}

	ctx.Output.WriteLine(fmt.Sprintf("Race: %s | Sex: %s", raceName, sexName))
	ctx.Output.WriteLine("")
	ctx.Output.WriteLine(fmt.Sprintf("HP:    %d/%d | Mana: %d/%d | Move: %d/%d", p.HP, p.MaxHP, p.Mana, p.MaxMana, p.Move, p.MaxMove))
	ctx.Output.WriteLine(fmt.Sprintf("Experience: %d | Gold: %d", p.Experience, p.Gold))
	ctx.Output.WriteLine("")

	attrNames := []string{"STR", "INT", "WIS", "DEX", "CON", "LCK", "CHM"}
	for i, name := range attrNames {
		ctx.Output.WriteLine(fmt.Sprintf("%s: %2d", name, p.Attributes[i]))
	}

	ctx.Output.WriteLine("")
	ctx.Output.WriteLine(fmt.Sprintf("Alignment: %d | Hitroll: %d | Damroll: %d | Armor: %d", p.Alignment, p.Hitroll, p.Damroll, p.Armor))
}

func cmdAbilities(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in to view abilities")
		return
	}

	p := ctx.Player

	// Initialize skills map if needed
	if p.Skills == nil {
		p.Skills = make(map[int]*skills.PlayerSkillProgress)
	}

	ctx.Output.WriteLine("=== ABILITIES ===")
	ctx.Output.WriteLine("")

	if len(p.Skills) == 0 {
		ctx.Output.WriteLine("You haven't learned any abilities yet.")
		ctx.Output.WriteLine("")
		return
	}

	// Separate spells and maneuvers
	allSpells := skills.AllSpells()
	spellsList := make([]*skills.Spell, 0)
	maneuversList := make([]*skills.Spell, 0)

	for spellID, progress := range p.Skills {
		if !progress.Learned {
			continue
		}

		spell := allSpells[spellID]
		if spell == nil {
			continue
		}

		// Categorize: spells (1000-8999) vs maneuvers (9000+)
		if spell.ID >= 9000 {
			maneuversList = append(maneuversList, spell)
		} else {
			spellsList = append(spellsList, spell)
		}
	}

	// Display spells
	if len(spellsList) > 0 {
		ctx.Output.WriteLine("SPELLS:")
		for _, spell := range spellsList {
			progress := p.Skills[spell.ID]
			ctx.Output.WriteLine(fmt.Sprintf("  [%d] %s", spell.ID, spell.Name))
			ctx.Output.WriteLine(fmt.Sprintf("      Mana: %d | Cooldown: %ds | Proficiency: %d%%",
				spell.ManaCost, spell.CooldownSeconds, progress.Proficiency))
			ctx.Output.WriteLine(fmt.Sprintf("      Casts: %d", progress.LifetimeCasts))
			ctx.Output.WriteLine("")
		}
	}

	// Display maneuvers
	if len(maneuversList) > 0 {
		ctx.Output.WriteLine("MANEUVERS:")
		for _, spell := range maneuversList {
			progress := p.Skills[spell.ID]
			ctx.Output.WriteLine(fmt.Sprintf("  [%d] %s", spell.ID, spell.Name))
			ctx.Output.WriteLine(fmt.Sprintf("      Cooldown: %ds | Proficiency: %d%%",
				spell.CooldownSeconds, progress.Proficiency))
			ctx.Output.WriteLine(fmt.Sprintf("      Uses: %d", progress.LifetimeCasts))
			ctx.Output.WriteLine("")
		}
	}

	if len(spellsList) == 0 && len(maneuversList) == 0 {
		ctx.Output.WriteLine("You haven't learned any abilities yet.")
		ctx.Output.WriteLine("")
	}
}

func cmdCast(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in to cast spells")
		return
	}

	args = strings.TrimSpace(args)
	if args == "" {
		ctx.Output.WriteLine("Cast what? (syntax: cast <spell> or cast <spell> <target>)")
		return
	}

	p := ctx.Player

	// Initialize skills map if needed
	if p.Skills == nil {
		p.Skills = make(map[int]*skills.PlayerSkillProgress)
	}

	// Parse spell name and target
	// Try to find the longest matching spell name
	var spell *skills.Spell
	var targetKeyword string
	var skillProgress *skills.PlayerSkillProgress
	var hasSkill bool

	// Try matching from longest to shortest spell name
	words := strings.Fields(args)
	for i := len(words); i > 0; i-- {
		potentialSpellName := strings.Join(words[:i], " ")
		spell = skills.GetSpellByName(potentialSpellName)

		if spell != nil {
			skillProgress, hasSkill = p.Skills[spell.ID]
			if hasSkill && skillProgress.Learned {
				// Found a matching learned spell
				if i < len(words) {
					targetKeyword = strings.Join(words[i:], " ")
				}
				break
			}
			spell = nil // Reset if not learned
		}
	}

	if spell == nil {
		ctx.Output.WriteLine(fmt.Sprintf("You don't know any spell called '%s'", args))
		return
	}

	// Check if player can cast (mana, cooldown)
	now := time.Now().UnixNano()
	canCast, reason := skillProgress.CanCast(spell, p.Mana, now)
	if !canCast {
		ctx.Output.WriteLine(reason)
		return
	}

	// Handle targeting based on spell type
	var targetMob *game.Mobile
	needsTarget := spell.Targeting.Mode == "hostile_single" || spell.Targeting.Mode == "hostile_area"

	if needsTarget {
		if targetKeyword == "" {
			ctx.Output.WriteLine("Cast on whom?")
			return
		}

		// Find target mob by keyword
		mob, found := ctx.World.FindMobInRoom(p, targetKeyword)
		if !found {
			ctx.Output.WriteLine(fmt.Sprintf("You don't see '%s' here.", targetKeyword))
			return
		}
		targetMob = mob
	}

	// Cast the spell!
	p.Mana -= spell.ManaCost
	skillProgress.UpdateCooldown()
	skillProgress.UpdateProficiency(1) // +1% proficiency per cast

	// Calculate damage if it's a damage spell
	totalDamage := 0
	if spell.Effects.Damage != "" && spell.Effects.Damage != "0" {
		// Parse damage formula: "1d6 + I/2" means 1d6 + INT/2
		damageFormula := spell.Effects.Damage

		// Roll base dice (e.g., "1d6")
		parts := strings.Split(damageFormula, "+")
		baseDamage := 0
		if len(parts) > 0 {
			dicePart := strings.TrimSpace(parts[0])
			baseDamage = rollDice(dicePart)
		}

		// Add attribute bonus (INT for spells)
		intBonus := p.Attributes[1] / 2 // INT is index 1
		totalDamage = baseDamage + intBonus

		// Add proficiency scaling
		proficiencyBonus := skillProgress.Proficiency / 20
		totalDamage += proficiencyBonus

		// Deal damage to target
		if targetMob != nil {
			died := ctx.World.DamageMob(p, targetMob, totalDamage)

			// Show messages
			msg := strings.ReplaceAll(spell.Messages.Cast, "$actor", p.Name)
			msg = strings.ReplaceAll(msg, "$spell", spell.Name)
			msg = strings.ReplaceAll(msg, "$target", targetMob.Short)
			msg = strings.TrimSuffix(msg, ".") // Remove trailing period before appending damage
			ctx.Output.WriteLine(fmt.Sprintf("%s for &R%d&w damage! (Proficiency: %d%%)",
				msg, totalDamage, skillProgress.Proficiency))
			ctx.Output.WriteLine(fmt.Sprintf("Mana remaining: %d/%d", p.Mana, p.MaxMana))

			// Broadcast to room
			roomMsg := strings.ReplaceAll(spell.Messages.Cast, "$actor", p.Name)
			roomMsg = strings.ReplaceAll(roomMsg, "$target", targetMob.Short)
			ctx.World.BroadcastCombatMessage(p, roomMsg)

			if died {
				deathMsg := fmt.Sprintf("&R%s falls to the ground, defeated!&w", targetMob.Short)
				ctx.Output.WriteLine(deathMsg)
				ctx.World.BroadcastCombatMessage(p, deathMsg)
			} else {
				hpMsg := fmt.Sprintf("%s has &Y%d/%d&w HP remaining.", targetMob.Short, targetMob.HP, targetMob.MaxHP)
				ctx.Output.WriteLine(hpMsg)
			}
		}
	} else {
		// Non-damage spell (utility, healing, etc.)
		msg := strings.ReplaceAll(spell.Messages.Cast, "$actor", p.Name)
		msg = strings.ReplaceAll(msg, "$spell", spell.Name)
		if targetMob != nil {
			msg = strings.ReplaceAll(msg, "$target", targetMob.Short)
		}
		ctx.Output.WriteLine(fmt.Sprintf("%s (Proficiency: %d%%)", msg, skillProgress.Proficiency))
		ctx.Output.WriteLine(fmt.Sprintf("Mana remaining: %d/%d", p.Mana, p.MaxMana))
	}
}

// rollDice rolls XdY dice (e.g., "1d6" rolls 1 six-sided die)
func rollDice(notation string) int {
	parts := strings.Split(notation, "d")
	if len(parts) != 2 {
		return 0
	}

	var numDice, dieSize int
	fmt.Sscanf(parts[0], "%d", &numDice)
	fmt.Sscanf(parts[1], "%d", &dieSize)

	total := 0
	for i := 0; i < numDice; i++ {
		total += rand.Intn(dieSize) + 1
	}
	return total
}

func cmdSlash(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in to use combat maneuvers")
		return
	}

	args = strings.TrimSpace(args)
	if args == "" {
		ctx.Output.WriteLine("Slash what? (syntax: slash <target>)")
		return
	}

	p := ctx.Player

	// Initialize skills map if needed
	if p.Skills == nil {
		p.Skills = make(map[int]*skills.PlayerSkillProgress)
	}

	// Get the Slash maneuver (skill ID 9002)
	const SlashSkillID = 9002
	spell := skills.GetSpell(SlashSkillID)
	if spell == nil {
		ctx.Output.WriteLine("Slash maneuver not found in skill system.")
		return
	}

	// Check if player has learned Slash
	skillProgress, hasSkill := p.Skills[SlashSkillID]
	if !hasSkill || !skillProgress.Learned {
		ctx.Output.WriteLine("You haven't learned the Slash maneuver yet.")
		return
	}

	// Check cooldown (no mana cost for physical maneuvers)
	now := time.Now().UnixNano()
	canUse, reason := skillProgress.CanCast(spell, p.Mana, now)
	if !canUse {
		// Override mana message for physical attacks since mana cost is 0
		if strings.Contains(reason, "mana") {
			canUse = true
		} else {
			ctx.Output.WriteLine(reason)
			return
		}
	}

	// Find target mob in current room
	targetKeyword := strings.ToLower(args)
	mob, found := ctx.World.FindMobInRoom(p, targetKeyword)
	if !found {
		ctx.Output.WriteLine(fmt.Sprintf("You don't see '%s' here.", args))
		return
	}

	// Calculate damage: 1d6 + STR/2
	// Damage formula from spell.Effects.Damage: "1d6 + S/2"
	baseDamage := rollDice("1d6")
	strBonus := p.Attributes[0] / 2 // STR is index 0
	totalDamage := baseDamage + strBonus

	// Apply proficiency scaling (higher proficiency = more consistent/higher damage)
	// For MVP, proficiency adds a small bonus: +1 damage per 20% proficiency
	proficiencyBonus := skillProgress.Proficiency / 20
	totalDamage += proficiencyBonus

	// Update cooldown and proficiency
	skillProgress.UpdateCooldown()
	skillProgress.UpdateProficiency(1) // +1% proficiency per use

	// Deal damage to mob
	died := ctx.World.DamageMob(p, mob, totalDamage)

	// Show messages
	playerMsg := fmt.Sprintf("You slash at %s for &R%d&w damage! (Proficiency: %d%%)",
		mob.Short, totalDamage, skillProgress.Proficiency)
	ctx.Output.WriteLine(playerMsg)

	roomMsg := fmt.Sprintf("%s slashes at %s!", p.Name, mob.Short)
	ctx.World.BroadcastCombatMessage(p, roomMsg)

	if died {
		deathMsg := fmt.Sprintf("&R%s falls to the ground, defeated!&w", mob.Short)
		ctx.Output.WriteLine(deathMsg)
		ctx.World.BroadcastCombatMessage(p, deathMsg)
	} else {
		hpMsg := fmt.Sprintf("%s has &Y%d/%d&w HP remaining.", mob.Short, mob.HP, mob.MaxHP)
		ctx.Output.WriteLine(hpMsg)
	}
}

func cmdStudy(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in to study")
		return
	}

	args = strings.TrimSpace(args)
	if args == "" {
		ctx.Output.WriteLine("Study what? (syntax: study <item>)")
		return
	}

	p := ctx.Player

	// Initialize skills map if needed
	if p.Skills == nil {
		p.Skills = make(map[int]*skills.PlayerSkillProgress)
	}

	// For MVP: map item keywords to spell IDs
	// In the future, this should be stored in object data
	itemToSpell := map[string]int{
		"arcane":    1001, // Arcane Bolt
		"bolt":      1001,
		"leviathan": 1002, // Leviathan's Fire
		"fire":      1002,
		"mend":      1003, // Mend
		"heal":      1003,
		"healing":   1003,
		"scroll":    1003, // default to heal for scrolls
		"shadow":    1004, // Shadow Veil
		"veil":      1004,
		"ephemeral": 1005, // Ephemeral Step
		"step":      1005,
		"path":      1006, // Path Shift
		"shift":     1006,
		"winter":    1007, // Winter's Whisper
		"whisper":   1007,
		"frost":     1007,
		"cold":      1007,
		"knowing":   1008, // Knowing
		"knowledge": 1008,
	}

	// Find which spell the player is trying to study
	var targetSpellID int
	found := false
	for keyword, spellID := range itemToSpell {
		if strings.Contains(strings.ToLower(args), keyword) {
			targetSpellID = spellID
			found = true
			break
		}
	}

	if !found {
		ctx.Output.WriteLine("You can't study that.")
		return
	}

	// Get the spell
	spell := skills.GetSpell(targetSpellID)
	if spell == nil {
		ctx.Output.WriteLine("That spell doesn't exist.")
		return
	}

	// Check if already learned
	if progress, ok := p.Skills[targetSpellID]; ok && progress.Learned {
		ctx.Output.WriteLine(fmt.Sprintf("You already know %s!", spell.Name))
		return
	}

	// Simulate finding and studying the item
	// For MVP, just auto-succeed with 30% proficiency
	// In full implementation, would:
	// 1. Check if item is in room
	// 2. Make proficiency check (DC = 55 - proficiency*0.8)
	// 3. Remove item from room

	p.Skills[targetSpellID] = &skills.PlayerSkillProgress{
		SpellID:       targetSpellID,
		Proficiency:   30,
		Learned:       true,
		LifetimeCasts: 0,
		LastCastTime:  0,
	}

	ctx.Output.WriteLine(fmt.Sprintf("&YYou carefully study the item and learn &W%s&Y!&w", spell.Name))
	ctx.Output.WriteLine(fmt.Sprintf("Proficiency: 30%% | Mana Cost: %d | Cooldown: %ds",
		spell.ManaCost, spell.CooldownSeconds))

	// In full implementation: ctx.World.RemoveObjectFromRoom(ctx.Player.Location, itemVnum)
}

func cmdSave(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in to save")
		return
	}

	record := persist.PlayerToRecord(ctx.Player)
	if err := persist.SavePlayer("players", record); err != nil {
		ctx.Output.WriteLine(fmt.Sprintf("&RError saving: %v&w", err))
		return
	}

	ctx.Output.WriteLine("&YYour progress has been saved.&w")
}

func cmdQuit(ctx Context, args string) {
	ctx.Output.WriteLine("Goodbye.")
	if ctx.Disconnect != nil {
		ctx.Disconnect("quit")
	}
}

func cmdExits(ctx Context, args string) {
	view, err := ctx.World.DescribeRoom(ctx.Player)
	if err != nil {
		ctx.Output.WriteLine("You are nowhere.")
		return
	}
	ctx.Output.WriteLine(FormatExits(view.Exits))
}

func cmdAutoexits(ctx Context, args string) {
	ctx.Player.AutoExits = !ctx.Player.AutoExits
	if ctx.Player.AutoExits {
		ctx.Output.WriteLine("Autoexits enabled.")
		return
	}
	ctx.Output.WriteLine("Autoexits disabled.")
}

func cmdAstat(ctx Context, args string) {
	view, err := ctx.World.DescribeRoom(ctx.Player)
	if err != nil {
		ctx.Output.WriteLine("You are nowhere.")
		return
	}

	ctx.Output.WriteLine("")
	ctx.Output.WriteLine("=== AREA STATISTICS ===")
	ctx.Output.WriteLine("")
	ctx.Output.WriteLine(fmt.Sprintf("Name:   %s", view.AreaName))
	ctx.Output.WriteLine(fmt.Sprintf("Author: %s", view.AreaAuthor))
	ctx.Output.WriteLine("")
}

func FormatExits(exits []string) string {
	if len(exits) == 0 {
		return "Exits: none"
	}

	order := []string{"north", "east", "south", "west", "up", "down", "northeast", "northwest", "southeast", "southwest"}

	present := map[string]bool{}
	for _, exit := range exits {
		present[exit] = true
	}

	display := make([]string, 0, len(exits))
	for _, key := range order {
		if present[key] {
			display = append(display, capitalize(key))
			delete(present, key)
		}
	}

	return "Exits: " + strings.Join(display, " ") + "."
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// Keeper Commands

func cmdMakeKeeper(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in.")
		return
	}

	if !ctx.Player.IsKeeper {
		ctx.Output.WriteLine("You do not have the authority to do that.")
		return
	}

	targetName := strings.TrimSpace(args)
	if targetName == "" {
		ctx.Output.WriteLine("Usage: makekeeper <player name>")
		return
	}

	// Find the player in the world
	target, ok := ctx.World.FindPlayer(targetName)
	if !ok {
		ctx.Output.WriteLine(fmt.Sprintf("Player '%s' not found.", targetName))
		return
	}

	if target.IsKeeper {
		ctx.Output.WriteLine(fmt.Sprintf("%s is already a Keeper of the realm.", target.Name))
		return
	}

	target.IsKeeper = true
	ctx.Output.WriteLine(fmt.Sprintf("You grant %s the responsibility of a Keeper.", target.Name))
	target.Output.WriteLine("You have been elevated to Keeper of the realm. Guard this position well.")
}

func cmdRemoveKeeper(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in.")
		return
	}

	if !ctx.Player.IsKeeper {
		ctx.Output.WriteLine("You do not have the authority to do that.")
		return
	}

	targetName := strings.TrimSpace(args)
	if targetName == "" {
		ctx.Output.WriteLine("Usage: removekeeper <player name>")
		return
	}

	// Find the player in the world
	target, ok := ctx.World.FindPlayer(targetName)
	if !ok {
		ctx.Output.WriteLine(fmt.Sprintf("Player '%s' not found.", targetName))
		return
	}

	if !target.IsKeeper {
		ctx.Output.WriteLine(fmt.Sprintf("%s is not a Keeper.", target.Name))
		return
	}

	target.IsKeeper = false
	ctx.Output.WriteLine(fmt.Sprintf("You strip %s of their Keeper responsibilities.", target.Name))
	target.Output.WriteLine("Your status as a Keeper has been revoked.")
}

func cmdSpawn(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in.")
		return
	}

	if !ctx.Player.IsKeeper {
		ctx.Output.WriteLine("You do not have the authority to do that.")
		return
	}

	vnumStr := strings.TrimSpace(args)
	if vnumStr == "" {
		ctx.Output.WriteLine("Usage: spawn <mob vnum>")
		ctx.Output.WriteLine("Example: spawn 90001")
		return
	}

	vnum, err := strconv.Atoi(vnumStr)
	if err != nil {
		ctx.Output.WriteLine("Invalid vnum. Must be a number.")
		return
	}

	// Spawn the mob
	mob, err := ctx.World.SpawnMob(ctx.Player, vnum)
	if err != nil {
		ctx.Output.WriteLine(fmt.Sprintf("Failed to spawn mob: %s", err.Error()))
		return
	}

	ctx.Output.WriteLine(fmt.Sprintf("&GYou summon %s into existence!&w", mob.Short))
	ctx.World.BroadcastCombatMessage(ctx.Player, fmt.Sprintf("%s summons %s into existence!", ctx.Player.Name, mob.Short))
}

func cmdTeleport(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in.")
		return
	}

	if !ctx.Player.IsKeeper {
		ctx.Output.WriteLine("You do not have the authority to do that.")
		return
	}

	// Parse room vnum from args
	args = strings.TrimSpace(args)
	if args == "" {
		ctx.Output.WriteLine("Usage: teleport <room vnum>")
		return
	}

	var vnum int
	_, err := fmt.Sscanf(args, "%d", &vnum)
	if err != nil || vnum <= 0 {
		ctx.Output.WriteLine("Invalid room vnum.")
		return
	}

	if !ctx.World.HasRoom(vnum) {
		ctx.Output.WriteLine(fmt.Sprintf("Room %d does not exist.", vnum))
		return
	}

	ctx.Player.Location = vnum
	ctx.Output.WriteLine(fmt.Sprintf("You teleport to room %d.", vnum))

	// Show the room
	view, err := ctx.World.DescribeRoom(ctx.Player)
	if err == nil {
		DisplayRoomView(ctx.Output, view, ctx.Player.AutoExits)
	}
}

func registerMovement(registry *Registry) {
	directions := map[string]string{
		"north": "north",
		"south": "south",
		"east":  "east",
		"west":  "west",
		"up":    "up",
		"down":  "down",
		"ne":    "northeast",
		"nw":    "northwest",
		"se":    "southeast",
		"sw":    "southwest",
		"n":     "north",
		"s":     "south",
		"e":     "east",
		"w":     "west",
		"u":     "up",
		"d":     "down",
	}

	for name, direction := range directions {
		dir := direction
		registry.Register(name, func(ctx Context, args string) {
			view, err := ctx.World.MovePlayer(ctx.Player, dir)
			if err != nil {
				ctx.Output.WriteLine("You cannot go that way.")
				return
			}

			DisplayRoomView(ctx.Output, view, ctx.Player.AutoExits)
		})
	}
}
