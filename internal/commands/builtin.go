package commands

import (
    "fmt"
    "strings"
)

func RegisterBuiltins(registry *Registry) {
    registry.Register("look", cmdLook)
    registry.Register("say", cmdSay)
    registry.Register("who", cmdWho)
    registry.Register("stats", cmdStats)
    registry.Register("score", cmdScore)
    registry.Register("exits", cmdExits)
    registry.Register("autoexits", cmdAutoexits)
    registerMovement(registry)
    registry.Register("help", func(ctx Context, args string) {
        commands := registry.List()
        ctx.Output.WriteLine("Commands: " + strings.Join(commands, ", "))
    })
    registry.Register("quit", cmdQuit)
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

    ctx.Output.WriteLine(view.Name)
    if view.AreaName != "" || view.AreaAuthor != "" {
        areaLine := "[area: " + view.AreaName
        if view.AreaAuthor != "" {
            areaLine += " by " + view.AreaAuthor
        }
        areaLine += "]"
        ctx.Output.WriteLine(areaLine)
    }
    if view.Description != "" {
        ctx.Output.WriteLine(view.Description)
    }

    if ctx.Player.AutoExits {
        ctx.Output.WriteLine(FormatExits(view.Exits))
    }

    if len(view.Others) == 0 {
        ctx.Output.WriteLine("You are alone here.")
        return
    }

    ctx.Output.WriteLine("Also here: " + strings.Join(view.Others, ", "))
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
    ctx.Output.WriteLine(fmt.Sprintf("Class: %d | Race: %d | Sex: %d", p.Class, p.Race, p.Sex))
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
    ctx.Output.WriteLine(fmt.Sprintf("Class: %d | Race: %d | Sex: %d", p.Class, p.Race, p.Sex))
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

            ctx.Output.WriteLine(view.Name)
            if view.Description != "" {
                ctx.Output.WriteLine(view.Description)
            }

            if ctx.Player.AutoExits {
                ctx.Output.WriteLine(FormatExits(view.Exits))
            }

            if len(view.Others) > 0 {
                ctx.Output.WriteLine("Also here: " + strings.Join(view.Others, ", "))
            }
        })
    }
}
