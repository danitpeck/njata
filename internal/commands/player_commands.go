package commands

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
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
	registry.Register("consider", cmdConsider)
	registry.Register("say", cmdSay)
	registry.Register("chat", cmdChat)
	registry.Register("who", cmdWho)
	registry.Register("stats", cmdStats)
	registry.Register("inventory", cmdInventory)
	registry.Register("inv", cmdInventory)
	registry.Register("equipment", cmdEquipment)
	registry.Register("wear", cmdWear)
	registry.Register("remove", cmdRemove)
	registry.Register("get", cmdGet)
	registry.Register("drop", cmdDrop)
	registry.Register("hair", cmdHair)
	registry.Register("eyes", cmdEyes)
	registry.Register("exits", cmdExits)
	registry.Register("autoexits", cmdAutoexits)
	registry.Register("astat", cmdAstat)
	registry.Register("abilities", cmdAbilities)
	registry.Register("cast", cmdCast)
	registry.Register("slash", cmdSlash)
	registry.Register("power", cmdPowerAttack)
	registry.Register("powerattack", cmdPowerAttack)
	registry.Register("riposte", cmdRiposte)
	registry.Register("cleave", cmdCleave)
	registry.Register("defensive", cmdDefensiveStance)
	registry.Register("defensivestance", cmdDefensiveStance)
	registry.Register("study", cmdStudy)
	registry.Register("train", cmdTrain)
	registry.Register("save", cmdSave)
	registry.Register("makekeeper", cmdMakeKeeper)
	registry.Register("removekeeper", cmdRemoveKeeper)
	registry.Register("spawn", cmdSpawn)
	registry.Register("teleport", cmdTeleport)
	registry.Register("restore", cmdRestore)
	registry.Register("help", cmdHelp)
	registry.Register("quit", cmdQuit)
	registerMovement(registry)
}

func cmdConsider(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in.")
		return
	}

	args = strings.TrimSpace(args)
	if args == "" {
		ctx.Output.WriteLine("Consider whom? (syntax: consider <target>)")
		return
	}

	if strings.EqualFold(args, "me") || strings.EqualFold(args, "self") {
		ctx.Output.WriteLine("You look capable enough to handle yourself.")
		return
	}

	if targetPlayer, ok := ctx.World.FindPlayerInRoom(ctx.Player, args); ok {
		assessment := compareCombatScores(combatScorePlayer(ctx.Player), combatScorePlayer(targetPlayer))
		ctx.Output.WriteLine(fmt.Sprintf("You size up %s. %s", game.CapitalizeName(targetPlayer.Name), assessment))
		return
	}

	mob, ok := ctx.World.FindMobInRoom(ctx.Player, args)
	if !ok {
		ctx.Output.WriteLine(fmt.Sprintf("You don't see '%s' here.", args))
		return
	}

	assessment := compareCombatScores(combatScorePlayer(ctx.Player), combatScoreMob(mob))
	ctx.Output.WriteLine(fmt.Sprintf("You size up %s. %s", mob.Short, assessment))
}

func combatScorePlayer(p *game.Player) int {
	if p == nil {
		return 1
	}

	score := 0
	score += p.MaxHP
	score += p.Strength * 5
	score += p.Dexterity * 3
	score += p.Constitution * 4
	score += p.Armor * 2
	if score < 1 {
		score = 1
	}
	return score
}

func combatScoreMob(mob *game.Mobile) int {
	if mob == nil {
		return 1
	}

	maxHP := mob.MaxHP
	if maxHP <= 0 {
		maxHP = mob.HP
	}

	score := 0
	score += maxHP
	score += mob.Attributes[0] * 5
	score += mob.Attributes[3] * 3
	score += mob.Attributes[4] * 4
	score += mob.Level * 10
	if score < 1 {
		score = 1
	}
	return score
}

