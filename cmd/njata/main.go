package main

import (
    "context"
    "flag"
    "fmt"
    "os"
    "os/signal"

    "njata/internal/area"
    "njata/internal/commands"
    "njata/internal/game"
    "njata/internal/netserver"
)

func main() {
    port := flag.Int("port", 4000, "listen port")
    flag.Parse()

    rooms, start, err := area.LoadRoomsFromDir("areas")
    if err != nil {
        fmt.Printf("Area load error: %v\n", err)
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
