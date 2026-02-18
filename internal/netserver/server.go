package netserver

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"njata/internal/commands"
	"njata/internal/game"
	"njata/internal/parser"
	"njata/internal/persist"
	"njata/internal/skills"
)

const playerDataDir = "players"

type Server struct {
	world    *game.World
	registry *commands.Registry
	port     int
	logger   func(string)
}

func NewServer(world *game.World, registry *commands.Registry, port int, logger func(string)) *Server {
	return &Server{
		world:    world,
		registry: registry,
		port:     port,
		logger:   logger,
	}
}

func (s *Server) Run(ctx context.Context) error {
	address := fmt.Sprintf(":%d", s.port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer listener.Close()

	if s.logger != nil {
		s.logger(fmt.Sprintf("Listening on %s", address))
	}

	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()

	// Start autosave ticker - saves all players every 5 minutes
	go s.startAutosaveTimer(ctx)

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				continue
			}
			return err
		}

		if s.logger != nil {
			s.logger(fmt.Sprintf("Connection from %s", conn.RemoteAddr()))
		}

		go s.handleConn(conn)
	}
}

func (s *Server) startAutosaveTimer(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// Save all players
			players := s.world.PlayersSnapshot()
			for _, player := range players {
				if player != nil {
					record := persist.PlayerToRecord(player)
					if err := persist.SavePlayer(playerDataDir, record); err != nil {
						if s.logger != nil {
							s.logger(fmt.Sprintf("autosave error for %s: %v", player.Name, err))
						}
					}
				}
			}
		}
	}
}

func (s *Server) handleConn(conn net.Conn) {
	session := NewSession(conn)
	defer session.Close()

	WriteBanner(session)
	session.WriteLine("")
	session.WriteLine("")

	var player *game.Player
	var isNewPlayer bool

	for {
		session.Write("Name: ")
		line, err := session.ReadLine()
		if err != nil {
			return
		}

		name := strings.TrimSpace(line)
		if err := game.ValidateName(name); err != nil {
			session.WriteLine("Invalid name. Use 3-16 letters or digits.")
			continue
		}

		// Try to load existing player
		record, exists, err := persist.LoadPlayer(playerDataDir, name)
		if err != nil && exists {
			session.WriteLine("Error loading character. Please try again.")
			continue
		}

		isNewPlayer = !exists

		// Create new player struct
		player = &game.Player{
			Name:       name,
			Output:     session,
			Disconnect: session.RequestDisconnect,
			AutoExits:  true,
			Skills:     make(map[int]*skills.PlayerSkillProgress),
			Inventory:  []*game.Object{},
			Equipment:  make(map[string]*game.Object),
		}

		// If existing player, load their stats
		if !isNewPlayer {
			persist.RecordToPlayer(player, record)
		} else {
			// New player: run character creation
			creation := NewCharacterCreation(session, player)
			if err := creation.Run(); err != nil {
				return
			}

			// Auto-teach new players Arcane Bolt (spell ID 1001) as a starter spell
			arcaneBoltSpell := skills.GetSpell(1001)
			if arcaneBoltSpell != nil {
				player.Skills[1001] = &skills.PlayerSkillProgress{
					SpellID:       1001,
					Proficiency:   50,
					Learned:       true,
					LifetimeCasts: 0,
					LastCastTime:  0,
				}
				session.WriteLine(fmt.Sprintf("You have learned &Y%s&w as your first spell!", arcaneBoltSpell.Name))
			}
		}

		if player.Location != 0 && !s.world.HasRoom(player.Location) {
			player.Location = 0
		}

		if err := s.world.AddPlayer(player); err != nil {
			session.WriteLine("That name is already in use.")
			player = nil
			continue
		}

		defer func() {
			if player != nil {
				record := persist.PlayerToRecord(player)
				if err := persist.SavePlayer(playerDataDir, record); err != nil && s.logger != nil {
					s.logger(fmt.Sprintf("save error for %s: %v", player.Name, err))
				}
				s.world.RemovePlayer(player.Name)
				s.world.BroadcastSystemToRoomExcept(player, fmt.Sprintf("%s has left the game.", game.CapitalizeName(player.Name)))
			}
		}()

		s.world.BroadcastSystemToRoomExcept(player, fmt.Sprintf("%s has entered the game.", game.CapitalizeName(player.Name)))
		session.WriteLine(fmt.Sprintf("Welcome back, %s!", game.CapitalizeName(player.Name)))

		view, err := s.world.DescribeRoom(player)
		if err == nil {
			commands.DisplayRoomView(session, view, player.AutoExits)
		}
		break
	}

	for {
		if session.IsDisconnectRequested() {
			return
		}

		session.Write("> ")
		line, err := session.ReadLine()
		if err != nil {
			return
		}

		command, args := parser.ParseInput(line)
		if command == "" {
			continue
		}

		ctx := commands.Context{
			World:      s.world,
			Player:     player,
			Output:     session,
			Disconnect: session.RequestDisconnect,
		}

		if !s.registry.Execute(ctx, command, args) {
			session.WriteLine("Huh? Type 'help' for commands.")
		}
	}
}
