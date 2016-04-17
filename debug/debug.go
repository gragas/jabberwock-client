package debug

import (
	"fmt"
	"github.com/gragas/go-sdl2/sdl"	
	"github.com/gragas/jabberwock-client/utils"
	"github.com/gragas/jabberwock-lib/entity"
	"github.com/gragas/jabberwock-lib/player"
	"github.com/gragas/jabberwock-lib/protocol"
	"github.com/gragas/jabberwock-server/game"
	"net"
)

var conn net.Conn

func Init(ip string, port int, quiet bool, debug bool, serverDebug bool) {
	done := make(chan string)
	go game.StartGame(ip, port, quiet, serverDebug, done)
	<-done
	var registered bool
	var attempts int
	player := player.Player{Entity: entity.Entity{ID: uint64(protocol.GenerateEntityID),
		Name: player.DefaultName,
		Health: player.DefaultHealth,
		MaxHealth: player.DefaultMaxHealth,
		Energy: player.DefaultEnergy,
		MaxEnergy: player.DefaultMaxEnergy,
		Spirit: player.DefaultSpirit,
		MaxSpirit: player.DefaultMaxSpirit,
		Summoning: player.DefaultSummoning,
		Alteration: player.DefaultAlteration,
		Willpower: player.DefaultWillpower,
		Divinity: player.DefaultDivinity,
		Lifebringing: player.DefaultLifebringing,
		X: player.DefaultX, Y: player.DefaultY,
		Xv: player.DefaultXv, Yv: player.DefaultYv,
		W: player.DefaultW, H: player.DefaultH}}
	for ; attempts < 5; attempts++ {
		conn, registered = utils.RegisterClient(ip, port, player, debug)
		if registered {
			break
		}
	}
	if !registered {
		fmt.Printf("CLIENT: Failed to register after %v attempts.\n", attempts)
	}
	utils.Loop = utils.LoopFuncs{pollEvents, update, draw}
}

func pollEvents() {
	for {
		if event := sdl.PollEvent(); event != nil {
			switch event.(type) {
			case *sdl.QuitEvent:
				utils.Running = false
				err := conn.Close()
				if err != nil {
					fmt.Printf("CLIENT: Failed to close connection.\n")
				} else {
					fmt.Printf("CLIENT: Successfully closed connection.\n")
				}
			}
		} else {
			break
		}
	}
}

func update() {
	
}

func draw() {
	
}