func compareCombatScores(playerScore int, targetScore int) string {
	if playerScore < 1 {
		playerScore = 1
	}
	if targetScore < 1 {
		targetScore = 1
	}

	ratio := float64(targetScore) / float64(playerScore)
	switch {
	case ratio <= 0.5:
		return "It looks trivial."
	case ratio <= 0.8:
		return "You would have the advantage."
	case ratio <= 1.2:
		return "It seems evenly matched."
	case ratio <= 1.6:
		return "It looks dangerous."
	default:
		return "You would likely be defeated."
	}
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

		if strings.EqualFold(keyword, "me") || strings.EqualFold(keyword, "self") {
			writePlayerLook(ctx.Output, ctx.Player)
			return
		}

		if target, ok := ctx.World.FindPlayerInRoom(ctx.Player, keyword); ok {
			writePlayerLook(ctx.Output, target)
			return
		}

		if mob, ok := ctx.World.FindMobInRoom(ctx.Player, keyword); ok {
			ctx.Output.WriteLine(mob.Long)
			return
		}

		if obj, ok := ctx.World.FindObjectInRoom(ctx.Player, keyword); ok {
			ctx.Output.WriteLine(obj.Long)
			return
		}

		if obj, ok := ctx.World.FindObjectInInventory(ctx.Player, keyword); ok {
			ctx.Output.WriteLine(obj.Long)
			return
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

func writePlayerLook(output game.Output, target *game.Player) {
	name := game.CapitalizeName(target.Name)
	output.WriteLine(fmt.Sprintf("You see %s.", name))

	hair := target.Hair
	if hair == "" {
		hair = "(none)"
	}

	eyes := target.Eyes
	if eyes == "" {
		eyes = "(none)"
	}

	output.WriteLine(fmt.Sprintf("Hair: %s", hair))
	output.WriteLine(fmt.Sprintf("Eyes: %s", eyes))
}

func cmdSay(ctx Context, args string) {
	if strings.TrimSpace(args) == "" {
		ctx.Output.WriteLine("Say what?")
		return
	}

	ctx.World.BroadcastSay(ctx.Player, args)
}

func cmdChat(ctx Context, args string) {
	if strings.TrimSpace(args) == "" {
		ctx.Output.WriteLine("Chat what?")
		return
	}

	ctx.World.BroadcastChat(ctx.Player, args)
}

func cmdWho(ctx Context, args string) {
	players := ctx.World.ListPlayers()
	ctx.Output.WriteLine(fmt.Sprintf("Players online (%d): %s", len(players), strings.Join(players, ", ")))
}

func cmdStats(ctx Context, args string) {
	p := ctx.Player
	ctx.Output.WriteLine(fmt.Sprintf("=== %s ===", game.CapitalizeName(p.Name)))

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
	ctx.Output.WriteLine(fmt.Sprintf("HP:    %d/%d | Mana: %d/%d", p.HP, p.MaxHP, p.Mana, p.MaxMana))
	ctx.Output.WriteLine(fmt.Sprintf("Gold: %d", p.Gold))
	ctx.Output.WriteLine("")

	ctx.Output.WriteLine(fmt.Sprintf("STR: %2d | DEX: %2d | CON: %2d", p.Strength, p.Dexterity, p.Constitution))
	ctx.Output.WriteLine(fmt.Sprintf("INT: %2d | WIS: %2d | CHA: %2d | LCK: %2d", p.Intelligence, p.Wisdom, p.Charisma, p.Luck))

	ctx.Output.WriteLine("")
	ctx.Output.WriteLine(fmt.Sprintf("Armor: %d", p.Armor))
}

func cmdInventory(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in to view inventory.")
		return
	}

	if len(ctx.Player.Inventory) == 0 {
		ctx.Output.WriteLine("You are carrying nothing.")
		return
	}

	ctx.Output.WriteLine("You are carrying:")
	for _, item := range ctx.Player.Inventory {
		if item == nil {
			continue
		}
		label := item.Short
		if label == "" {
			label = "something"
		}
		ctx.Output.WriteLine("  " + label)
	}
}

func cmdEquipment(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in to view equipment.")
		return
	}

	equipment := ctx.World.EquipmentSnapshot(ctx.Player)
	itemCount := 0
	for _, obj := range equipment {
		if obj != nil {
			itemCount++
		}
	}
	if itemCount == 0 {
		ctx.Output.WriteLine("You are wearing nothing.")
		return
	}

	ctx.Output.WriteLine("You are wearing:")
	for _, slot := range game.EquipSlotOrder {
		obj := equipment[slot]
		if obj == nil {
			continue
		}
		label := obj.Short
		if label == "" {
			label = "something"
		}
		ctx.Output.WriteLine(fmt.Sprintf("  %-5s %s", slot+":", label))
	}
}

func cmdWear(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in to wear items.")
		return
	}

	trimmed := strings.TrimSpace(args)
	if strings.EqualFold(trimmed, "all") {
		items := append([]*game.Object(nil), ctx.Player.Inventory...)
		worn := 0
		for _, obj := range items {
			if obj == nil {
				continue
			}
			slot, ok := resolveEquipSlot(obj, "")
			if !ok {
				continue
			}
			equipment := ctx.World.EquipmentSnapshot(ctx.Player)
			if equipment[slot] != nil {
				continue
			}
			if ctx.World.EquipObject(ctx.Player, obj, slot) {
				// Apply armor bonus
				if obj.ArmorVal != 0 {
					ctx.Player.Armor += obj.ArmorVal
				}
				label := obj.Short
				if label == "" {
					label = "something"
				}
				if slot == game.EquipWield {
					ctx.Output.WriteLine(fmt.Sprintf("You equip %s.", label))
				} else {
					ctx.Output.WriteLine(fmt.Sprintf("You wear %s on your %s.", label, slot))
				}
				worn++
			}
		}
		if worn == 0 {
			ctx.Output.WriteLine("You have nothing you can wear.")
			return
		}
		return
	}

	fields := strings.Fields(trimmed)
	if len(fields) == 0 {
		ctx.Output.WriteLine("Wear what?")
		return
	}

	keyword := fields[0]
	var slotOverride string
	if len(fields) > 1 {
		slotOverride = fields[1]
	}

	obj, found := ctx.World.FindObjectInInventory(ctx.Player, keyword)
	if !found {
		ctx.Output.WriteLine("You aren't carrying that.")
		return
	}

	slot, ok := resolveEquipSlot(obj, slotOverride)
	if !ok {
		ctx.Output.WriteLine("You can't wear that.")
		return
	}

	equipment := ctx.World.EquipmentSnapshot(ctx.Player)
	if equipment[slot] != nil {
		ctx.Output.WriteLine(fmt.Sprintf("You are already wearing something on your %s.", slot))
		return
	}

	if !ctx.World.EquipObject(ctx.Player, obj, slot) {
		ctx.Output.WriteLine("You can't wear that.")
		return
	}

	// Apply armor bonus
	if obj.ArmorVal != 0 {
		ctx.Player.Armor += obj.ArmorVal
	}

	label := obj.Short
	if label == "" {
		label = "something"
	}
	if slot == game.EquipWield {
		ctx.Output.WriteLine(fmt.Sprintf("You equip %s.", label))
	} else {
		ctx.Output.WriteLine(fmt.Sprintf("You wear %s on your %s.", label, slot))
	}
}

