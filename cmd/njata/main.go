package main

import (
    "context"
    "flag"
    "fmt"
    "os"
    "os/signal"
    "time"

    "njata/internal/area"
    "njata/internal/classes"
    "njata/internal/commands"
    "njata/internal/config"
    "njata/internal/game"
    "njata/internal/netserver"
    "njata/internal/races"
    "njata/internal/skills"
)

func main() {
    port := flag.Int("port", 4000, "listen port")
    configPath := flag.String("config", "config.json", "config file path")
    flag.Parse()

    cfg, err := config.Load(*configPath)
    if err != nil {
        fmt.Printf("Config load error: %v\n", err)
        os.Exit(1)
    }

    rooms, mobiles, objects, start, err := area.LoadRoomsFromDir("areas")
    if err != nil {
        fmt.Printf("Area load error: %v\n", err)
    }

    if cfg.StartRoomVnum != 0 {
        if _, ok := rooms[cfg.StartRoomVnum]; ok {
            start = cfg.StartRoomVnum
        } else {
            fmt.Printf("Config start_room_vnum %d not found; using %d\n", cfg.StartRoomVnum, start)
        }
    }

    if err := races.Load("races"); err != nil {
        fmt.Printf("Race load error: %v\n", err)
        os.Exit(1)
    }

    if err := classes.Load("classes"); err != nil {
        fmt.Printf("Class load error: %v\n", err)
        os.Exit(1)
    }

    if err := skills.Load("skills/skills.json"); err != nil {
        fmt.Printf("Skills load error: %v\n", err)
        os.Exit(1)
    }

    world := game.CreateWorldFromRooms(rooms, start)
    world.SetPrototypes(mobiles, objects)
    registry := commands.NewRegistry()
    commands.RegisterBuiltins(registry)

    logger := func(message string) {
        fmt.Println(message)
    }

    server := netserver.NewServer(world, registry, *port, logger)

    ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
    defer stop()

    if cfg.RespawnDefaultMinutes > 0 {
        interval := time.Duration(cfg.RespawnDefaultMinutes) * time.Minute
        ticker := time.NewTicker(interval)
        go func() {
            defer ticker.Stop()
            for {
                select {
                case <-ctx.Done():
                    return
                case <-ticker.C:
                    world.RespawnTick(cfg.RespawnDefaultMinutes, logger)
                }
            }
        }()
    }

    if err := server.Run(ctx); err != nil && ctx.Err() == nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
