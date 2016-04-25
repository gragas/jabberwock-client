package utils

import (
	"bufio"
	"fmt"
	"github.com/gragas/go-sdl2/sdl"
	"github.com/gragas/jabberwock-lib/player"
	"github.com/gragas/jabberwock-lib/protocol"
	"net"
	"strconv"
)

var Delta int64
var Running bool
var Window *sdl.Window
var Surface *sdl.Surface
var Renderer *sdl.Renderer
var Loop LoopFuncs

type LoopFuncs struct {
	PollEvents    func()
	Draw          func(dest *sdl.Surface)
}

func RegisterClient(ip string, port int, sendPlayer *player.Player, debug bool) (net.Conn, *bufio.Reader, *player.Player, bool) {
	if debug {
		fmt.Printf("CLIENT: Registering with server at \033[0;31m")
		fmt.Printf("%v\033[0m:\033[0;34m%v\033[0m\n", ip, port)
	}

	/* dial the server */
	address := ip + ":" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Printf("CLIENT ERROR: %v\n", err)
		return nil, nil, nil, false
	}
	/*******************/

	/* send it a registration message */
	fmt.Fprintf(conn, string(protocol.Register) + sendPlayer.String() + string(protocol.EndOfMessage))
	/************************************/

	/* get the new player back */
	reader := bufio.NewReader(conn)
	msg, err := reader.ReadString(byte(protocol.EndOfMessage))
	if err != nil {
		fmt.Printf("CLIENT ERROR: %v\n", err)
		return nil, nil, nil, false
	}
	if len(msg) < 2 {
		fmt.Printf("CLIENT: Registration response was too short: %v\n", msg)
		return nil, nil, nil, false
	}
	if protocol.Code(msg[0]) != protocol.Success {
		fmt.Printf("CLIENT: Recieved invalid registration response: %v\n", msg)
		return nil, nil, nil, false
	}
	contents := msg[1:len(msg)-1]
	recvPlayer := new(player.Player)
	err = recvPlayer.FromBytes([]byte(contents))
	if err != nil {
		panic(err)
	}
	/***************************/

	/* handshake with the server */
	fmt.Fprintf(conn, string(protocol.Handshake) + string(protocol.EndOfMessage))
	/*****************************/
	if debug {
		fmt.Printf("CLIENT: Successfully registered client with server at %s\n", address)
	}
	return conn, reader, recvPlayer, true
}