func cmdRemove(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in to remove items.")
		return
	}

	trimmed := strings.TrimSpace(args)
	if strings.EqualFold(trimmed, "all") {
		removed := 0
		for _, slot := range game.EquipSlotOrder {
			obj, ok := ctx.World.UnequipObject(ctx.Player, slot)
			if !ok {
				continue
			}
			// Remove armor bonus
			if obj.ArmorVal != 0 {
				ctx.Player.Armor -= obj.ArmorVal
			}
			label := obj.Short
			if label == "" {
				label = "something"
			}
			ctx.Output.WriteLine(fmt.Sprintf("You remove %s from your %s.", label, slot))
			removed++
		}
		if removed == 0 {
			ctx.Output.WriteLine("You are wearing nothing.")
			return
		}
		return
	}

	fields := strings.Fields(trimmed)
	if len(fields) == 0 {
		ctx.Output.WriteLine("Remove what?")
		return
	}

	keyword := fields[0]
	if slot, ok := normalizeEquipSlot(keyword); ok {
		obj, ok := ctx.World.UnequipObject(ctx.Player, slot)
		if !ok {
			ctx.Output.WriteLine(fmt.Sprintf("You are not wearing anything on your %s.", slot))
			return
		}
		// Remove armor bonus
		if obj.ArmorVal != 0 {
			ctx.Player.Armor -= obj.ArmorVal
		}
		label := obj.Short
		if label == "" {
			label = "something"
		}
		ctx.Output.WriteLine(fmt.Sprintf("You remove %s from your %s.", label, slot))
		return
	}

	obj, found := ctx.World.FindObjectInEquipment(ctx.Player, keyword)
	if !found {
		ctx.Output.WriteLine("You are not wearing that.")
		return
	}

	slot, ok := ctx.World.FindEquippedSlot(ctx.Player, obj)
	if !ok {
		ctx.Output.WriteLine("You are not wearing that.")
		return
	}

	if _, ok := ctx.World.UnequipObject(ctx.Player, slot); !ok {
		ctx.Output.WriteLine("You are not wearing that.")
		return
	}

	// Remove armor bonus
	if obj.ArmorVal != 0 {
		ctx.Player.Armor -= obj.ArmorVal
	}

	label := obj.Short
	if label == "" {
		label = "something"
	}
	ctx.Output.WriteLine(fmt.Sprintf("You remove %s from your %s.", label, slot))
}

func normalizeEquipSlot(slot string) (string, bool) {
	trimmed := strings.ToLower(strings.TrimSpace(slot))
	if trimmed == "" {
		return "", false
	}

	switch trimmed {
	case game.EquipHead:
		return game.EquipHead, true
	case game.EquipBody:
		return game.EquipBody, true
	case game.EquipNeck:
		return game.EquipNeck, true
	case game.EquipBack:
		return game.EquipBack, true
	case game.EquipWaist:
		return game.EquipWaist, true
	default:
		return "", false
	}
}

func resolveEquipSlot(obj *game.Object, slotOverride string) (string, bool) {
	if obj == nil {
		return "", false
	}

	objType := strings.ToLower(strings.TrimSpace(obj.Type))
	if objType != "armor" && objType != "_worn" && objType != "weapon" {
		return "", false
	}

	// Prefer explicit equip_slot field if set
	var slotToUse string
	if obj.EquipSlot != "" {
		slotToUse = obj.EquipSlot
	} else {
		slotToUse = inferEquipSlot(obj)
	}

	if slotOverride == "" {
		return slotToUse, true
	}

	slot, ok := normalizeEquipSlot(slotOverride)
	if !ok {
		return "", false
	}

	if slot != slotToUse {
		return "", false
	}

	return slot, true
}

func inferEquipSlot(obj *game.Object) string {
	if obj == nil {
		return game.EquipBody
	}

	// Check for weapons
	if strings.ToLower(strings.TrimSpace(obj.Type)) == "weapon" {
		return game.EquipWield
	}

	if objectHasAnyKeyword(obj, []string{"helm", "helmet", "hat", "cap", "hood", "mask", "crown", "circlet", "tiara", "coif", "bonnet", "earmuff", "earmuffs", "glasses", "goggles"}) {
		return game.EquipHead
	}
	if objectHasAnyKeyword(obj, []string{"neck", "necklace", "amulet", "pendant", "torc", "choker", "collar", "gorget", "scarf"}) {
		return game.EquipNeck
	}
	if objectHasAnyKeyword(obj, []string{"cloak", "cape", "back", "mantle", "shawl", "poncho"}) {
		return game.EquipBack
	}
	if objectHasAnyKeyword(obj, []string{"belt", "waist", "girdle", "sash", "cincture"}) {
		return game.EquipWaist
	}

	return game.EquipBody
}

