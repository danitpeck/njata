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
    classID := p.Class

    // Get skills available to this class
    availableSkills := skills.GetSkillsByClass(classID)
    if len(availableSkills) == 0 {
        ctx.Output.WriteLine("You have no spells available.")
        return
    }

    ctx.Output.WriteLine("=== SPELLBOOK ===")
    ctx.Output.WriteLine("")

    // Initialize skills map if needed
    if p.Skills == nil {
        p.Skills = make(map[int]int)
    }
    if p.SkillCooldowns == nil {
        p.SkillCooldowns = make(map[int]int64)
    }

    // Learn all class skills if not already learned
    for _, skill := range availableSkills {
        if _, learned := p.Skills[skill.ID]; !learned {
            p.Skills[skill.ID] = 100 // Start at full proficiency
        }
    }

    for _, skill := range availableSkills {
        proficiency := p.Skills[skill.ID]
        ctx.Output.WriteLine(fmt.Sprintf("[%d] %s - %s", skill.ID, skill.Name, skill.Description))
        ctx.Output.WriteLine(fmt.Sprintf("    Type: %s | Mana: %d | Cooldown: %d | Proficiency: %d%%", 
            skill.Type, skill.ManaCost, skill.CooldownSeconds, proficiency))
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
        ctx.Output.WriteLine("Cast what? (syntax: cast <spell name or number>)")
        return
    }

    p := ctx.Player

    // Initialize maps if needed
    if p.Skills == nil {
        p.Skills = make(map[int]int)
    }
    if p.SkillCooldowns == nil {
        p.SkillCooldowns = make(map[int]int64)
    }

    // Find the skill by name or ID
    var targetSkill *skills.Skill
    allSkills := skills.ListAll()
    
    // Try to match by name
    for _, s := range allSkills {
        if strings.EqualFold(s.Name, args) {
            targetSkill = s
            break
        }
    }

    if targetSkill == nil {
        ctx.Output.WriteLine(fmt.Sprintf("You don't know any spell called '%s'", args))
        return
    }

    // Check if player learned this skill
    proficiency, hasSkill := p.Skills[targetSkill.ID]
    if !hasSkill {
        ctx.Output.WriteLine(fmt.Sprintf("You haven't learned %s yet.", targetSkill.Name))
        return
    }

    // Check mana
    if p.Mana < targetSkill.ManaCost {
        ctx.Output.WriteLine(fmt.Sprintf("You need %d mana to cast %s, but only have %d.", targetSkill.ManaCost, targetSkill.Name, p.Mana))
        return
    }

    // Check cooldown
    now := time.Now().UnixNano()
    lastCast := p.SkillCooldowns[targetSkill.ID]
    timeSinceCast := (now - lastCast) / 1e9 // Convert to seconds
    cooldown := int64(targetSkill.CooldownSeconds)

    if timeSinceCast < cooldown {
        ctx.Output.WriteLine(fmt.Sprintf("%s is still on cooldown for %d more seconds.", targetSkill.Name, cooldown-timeSinceCast))
        return
    }

    // Cast the spell!
    p.Mana -= targetSkill.ManaCost
    p.SkillCooldowns[targetSkill.ID] = now

    // Success message
    ctx.Output.WriteLine(fmt.Sprintf("You cast %s!", targetSkill.Name))
    ctx.Output.WriteLine(fmt.Sprintf("Effect: %s (Proficiency: %d%%)", targetSkill.Description, proficiency))
    ctx.Output.WriteLine(fmt.Sprintf("Mana remaining: %d/%d", p.Mana, p.MaxMana))
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
