package commands

import "njata/internal/game"

type Context struct {
    World      *game.World
    Player     *game.Player
    Output     game.Output
    Disconnect func(reason string)
}
