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
var entities []entity.Entity
var players []*player.Player

func Init(ip string, port int, initialized chan bool, quiet bool, debug bool, serverDebug bool) {
	done := make(chan bool)
	go game.StartGame(ip, port, quiet, serverDebug, done)
	<-done
	sendPlayer := player.NewDefaultPlayer()
	var registered bool
	var attempts int
	for ; attempts < 5; attempts++ {
		conn, reader, clientPlayer, registered = utils.RegisterClient(ip, port, sendPlayer, debug)
		if registered {
			players = append(players, clientPlayer)
			entities = append(entities, clientPlayer)
			break
		}
	}
	if !registered {
		fmt.Printf("CLIENT: Failed to register after %v attempts.\n", attempts)
		panic(errors.New("ERROR: Client failed to register with server.\n"))
	}
	surf, rect := entity.NewDefaultEntityView(clientPlayer)
	clientPlayerView = &player.PlayerView{PlayerPtr: clientPlayer, Surface: surf, Rect: rect}
	if debug {
		fmt.Printf("CLIENT: Client player is %v\n", entity.ShortString(clientPlayer))
	}
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
			panic(err)
		}
		for _, p := range players {
			if clientPlayer.GetID() == p.GetID() {
					clientPlayer = p
			}
		}
	default:
		fmt.Printf("CLIENT: Invalid protocol.Code.\ncode: %v\nmsg: %v\n", protocol.Code(msg[0]), msg)
	}
}

func draw(dest *sdl.Surface) {
	utils.Surface.FillRect(nil, sdl.MapRGBA(utils.Surface.Format, 0xFF, 0xEE, 0xCC, 0xFF))
	if clientPlayerView != nil {
		clientPlayerView.Draw(dest)
	}
}
