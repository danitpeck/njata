package commands

import (
    "fmt"
    "strings"
)

func RegisterBuiltins(registry *Registry) {
    registry.Register("look", cmdLook)
    registry.Register("say", cmdSay)
    registry.Register("who", cmdWho)
    registry.Register("help", func(ctx Context, args string) {
        commands := registry.List()
        ctx.Output.WriteLine("Commands: " + strings.Join(commands, ", "))
    })
    registry.Register("quit", cmdQuit)
}

func cmdLook(ctx Context, args string) {
    room := ctx.World.RoomSnapshot()
    ctx.Output.WriteLine(room.Name)
    ctx.Output.WriteLine(room.Description)

    others := ctx.World.ListPlayersExcept(ctx.Player.Name)
    if len(others) == 0 {
        ctx.Output.WriteLine("You are alone here.")
        return
    }

    ctx.Output.WriteLine("Also here: " + strings.Join(others, ", "))
}

func cmdSay(ctx Context, args string) {
    if strings.TrimSpace(args) == "" {
        ctx.Output.WriteLine("Say what?")
        return
    }

    ctx.World.BroadcastSay(ctx.Player.Name, args)
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