func objectHasAnyKeyword(obj *game.Object, keywords []string) bool {
	if obj == nil {
		return false
	}

	shortText := strings.ToLower(obj.Short)
	for _, keyword := range keywords {
		key := strings.ToLower(keyword)
		for _, objKeyword := range obj.Keywords {
			if strings.Contains(strings.ToLower(objKeyword), key) {
				return true
			}
		}
		if strings.Contains(shortText, key) {
			return true
		}
	}

	return false
}

func objectMatchesKeyword(obj *game.Object, keyword string) bool {
	if obj == nil {
		return false
	}

	key := strings.ToLower(strings.TrimSpace(keyword))
	if key == "" {
		return false
	}

	for _, objKeyword := range obj.Keywords {
		if strings.ToLower(objKeyword) == key {
			return true
		}
	}

	return strings.Contains(strings.ToLower(obj.Short), key)
}

func cmdGet(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in to get items.")
		return
	}

	keyword := strings.TrimSpace(args)
	if keyword == "" {
		ctx.Output.WriteLine("Get what?")
		return
	}

	lower := strings.ToLower(keyword)
	if lower == "all" || strings.HasPrefix(lower, "all ") || strings.HasPrefix(lower, "all.") {
		filter := ""
		if lower != "all" {
			filter = strings.TrimSpace(keyword[4:])
		}

		objects := ctx.World.RoomObjectsSnapshot(ctx.Player)
		picked := 0
		for _, obj := range objects {
			if obj == nil {
				continue
			}
			if filter != "" && !objectMatchesKeyword(obj, filter) {
				continue
			}
			if obj.Flags != nil && (obj.Flags["notake"] || obj.Flags["no_take"]) {
				continue
			}
			if obj.Type == "fountain" || obj.Type == "furniture" {
				continue
			}
			if !ctx.World.RemoveObjectFromRoom(ctx.Player, obj) {
				continue
			}
			ctx.World.AddObjectToInventory(ctx.Player, obj)
			label := obj.Short
			if label == "" {
				label = "something"
			}
			ctx.Output.WriteLine(fmt.Sprintf("You pick up %s.", label))
			picked++
		}

		if picked == 0 {
			ctx.Output.WriteLine("You see nothing here.")
		}
		return
	}

	obj, found := ctx.World.FindObjectInRoom(ctx.Player, keyword)
	if !found {
		ctx.Output.WriteLine("You don't see that here.")
		return
	}

	if obj.Flags != nil && (obj.Flags["notake"] || obj.Flags["no_take"]) {
		ctx.Output.WriteLine("You can't take that.")
		return
	}

	if obj.Type == "fountain" || obj.Type == "furniture" {
		ctx.Output.WriteLine("You can't take that.")
		return
	}

	if !ctx.World.RemoveObjectFromRoom(ctx.Player, obj) {
		ctx.Output.WriteLine("You can't take that.")
		return
	}

	ctx.World.AddObjectToInventory(ctx.Player, obj)
	label := obj.Short
	if label == "" {
		label = "something"
	}
	ctx.Output.WriteLine(fmt.Sprintf("You pick up %s.", label))
}

func cmdDrop(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in to drop items.")
		return
	}

	keyword := strings.TrimSpace(args)
	if keyword == "" {
		ctx.Output.WriteLine("Drop what?")
		return
	}

	lower := strings.ToLower(keyword)
	if lower == "all" || strings.HasPrefix(lower, "all ") || strings.HasPrefix(lower, "all.") {
		filter := ""
		if lower != "all" {
			filter = strings.TrimSpace(keyword[4:])
		}

		items := append([]*game.Object(nil), ctx.Player.Inventory...)
		dropped := 0
		for _, obj := range items {
			if obj == nil {
				continue
			}
			if filter != "" && !objectMatchesKeyword(obj, filter) {
				continue
			}
			if !ctx.World.RemoveObjectFromInventory(ctx.Player, obj) {
				continue
			}
			ctx.World.AddObjectToRoom(ctx.Player, obj)
			label := obj.Short
			if label == "" {
				label = "something"
			}
			ctx.Output.WriteLine(fmt.Sprintf("You drop %s.", label))
			dropped++
		}

		if dropped == 0 {
			ctx.Output.WriteLine("You are carrying nothing.")
		}
		return
	}

	obj, found := ctx.World.FindObjectInInventory(ctx.Player, keyword)
	if !found {
		ctx.Output.WriteLine("You aren't carrying that.")
		return
	}

	if !ctx.World.RemoveObjectFromInventory(ctx.Player, obj) {
		ctx.Output.WriteLine("You can't drop that.")
		return
	}

	ctx.World.AddObjectToRoom(ctx.Player, obj)
	label := obj.Short
	if label == "" {
		label = "something"
	}
	ctx.Output.WriteLine(fmt.Sprintf("You drop %s.", label))
}

