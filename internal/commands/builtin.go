package commands

import (
    "fmt"
    "strings"
)

func RegisterBuiltins(registry *Registry) {
    registry.Register("look", cmdLook)
    registry.Register("say", cmdSay)
    registry.Register("who", cmdWho)
    registerMovement(registry)
    registry.Register("help", func(ctx Context, args string) {
        commands := registry.List()
        ctx.Output.WriteLine("Commands: " + strings.Join(commands, ", "))
    })
    registry.Register("quit", cmdQuit)
}

func cmdLook(ctx Context, args string) {
    view, err := ctx.World.DescribeRoom(ctx.Player)
    if err != nil {
        ctx.Output.WriteLine("You are nowhere.")
        return
    }

    ctx.Output.WriteLine(view.Name)
    if view.Description != "" {
        ctx.Output.WriteLine(view.Description)
    }

    if len(view.Exits) == 0 {
        ctx.Output.WriteLine("Exits: none")
    } else {
        ctx.Output.WriteLine("Exits: " + strings.Join(view.Exits, ", "))
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

func cmdQuit(ctx Context, args string) {
    ctx.Output.WriteLine("Goodbye.")
    if ctx.Disconnect != nil {
        ctx.Disconnect("quit")
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

            ctx.Output.WriteLine(view.Name)
            if view.Description != "" {
                ctx.Output.WriteLine(view.Description)
            }

            if len(view.Exits) == 0 {
                ctx.Output.WriteLine("Exits: none")
            } else {
                ctx.Output.WriteLine("Exits: " + strings.Join(view.Exits, ", "))
            }

            if len(view.Others) > 0 {
                ctx.Output.WriteLine("Also here: " + strings.Join(view.Others, ", "))
            }
        })
    }
}
