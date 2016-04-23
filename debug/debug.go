package debug

import (
	"bufio"
	"encoding/json"
	"errors"
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
var reader *bufio.Reader
var receiverToHandler chan string
var clientPlayer *player.Player
var clientPlayerView *player.PlayerView
var entities map[uint64]entity.Entity
var players map[uint64]*player.Player
var playerViews map[uint64]*player.PlayerView

func Init(ip string, port int, initialized chan bool, quiet bool, debug bool, serverDebug bool) {
	/* initialize "global" variables */
	entities = make(map[uint64]entity.Entity)
	players = make(map[uint64]*player.Player)
	playerViews = make(map[uint64]*player.PlayerView)
	/*********************************/

	/* start up a server in the background */
	done := make(chan bool)
	go game.StartGame(ip, port, quiet, serverDebug, done)
	<-done
	/***************************************/

	/* register with the server */
	sendPlayer := player.NewDefaultPlayer()
	var registered bool
	conn, reader, clientPlayer, registered = utils.RegisterClient(ip, port, sendPlayer, debug)
	if registered {
		players[clientPlayer.GetID()] = clientPlayer
		entities[clientPlayer.GetID()] = clientPlayer
	} else {
		panic(errors.New("ERROR: Client failed to register with server.\n"))
	}
	/****************************/

	/* setup a player view */
	surf, rect := entity.NewDefaultEntityView(clientPlayer)
	clientPlayerView = &player.PlayerView{PlayerPtr: clientPlayer, Surface: surf, Rect: rect}
	playerViews[clientPlayer.GetID()] = clientPlayerView
	/***********************/

	utils.Loop = utils.LoopFuncs{pollEvents, draw}
	initialized <- true
	receiver(debug)
}

func quit() {
	utils.Running = false
	err := conn.Close()
	if err != nil {
		fmt.Printf("CLIENT: Failed to close connection.\n")
	} else {
		fmt.Printf("CLIENT: Successfully closed connection with %v.\n", conn.RemoteAddr())
	}
}

func pollEvents() {
	for {
		if event := sdl.PollEvent(); event != nil {
			switch event.(type) {
			case *sdl.QuitEvent:
				quit()
			case *sdl.KeyDownEvent:
				sym := event.(*sdl.KeyDownEvent).Keysym.Sym
				switch sym {
				case sdl.K_ESCAPE:
					quit()
				case sdl.K_w:
					if clientPlayer != nil {
						entity.MoveNet(clientPlayer, conn, protocol.EntityStartMoveUp)
					}
				case sdl.K_s:
					if clientPlayer != nil {
						entity.MoveNet(clientPlayer, conn, protocol.EntityStartMoveDown)
					}
				case sdl.K_d:
					if clientPlayer != nil {
						entity.MoveNet(clientPlayer, conn, protocol.EntityStartMoveRight)
					}
				case sdl.K_a:
					if clientPlayer != nil {
						entity.MoveNet(clientPlayer, conn, protocol.EntityStartMoveLeft)
					}
				default:
					fmt.Println("CLIENT: Key down:", sym)
				}
			case *sdl.KeyUpEvent:
				sym := event.(*sdl.KeyUpEvent).Keysym.Sym
				switch sym {
				case sdl.K_w:
					if clientPlayer != nil {
						entity.MoveNet(clientPlayer, conn, protocol.EntityStopMoveUp)
					}
				case sdl.K_s:
					if clientPlayer != nil {
						entity.MoveNet(clientPlayer, conn, protocol.EntityStopMoveDown)
					}
				case sdl.K_d:
					if clientPlayer != nil {
						entity.MoveNet(clientPlayer, conn, protocol.EntityStopMoveRight)
					}
				case sdl.K_a:
					if clientPlayer != nil {
						entity.MoveNet(clientPlayer, conn, protocol.EntityStopMoveLeft)
					}
				default:
					fmt.Println("CLIENT: Key up:", sym)
				}
			}
		} else {
			break
		}
	}
}

func receiver(debug bool) {
	for {
		msg, err := reader.ReadString(byte(protocol.EndOfMessage))
		if err != nil {
			fmt.Printf("CLIENT: Disconnected from server %v.\n", conn.RemoteAddr())
			return
		}
		if len(msg) < 2 {
			if debug {
				fmt.Printf("CLIENT: Received malformed message.\n")
			}
			continue
		}
		update(msg[:len(msg)-1], debug)
	}
}

func update(msg string, debug bool) {
	switch protocol.Code(msg[0]) {
	case protocol.UpdatePlayers:
		err := json.Unmarshal([]byte(msg[1:]), &players)
		if err != nil {
			fmt.Printf("CLIENT: Received invalid msg with UpdatePlayers code: %v\n", msg)
			return
		}
		if players[clientPlayer.GetID()] != nil {
			clientPlayer = players[clientPlayer.GetID()]
		}
		surf, rect := entity.NewDefaultEntityView(clientPlayer)
		clientPlayerView.Surface = surf
		clientPlayerView.Rect = rect
	default:
		fmt.Printf("CLIENT: Invalid protocol.Code.\ncode: %v\nmsg: %v\n", protocol.Code(msg[0]), msg)
	}
}

func draw(dest *sdl.Surface) {
	utils.Surface.FillRect(nil, sdl.MapRGBA(utils.Surface.Format, 0xFF, 0xEE, 0xCC, 0xFF))
	for _, pv := range playerViews {
		pv.Draw(dest)
	}
}