func cmdHair(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in to set your hair description.")
		return
	}

	trimmed := strings.TrimSpace(args)
	if trimmed == "" {
		if ctx.Player.Hair != "" {
			ctx.Output.WriteLine(fmt.Sprintf("Your hair description is: %s", ctx.Player.Hair))
		} else {
			ctx.Output.WriteLine("You have no hair description.")
		}
		return
	}

	if strings.EqualFold(trimmed, "clear") {
		ctx.Player.Hair = ""
		ctx.Output.WriteLine("Hair description cleared.")
		return
	}

	trimmed = strings.ReplaceAll(trimmed, "~", "")
	ctx.Player.Hair = trimmed
	ctx.Output.WriteLine(fmt.Sprintf("Your hair description is: %s", ctx.Player.Hair))
}

func cmdEyes(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in to set your eyes description.")
		return
	}

	trimmed := strings.TrimSpace(args)
	if trimmed == "" {
		if ctx.Player.Eyes != "" {
			ctx.Output.WriteLine(fmt.Sprintf("Your eyes description is: %s", ctx.Player.Eyes))
		} else {
			ctx.Output.WriteLine("You have no eyes description.")
		}
		return
	}

	if strings.EqualFold(trimmed, "clear") {
		ctx.Player.Eyes = ""
		ctx.Output.WriteLine("Eyes description cleared.")
		return
	}

	trimmed = strings.ReplaceAll(trimmed, "~", "")
	ctx.Player.Eyes = trimmed
	ctx.Output.WriteLine(fmt.Sprintf("Your eyes description is: %s", ctx.Player.Eyes))
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

		// Categorize: spells (1000-1999) vs maneuvers (2000-2999)
		if spell.ID >= 2000 && spell.ID < 3000 {
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
		intBonus := p.Intelligence / 2 // INT bonus
		totalDamage = baseDamage + intBonus

		// Add proficiency scaling
		proficiencyBonus := skillProgress.Proficiency / 20
		totalDamage += proficiencyBonus

		// Deal damage to target
		if targetMob != nil {
			died, loot := ctx.World.DamageMob(p, targetMob, totalDamage)

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
				if len(loot) > 0 {
					ctx.Output.WriteLine("You loot: " + strings.Join(loot, ", ") + ".")
				}
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

	// Get the Slash maneuver (skill ID 2001)
	const SlashSkillID = 2001
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
	strBonus := p.Strength / 2 // STR bonus
	totalDamage := baseDamage + strBonus

	// Apply proficiency scaling (higher proficiency = more consistent/higher damage)
	// For MVP, proficiency adds a small bonus: +1 damage per 20% proficiency
	proficiencyBonus := skillProgress.Proficiency / 20
	totalDamage += proficiencyBonus

	// Update cooldown and proficiency
	skillProgress.UpdateCooldown()
	skillProgress.UpdateProficiency(1) // +1% proficiency per use

	// Deal damage to mob
	died, loot := ctx.World.DamageMob(p, mob, totalDamage)

	// Show messages
	playerMsg := fmt.Sprintf("You slash at %s for &R%d&w damage! (Proficiency: %d%%)",
		mob.Short, totalDamage, skillProgress.Proficiency)
	ctx.Output.WriteLine(playerMsg)

	roomMsg := fmt.Sprintf("%s slashes at %s!", game.CapitalizeName(p.Name), mob.Short)
	ctx.World.BroadcastCombatMessage(p, roomMsg)

	if died {
		deathMsg := fmt.Sprintf("&R%s falls to the ground, defeated!&w", mob.Short)
		ctx.Output.WriteLine(deathMsg)
		ctx.World.BroadcastCombatMessage(p, deathMsg)
		if len(loot) > 0 {
			ctx.Output.WriteLine("You loot: " + strings.Join(loot, ", ") + ".")
		}
	} else {
		hpMsg := fmt.Sprintf("%s has &Y%d/%d&w HP remaining.", mob.Short, mob.HP, mob.MaxHP)
		ctx.Output.WriteLine(hpMsg)
	}
}

func cmdPowerAttack(ctx Context, args string) {
	args = strings.TrimSpace(args)
	args = strings.TrimPrefix(args, "attack ")
	performManeuver(ctx, strings.TrimSpace(args), 2002, "power")
}

func cmdRiposte(ctx Context, args string) {
	performManeuver(ctx, strings.TrimSpace(args), 2004, "riposte")
}

func cmdCleave(ctx Context, args string) {
	performManeuver(ctx, strings.TrimSpace(args), 2005, "cleave")
}

func cmdDefensiveStance(ctx Context, args string) {
	args = strings.TrimSpace(args)
	args = strings.TrimPrefix(args, "stance")
	performManeuver(ctx, strings.TrimSpace(args), 2003, "defensive")
}

func performManeuver(ctx Context, args string, skillID int, commandName string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in to use combat maneuvers")
		return
	}

	p := ctx.Player

	if p.Skills == nil {
		p.Skills = make(map[int]*skills.PlayerSkillProgress)
	}

	spell := skills.GetSpell(skillID)
	if spell == nil {
		ctx.Output.WriteLine("That maneuver is not available.")
		return
	}

	skillProgress, hasSkill := p.Skills[skillID]
	if !hasSkill || !skillProgress.Learned {
		ctx.Output.WriteLine(fmt.Sprintf("You haven't learned %s yet.", spell.Name))
		return
	}

	now := time.Now().UnixNano()
	canUse, reason := skillProgress.CanCast(spell, p.Mana, now)
	if !canUse {
		if strings.Contains(reason, "mana") {
			canUse = true
		} else {
			ctx.Output.WriteLine(reason)
			return
		}
	}

	needsTarget := spell.Targeting.Mode == "hostile_single" || spell.Targeting.Mode == "hostile_area"
	if needsTarget && args == "" {
		ctx.Output.WriteLine(fmt.Sprintf("%s whom? (syntax: %s <target>)", spell.Name, commandName))
		return
	}

	var targetMob *game.Mobile
	if needsTarget {
		mob, found := ctx.World.FindMobInRoom(p, args)
		if !found {
			ctx.Output.WriteLine(fmt.Sprintf("You don't see '%s' here.", args))
			return
		}
		targetMob = mob
	}

	skillProgress.UpdateCooldown()
	skillProgress.UpdateProficiency(1)

	if spell.Effects.Damage != "" && spell.Effects.Damage != "0" && targetMob != nil {
		damageFormula := spell.Effects.Damage
		parts := strings.Split(damageFormula, "+")
		baseDamage := 0
		if len(parts) > 0 {
			dicePart := strings.TrimSpace(parts[0])
			baseDamage = rollDice(dicePart)
		}
		statBonus := maneuverStatBonus(parts, p)
		proficiencyBonus := skillProgress.Proficiency / 20
		totalDamage := baseDamage + statBonus + proficiencyBonus

		died, loot := ctx.World.DamageMob(p, targetMob, totalDamage)

		msg := strings.ReplaceAll(spell.Messages.Cast, "$actor", p.Name)
		msg = strings.ReplaceAll(msg, "$target", targetMob.Short)
		msg = strings.TrimSuffix(msg, ".")
		ctx.Output.WriteLine(fmt.Sprintf("%s for &R%d&w damage! (Proficiency: %d%%)",
			msg, totalDamage, skillProgress.Proficiency))

		roomMsg := strings.ReplaceAll(spell.Messages.Cast, "$actor", p.Name)
		roomMsg = strings.ReplaceAll(roomMsg, "$target", targetMob.Short)
		ctx.World.BroadcastCombatMessage(p, roomMsg)

		if died {
			deathMsg := fmt.Sprintf("&R%s falls to the ground, defeated!&w", targetMob.Short)
			ctx.Output.WriteLine(deathMsg)
			ctx.World.BroadcastCombatMessage(p, deathMsg)
			if len(loot) > 0 {
				ctx.Output.WriteLine("You loot: " + strings.Join(loot, ", ") + ".")
			}
		} else {
			hpMsg := fmt.Sprintf("%s has &Y%d/%d&w HP remaining.", targetMob.Short, targetMob.HP, targetMob.MaxHP)
			ctx.Output.WriteLine(hpMsg)
		}

		return
	}

	msg := strings.ReplaceAll(spell.Messages.Cast, "$actor", p.Name)
	if targetMob != nil {
		msg = strings.ReplaceAll(msg, "$target", targetMob.Short)
	}
	ctx.Output.WriteLine(fmt.Sprintf("%s (Proficiency: %d%%)", msg, skillProgress.Proficiency))
}

