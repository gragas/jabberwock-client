package utils

import (
	"fmt"
	"github.com/gragas/go-sdl2/sdl"	
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

func RegisterClient(ip string, port int, debug bool) {
	if debug {
		fmt.Printf("Registering client with server at \033[0;31m")
		fmt.Printf("%v\033[0m:\033[0;34m%v\033[0m\n", ip, port)
	}
}
