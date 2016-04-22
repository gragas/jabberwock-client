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

	registrationFailure := func(address string, malformed bool) {
		fmt.Printf("CLIENT: Could not register with server at %s!\n", address)
		if malformed {
			fmt.Printf("CLIENT: Additional info: received malformed response.\n")
		}
	}

	address := ip + ":" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		registrationFailure(address, false)
		return nil, nil, nil, false
	}
	sendstr := string(protocol.Register)
	sendstr += sendPlayer.String() + string(protocol.EndOfMessage)
	fmt.Fprintf(conn, sendstr)
	reader := bufio.NewReader(conn)
	status, err := reader.ReadString(byte(protocol.EndOfMessage))
	if err != nil || protocol.Code(status[0]) != protocol.Success {
		registrationFailure(address, false)
		return nil, nil, nil, false
	}
	if len(status) < 2 {
		registrationFailure(address, true)
		return nil, nil, nil, false
	}
	if debug {
		fmt.Printf("CLIENT: Successfully registered client with server at %s\n", address)
		fmt.Printf("CLIENT: Response from server was '%v'\n",
			protocol.Code(status[0]).String()+status[1:len(status)-1])
	}
	var recvPlayer player.Player
	err = recvPlayer.FromBytes([]byte(status[1 : len(status)-1]))
	if err != nil {
		panic(err)
	}
	return conn, reader, &recvPlayer, true
}