func maneuverStatBonus(parts []string, p *game.Player) int {
	if len(parts) < 2 {
		return 0
	}

	bonusPart := strings.TrimSpace(parts[1])
	if bonusPart == "" {
		return 0
	}

	statChar := strings.ToUpper(string(bonusPart[0]))
	divisor := 1
	if slashIdx := strings.Index(bonusPart, "/"); slashIdx != -1 {
		if value, err := strconv.Atoi(strings.TrimSpace(bonusPart[slashIdx+1:])); err == nil && value > 0 {
			divisor = value
		}
	}

	statValue := 0
	switch statChar {
	case "S":
		statValue = p.Strength
	case "D":
		statValue = p.Dexterity
	case "C":
		statValue = p.Constitution
	case "I":
		statValue = p.Intelligence
	case "W":
		statValue = p.Wisdom
	case "L":
		statValue = p.Luck
	}

	return statValue / divisor
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

	// Look for the item in inventory first, then in the room
	var obj *game.Object
	var found bool

	obj, found = ctx.World.FindObjectInInventory(p, args)
	if !found {
		obj, found = ctx.World.FindObjectInRoom(p, args)
		if !found {
			ctx.Output.WriteLine("You don't see that here.")
			return
		}
	}

	// Determine spell ID and proficiency amount
	var spellID int
	var teachAmount int

	if obj.TeachesSpellID != 0 {
		// Use new teaching fields
		spellID = obj.TeachesSpellID
		teachAmount = obj.TeachesAmount
		if teachAmount == 0 {
			teachAmount = 30 // default
		}
	} else if obj.Value[3] != 0 {
		// Fall back to legacy Value[3]
		spellID = obj.Value[3]
		teachAmount = 30
	} else {
		ctx.Output.WriteLine("That item doesn't teach any spells.")
		return
	}

	// Get the spell
	spell := skills.GetSpell(spellID)
	if spell == nil {
		ctx.Output.WriteLine("That spell doesn't exist. (Internal error)")
		return
	}

	// Check if already learned
	if progress, ok := p.Skills[spellID]; ok && progress.Learned {
		ctx.Output.WriteLine(fmt.Sprintf("You already know %s!", spell.Name))
		return
	}

	// Learn or improve the spell
	var proficiency int
	if _, hasSkill := p.Skills[spellID]; hasSkill {
		// Already has some proficiency, use teaching amount (usually less aggressive than +10)
		proficiency = p.Skills[spellID].Proficiency + teachAmount
		if proficiency > 100 {
			proficiency = 100
		}
	} else {
		// New spell, start at the teaching amount
		proficiency = teachAmount
	}

	p.Skills[spellID] = &skills.PlayerSkillProgress{
		SpellID:       spellID,
		Proficiency:   proficiency,
		Learned:       true,
		LifetimeCasts: 0,
		LastCastTime:  0,
	}

	ctx.Output.WriteLine(fmt.Sprintf("&YYou carefully study %s and learn &W%s&Y!&w", obj.Short, spell.Name))
	ctx.Output.WriteLine(fmt.Sprintf("Proficiency: %d%% | Mana Cost: %d | Cooldown: %ds",
		proficiency, spell.ManaCost, spell.CooldownSeconds))

	// If item is consumable, destroy it
	if obj.Consumable {
		// Remove from inventory
		newInventory := make([]*game.Object, 0)
		for _, invObj := range p.Inventory {
			if invObj != obj {
				newInventory = append(newInventory, invObj)
			}
		}
		p.Inventory = newInventory
		// Strip article (a/an/the) from short description for better grammar
		itemDesc := strings.TrimPrefix(obj.Short, "a ")
		itemDesc = strings.TrimPrefix(itemDesc, "an ")
		itemDesc = strings.TrimPrefix(itemDesc, "the ")
		ctx.Output.WriteLine(fmt.Sprintf("&YThe %s crumbles to dust after you finish studying it.&w", itemDesc))
	}
}

