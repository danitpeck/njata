package commands

import (
    "fmt"
    "strings"
    "time"

    "njata/internal/classes"
    "njata/internal/game"
    "njata/internal/races"
    "njata/internal/skills"
)

func RegisterBuiltins(registry *Registry) {
    registry.Register("look", cmdLook)
    registry.Register("say", cmdSay)
    registry.Register("who", cmdWho)
    registry.Register("stats", cmdStats)
    registry.Register("score", cmdScore)
    registry.Register("exits", cmdExits)
    registry.Register("autoexits", cmdAutoexits)
    registry.Register("astat", cmdAstat)
    registry.Register("spellbook", cmdSpellbook)
    registry.Register("cast", cmdCast)
    registry.Register("study", cmdStudy)
    registry.Register("makekeeper", cmdMakeKeeper)
    registry.Register("removekeeper", cmdRemoveKeeper)
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
    
    className := "Unknown"
    if c := classes.GetByID(p.Class); c != nil {
        className = c.Name
    }
    
    raceName := "Unknown"
    if r := races.GetByID(p.Race); r != nil {
        raceName = r.Name
    }
    
    sexNames := []string{"neuter", "male", "female"}
    sexName := "unknown"
    if p.Sex >= 0 && p.Sex < len(sexNames) {
        sexName = sexNames[p.Sex]
    }
    
    ctx.Output.WriteLine(fmt.Sprintf("Race: %s | Class: %s | Sex: %s", raceName, className, sexName))
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

func cmdScore(ctx Context, args string) {
    p := ctx.Player
    
    ctx.Output.WriteLine("")
    ctx.Output.WriteLine(fmt.Sprintf("==== %s (Level %d) ====", p.Name, p.Level))
    
    className := "Unknown"
    if c := classes.GetByID(p.Class); c != nil {
        className = c.Name
    }
    
    raceName := "Unknown"
    if r := races.GetByID(p.Race); r != nil {
        raceName = r.Name
    }
    
    sexNames := []string{"neuter", "male", "female"}
    sexName := "unknown"
    if p.Sex >= 0 && p.Sex < len(sexNames) {
        sexName = sexNames[p.Sex]
    }
    
    ctx.Output.WriteLine(fmt.Sprintf("Race: %s | Class: %s | Sex: %s", raceName, className, sexName))
    ctx.Output.WriteLine("")
    
    // Attributes
    attrNames := []string{"STR", "DEX", "CON", "INT", "WIS", "CHA", "LCK"}
    ctx.Output.WriteLine("ATTRIBUTES:")
    for i, name := range attrNames {
        ctx.Output.WriteLine(fmt.Sprintf("  %s  : %2d", name, p.Attributes[i]))
    }
    ctx.Output.WriteLine("")
    
    // Vitals
    ctx.Output.WriteLine("VITALS:")
    ctx.Output.WriteLine(fmt.Sprintf("  Hitpoints: %d / %d", p.HP, p.MaxHP))
    ctx.Output.WriteLine(fmt.Sprintf("  Mana:      %d / %d", p.Mana, p.MaxMana))
    ctx.Output.WriteLine(fmt.Sprintf("  Movement:  %d / %d", p.Move, p.MaxMove))
    ctx.Output.WriteLine("")
    
    // Experience & Gold
    ctx.Output.WriteLine("EXPERIENCE & WEALTH:")
    ctx.Output.WriteLine(fmt.Sprintf("  Experience: %d", p.Experience))
    ctx.Output.WriteLine(fmt.Sprintf("  Gold:       %d", p.Gold))
    ctx.Output.WriteLine("")
    
    // Combat
    ctx.Output.WriteLine("COMBAT STATS:")
    ctx.Output.WriteLine(fmt.Sprintf("  Alignment:  %d", p.Alignment))
    ctx.Output.WriteLine(fmt.Sprintf("  Hitroll:    %d", p.Hitroll))
    ctx.Output.WriteLine(fmt.Sprintf("  Damroll:    %d", p.Damroll))
    ctx.Output.WriteLine(fmt.Sprintf("  Armor:      %d", p.Armor))
    ctx.Output.WriteLine("")
}


func cmdSpellbook(ctx Context, args string) {
    if ctx.Player == nil {
        ctx.Output.WriteLine("You must be logged in to use spellbook")
        return
    }

    p := ctx.Player

    // Initialize skills map if needed
    if p.Skills == nil {
        p.Skills = make(map[int]*skills.PlayerSkillProgress)
    }

    ctx.Output.WriteLine("=== SPELLBOOK ===")
    ctx.Output.WriteLine("")

    if len(p.Skills) == 0 {
        ctx.Output.WriteLine("You haven't learned any spells yet.")
        ctx.Output.WriteLine("")
        return
    }

    // Display all learned spells
    allSpells := skills.AllSpells()
    for spellID, progress := range p.Skills {
        if !progress.Learned {
            continue
        }

        spell := allSpells[spellID]
        if spell == nil {
            continue
        }

        ctx.Output.WriteLine(fmt.Sprintf("[%d] %s", spell.ID, spell.Name))
        ctx.Output.WriteLine(fmt.Sprintf("    Mana: %d | Cooldown: %ds | Proficiency: %d%%", 
            spell.ManaCost, spell.CooldownSeconds, progress.Proficiency))
        ctx.Output.WriteLine(fmt.Sprintf("    Casts: %d", progress.LifetimeCasts))
    }
    ctx.Output.WriteLine("")
}


func cmdCast(ctx Context, args string) {
    if ctx.Player == nil {
        ctx.Output.WriteLine("You must be logged in to cast spells")
        return
    }

    args = strings.TrimSpace(args)
    if args == "" {
        ctx.Output.WriteLine("Cast what? (syntax: cast <spell name>)")
        return
    }

    p := ctx.Player

    // Initialize skills map if needed
    if p.Skills == nil {
        p.Skills = make(map[int]*skills.PlayerSkillProgress)
    }

    // Find the spell by name
    spell := skills.GetSpellByName(args)
    if spell == nil {
        ctx.Output.WriteLine(fmt.Sprintf("You don't know any spell called '%s'", args))
        return
    }

    // Check if player has learned this spell
    skillProgress, hasSkill := p.Skills[spell.ID]
    if !hasSkill || !skillProgress.Learned {
        ctx.Output.WriteLine(fmt.Sprintf("You haven't learned %s yet.", spell.Name))
        return
    }

    // Check if player can cast (mana, cooldown)
    now := time.Now().UnixNano()
    canCast, reason := skillProgress.CanCast(spell, p.Mana, now)
    if !canCast {
        ctx.Output.WriteLine(reason)
        return
    }

    // Cast the spell!
    p.Mana -= spell.ManaCost
    skillProgress.UpdateCooldown()
    skillProgress.UpdateProficiency(1) // +1% proficiency per cast

    // Success message
    msg := strings.ReplaceAll(spell.Messages.Cast, "$actor", p.Name)
    msg = strings.ReplaceAll(msg, "$spell", spell.Name)
    ctx.Output.WriteLine(fmt.Sprintf("%s (Proficiency: %d%%)", msg, skillProgress.Proficiency))
    ctx.Output.WriteLine(fmt.Sprintf("Mana remaining: %d/%d", p.Mana, p.MaxMana))

    // For now, spell effects (damage, healing, etc.) would be implemented here
    // This is a placeholder for actual combat/effect system integration
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
        "wand":          1001, // Magic Missile
        "missile":       1001,
        "fireball":      1002,
        "heal":          1003,
        "healing":       1003,
        "scroll":        1003, // default to heal for scrolls
        "blindness":     1004,
        "invisibility":  1005,
        "teleport":      1006,
        "frost":         1007,
        "cold":          1007,
        "identify":      1008,
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
