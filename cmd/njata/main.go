package main

import (
    "context"
    "flag"
    "fmt"
    "os"
    "os/signal"

    "njata/internal/area"
    "njata/internal/classes"
    "njata/internal/commands"
    "njata/internal/game"
    "njata/internal/netserver"
    "njata/internal/races"
)

func main() {
    port := flag.Int("port", 4000, "listen port")
    flag.Parse()

    rooms, start, err := area.LoadRoomsFromDir("areas")
    if err != nil {
        fmt.Printf("Area load error: %v\n", err)
    }

    if err := races.Load("races"); err != nil {
        fmt.Printf("Race load error: %v\n", err)
        os.Exit(1)
    }

    if err := classes.Load("classes"); err != nil {
        fmt.Printf("Class load error: %v\n", err)
        os.Exit(1)
    }

    world := game.CreateWorldFromRooms(rooms, start)
    registry := commands.NewRegistry()
    commands.RegisterBuiltins(registry)

    logger := func(message string) {
        fmt.Println(message)
    }

    server := netserver.NewServer(world, registry, *port, logger)

    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    if err := server.Run(ctx); err != nil && ctx.Err() == nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