func cmdTrain(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in to train")
		return
	}

	args = strings.TrimSpace(args)
	if args == "" {
		ctx.Output.WriteLine("Train with whom? (syntax: train <trainer> [maneuver])")
		return
	}

	parts := strings.Fields(args)
	trainerKeyword := parts[0]

	p := ctx.Player
	mob, found := ctx.World.FindMobInRoom(p, trainerKeyword)
	if !found {
		ctx.Output.WriteLine("That trainer is not here.")
		return
	}

	if !mob.IsTrainer {
		ctx.Output.WriteLine(fmt.Sprintf("%s doesn't offer training.", mob.Short))
		return
	}

	spellID := mob.TeachesSpellID
	spell := skills.GetSpell(spellID)
	if spell == nil {
		ctx.Output.WriteLine("That trainer seems confused about what to teach. (Internal error)")
		return
	}

	if len(parts) > 1 {
		desired := strings.ToLower(strings.Join(parts[1:], " "))
		if !strings.Contains(strings.ToLower(spell.Name), desired) {
			ctx.Output.WriteLine(fmt.Sprintf("%s teaches %s. Try: train %s", mob.Short, spell.Name, trainerKeyword))
			return
		}
	}

	if p.Skills == nil {
		p.Skills = make(map[int]*skills.PlayerSkillProgress)
	}

	if progress, ok := p.Skills[spellID]; ok && progress.Learned {
		ctx.Output.WriteLine(fmt.Sprintf("You already know %s.", spell.Name))
		return
	}

	if mob.RequiredStatName != "" && mob.RequiredStatValue > 0 {
		statValue, ok := getPlayerStatValue(p, mob.RequiredStatName)
		if !ok {
			ctx.Output.WriteLine("That trainer seems confused about what to teach. (Invalid stat requirement)")
			return
		}
		if statValue < mob.RequiredStatValue {
			deficit := mob.RequiredStatValue - statValue
			ctx.Output.WriteLine(fmt.Sprintf("%s says: You're not ready yet. You need +%d %s to learn from me.",
				mob.Short, deficit, mob.RequiredStatName))
			return
		}
	}

	proficiency := 10
	p.Skills[spellID] = &skills.PlayerSkillProgress{
		SpellID:       spellID,
		Proficiency:   proficiency,
		Learned:       true,
		LifetimeCasts: 0,
		LastCastTime:  0,
	}

	ctx.Output.WriteLine(fmt.Sprintf("&YYou train with %s and learn &W%s&Y!&w", mob.Short, spell.Name))
	ctx.Output.WriteLine(fmt.Sprintf("Proficiency: %d%% | Cooldown: %ds",
		proficiency, spell.CooldownSeconds))
	if mob.TrainerMessage != "" {
		ctx.Output.WriteLine(mob.TrainerMessage)
	}

	ctx.World.BroadcastSystemToRoomExcept(p, fmt.Sprintf("%s trains with %s.",
		game.CapitalizeName(p.Name), mob.Short))
}

