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
	"github.com/gragas/jabberwock-lib/textures"
	"net"
	"os"
	"strconv"
	"time"
)

var conn net.Conn
var reader *bufio.Reader
var receiverToHandler chan string
var clientPlayer *player.Player
var clientPlayerView *player.PlayerView
var entities map[uint64]entity.Entity
var players map[uint64]*player.Player
var jsonPlayers map[string]*player.Player // marshalled players
var playerViews map[uint64]*player.PlayerView

var flag [2]bool
var turn int

func Init(ip string, port int, initialized chan bool, quiet bool, debug bool, serverDebug bool) {
	/* initialize the textures cache */
	textures.Init()

	/* initialize "global" variables */
	entities = make(map[uint64]entity.Entity)
	players = make(map[uint64]*player.Player)
	jsonPlayers = make(map[string]*player.Player)
	playerViews = make(map[uint64]*player.PlayerView)
	/*********************************/

	/* start up a server in the background */
	// done := make(chan bool)
	// go game.StartGame(ip, port, quiet, serverDebug, false, done)
	// <-done
	/***************************************/

	/* register with the server */
	sendPlayer := player.NewDefaultPlayer()
	var registered bool
	conn, reader, clientPlayer, registered = utils.RegisterClient(ip, port, sendPlayer, debug)
	if registered {
		players[clientPlayer.GetID()] = clientPlayer
		jsonPlayers[strconv.FormatUint(clientPlayer.GetID(), 10)] = clientPlayer
		entities[clientPlayer.GetID()] = clientPlayer
	} else {
		panic(errors.New("ERROR: Client failed to register with server.\n"))
	}
	/****************************/

	/* setup a player view */
	clientPlayerView = clientPlayer.NewDefaultPlayerView(utils.Renderer)
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
	cleanupRenderer()
	sdl.Quit()
	os.Exit(0)
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
						entity.MoveNet(clientPlayer, conn, true, entity.Up)
					}
				case sdl.K_s:
					if clientPlayer != nil {
						entity.MoveNet(clientPlayer, conn, true, entity.Down)
					}
				case sdl.K_d:
					if clientPlayer != nil {
						entity.MoveNet(clientPlayer, conn, true, entity.Right)
					}
				case sdl.K_a:
					if clientPlayer != nil {
						entity.MoveNet(clientPlayer, conn, true, entity.Left)
					}
				default:
					fmt.Println("CLIENT: Key down:", sym)
				}
			case *sdl.KeyUpEvent:
				sym := event.(*sdl.KeyUpEvent).Keysym.Sym
				switch sym {
				case sdl.K_w:
					if clientPlayer != nil {
						entity.MoveNet(clientPlayer, conn, false, entity.Up)
					}
				case sdl.K_s:
					if clientPlayer != nil {
						entity.MoveNet(clientPlayer, conn, false, entity.Down)
					}
				case sdl.K_d:
					if clientPlayer != nil {
						entity.MoveNet(clientPlayer, conn, false, entity.Right)
					}
				case sdl.K_a:
					if clientPlayer != nil {
						entity.MoveNet(clientPlayer, conn, false, entity.Left)
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
	// Peterson's Solution
	flag[0] = true
	turn = 1
	for flag[1] && turn == 1 {}
	
	switch protocol.Code(msg[0]) {
	case protocol.UpdatePlayers:
		err := json.Unmarshal([]byte(msg[1:]), &jsonPlayers)
		if err != nil {
			fmt.Printf("CLIENT: Received invalid msg with UpdatePlayers code: %v\n", msg)
			return
		}
		// now update players to reflect the new jsonPlayers
		for strkey, v := range jsonPlayers {
			k, err := strconv.ParseUint(strkey, 10, 64)
			if err != nil {
				panic(err)
			}
			players[k] = v
			if playerViews[k] == nil {
				playerViews[k] = players[k].NewDefaultPlayerView(utils.Renderer)
			} else {
				playerViews[k].PlayerPtr = players[k]
			}
			entities[k] = players[k]
		}
		// make sure the client player is pointing in the right direction
		if players[clientPlayer.GetID()] != nil {
			clientPlayer = players[clientPlayer.GetID()]
		}
	case protocol.Disconnect:
		id, err := strconv.ParseUint(msg[1:], 10, 64)
		if err != nil {
			fmt.Printf("CLIENT: Could not parse entity ID of disconnected client.\n")
		} else {
			delete(players, id)
			delete(playerViews, id)
			delete(jsonPlayers ,strconv.FormatUint(id, 10))
			delete(entities, id)
		}
	default:
		fmt.Printf("CLIENT: Invalid protocol.Code.\ncode: %v\nmsg: %v\n", protocol.Code(msg[0]), msg)
	}

	flag[0] = false
}

func draw(dest *sdl.Surface) {
	// Peterson's Solution
	flag[1] = true
	turn = 0
	for flag[0] && turn == 0 {}

	err := utils.Renderer.SetRenderTarget(nil); if err != nil { panic(err) }
	err = utils.Renderer.SetDrawBlendMode(sdl.BLENDMODE_BLEND); if err != nil { panic(err) }
	err = utils.Renderer.SetDrawColor(0xFF, 0xFF, 0xFF, 0xFF)
	utils.Renderer.Clear()
	for _, pv := range playerViews {
		pv.Draw(utils.Renderer, dest, time.Duration(utils.Delta))
	}

	flag[1] = false
}

func cleanupRenderer() {
	for _, pv := range playerViews {
		pv.Texture.Destroy()
	}
	utils.Renderer.Destroy()
}
