package utils

import (
	"bufio"
	"fmt"
	"github.com/gragas/go-sdl2/sdl"	
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
	Update func()
	Draw func()
}

func RegisterClient(ip string, port int, debug bool) (net.Conn, bool) {
	if debug {
		fmt.Printf("CLIENT: Registering with server at \033[0;31m")
		fmt.Printf("%v\033[0m:\033[0;34m%v\033[0m\n", ip, port)
	}
	
	registrationFailure := func(address string) (net.Conn, bool) {
		fmt.Printf("CLIENT: Could not register with server at %s!\n", address)
		return nil, false
	}

	address := ip + ":" + strconv.Itoa(port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return registrationFailure(address)
	}
	fmt.Fprintf(conn, "REGISTER\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		return registrationFailure(address)
	}
	if debug {
		fmt.Printf("CLIENT: Successfully registered client with server at %s\n",
			address)
		fmt.Printf("CLIENT: Response from server was %v\n", status[:8])
	}
	return conn, true
}
