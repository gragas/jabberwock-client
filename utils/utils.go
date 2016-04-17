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
	PollEvents func()
	Update     func()
	Draw       func()
}

func RegisterClient(ip string, port int, player player.Player,
	debug bool) (net.Conn, bool) {
	if debug {
		fmt.Printf("CLIENT: Registering with server at \033[0;31m")
		fmt.Printf("%v\033[0m:\033[0;34m%v\033[0m\n", ip, port)
	}

	registrationFailure := func(address string, malformed bool) (net.Conn, bool) {
		fmt.Printf("CLIENT: Could not register with server at %s!\n", address)
		if malformed {
			fmt.Printf("CLIENT: Additional info: received malformed response.\n")
		}
		return nil, false
	}

	address := ip + ":" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return registrationFailure(address, false)
	}
	sendstr := string(protocol.Register)
	sendstr += player.String() + string(protocol.EndOfMessage)
	fmt.Fprintf(conn, sendstr)
	status, err := bufio.NewReader(conn).ReadString(byte(protocol.EndOfMessage))
	if err != nil || protocol.Code(status[0]) != protocol.Success {
		return registrationFailure(address, false)
	}
	if len(status) < 2 {
		return registrationFailure(address, true)
	}
	if debug {
		fmt.Printf("CLIENT: Successfully registered client with server at %s\n",
			address)
		fmt.Printf("CLIENT: Response from server was '%v'\n",
			protocol.Code(status[0]).String()+status[1:len(status)-1])
	}
	return conn, true
}