func getPlayerStatValue(p *game.Player, statName string) (int, bool) {
	switch strings.ToLower(statName) {
	case "strength", "str":
		return p.Strength, true
	case "dexterity", "dex":
		return p.Dexterity, true
	case "constitution", "con":
		return p.Constitution, true
	case "intelligence", "int":
		return p.Intelligence, true
	case "wisdom", "wis":
		return p.Wisdom, true
	case "charisma", "cha":
		return p.Charisma, true
	case "luck", "lck":
		return p.Luck, true
	default:
		return 0, false
	}
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
	ctx.World.BroadcastCombatMessage(ctx.Player, fmt.Sprintf("%s summons %s into existence!", game.CapitalizeName(ctx.Player.Name), mob.Short))
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

func cmdRestore(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in.")
		return
	}

	if !ctx.Player.IsKeeper {
		ctx.Output.WriteLine("You do not have the authority to do that.")
		return
	}

	ctx.Player.HP = ctx.Player.MaxHP
	ctx.Player.Mana = ctx.Player.MaxMana

	ctx.Output.WriteLine(fmt.Sprintf("&G✓ Restored to full HP and Mana!&w"))
	ctx.Output.WriteLine(fmt.Sprintf("HP: %d/%d | Mana: %d/%d", ctx.Player.HP, ctx.Player.MaxHP, ctx.Player.Mana, ctx.Player.MaxMana))
}

func cmdHelp(ctx Context, args string) {
	if ctx.Player == nil {
		ctx.Output.WriteLine("You must be logged in.")
		return
	}

	if args == "" {
		ctx.Output.WriteLine("Usage: help <topic_or_spell_name>")
		ctx.Output.WriteLine("")
		ctx.Output.WriteLine("View information about a topic (like 'rules') or spell/maneuver.")
		ctx.Output.WriteLine("Example: help rules")
		ctx.Output.WriteLine("Example: help arcane bolt")
		ctx.Output.WriteLine("")
		ctx.Output.WriteLine("Type 'abilities' to see all your learned abilities.")
		return
	}

	// Try to load help.json
	helpTopics := loadHelpTopics()
	topicKey := strings.ToLower(args)

	// Check if this is a help topic
	if topic, ok := helpTopics[topicKey]; ok {
		ctx.Output.WriteLine(fmt.Sprintf("&Y=== %s ===&w", topic.Title))
		ctx.Output.WriteLine("")
		ctx.Output.WriteLine(topic.Content)
		return
	}

	// Fall back to spell/ability help
	spell, matchCount := skills.FindSpellByPartial(args)

	if matchCount == 0 {
		ctx.Output.WriteLine(fmt.Sprintf("Help topic '%s' not found. Type 'help' for usage.", args))
		return
	}

	if matchCount > 1 {
		ctx.Output.WriteLine(fmt.Sprintf("Ambiguous: '%s' matches %d abilities. Please be more specific.", args, matchCount))
		return
	}

	if spell == nil {
		ctx.Output.WriteLine("Ability not found.")
		return
	}

	// Display spell card
	ctx.Output.WriteLine(fmt.Sprintf("&Y=== %s ===&w", spell.Name))
	ctx.Output.WriteLine("")

	// Description
	if spell.Description != "" {
		ctx.Output.WriteLine(spell.Description)
		ctx.Output.WriteLine("")
	}

	// Show player's proficiency if they know it
	p := ctx.Player
	if p.Skills != nil {
		if progress, ok := p.Skills[spell.ID]; ok && progress.Learned {
			ctx.Output.WriteLine(fmt.Sprintf("&C[YOUR PROFICIENCY]&w"))
			ctx.Output.WriteLine(fmt.Sprintf("  Proficiency: %d%% | Casts: %d", progress.Proficiency, progress.LifetimeCasts))
			ctx.Output.WriteLine("")
		}
	}

	// Cost and cooldown
	ctx.Output.WriteLine(fmt.Sprintf("&CMana Cost:&w %d | &CCooldown:&w %ds", spell.ManaCost, spell.CooldownSeconds))

	// Targeting info
	ctx.Output.WriteLine(fmt.Sprintf("&CTargeting:&w %s (Range: %d)", spell.Targeting.Mode, spell.Targeting.Range))

	// Damage/Effect formula
	if spell.Effects.Damage != "0" && spell.Effects.Damage != "" {
		ctx.Output.WriteLine(fmt.Sprintf("&CEffect:&w %s damage (%s)", spell.Effects.DamageType, spell.Effects.Damage))
		if spell.Effects.SaveType != "none" && spell.Effects.SaveType != "" {
			ctx.Output.WriteLine(fmt.Sprintf("&CSave:&w %s (DC %d)", spell.Effects.SaveType, spell.Effects.SaveDC))
		}
	} else if spell.Effects.Healing != "" {
		ctx.Output.WriteLine(fmt.Sprintf("&CHealing:&w +%s HP", spell.Effects.Healing))
	}

	ctx.Output.WriteLine("")
}

type HelpTopic struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

func loadHelpTopics() map[string]HelpTopic {
	helpMap := make(map[string]HelpTopic)

	data, err := os.ReadFile("system/help.json")
	if err != nil {
		// If file doesn't exist, return empty map; help will still work for spells
		return helpMap
	}

	var helpData map[string]HelpTopic
	err = json.Unmarshal(data, &helpData)
	if err != nil {
		// If JSON is invalid, return empty map
		return helpMap
	}

	return helpData
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
