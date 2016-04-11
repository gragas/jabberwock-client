package debug

import (
	"fmt"
	"github.com/gragas/go-sdl2/sdl"	
	"github.com/gragas/jabberwock-client/utils"
	"github.com/gragas/jabberwock-server/game"
)

func Init(ip string, port int, quiet bool, debug bool, serverDebug bool) {
	done := make(chan string)
	go game.StartGame(ip, port, quiet, serverDebug, done)
	<-done // make sure the server is setup, then register with it
	var registered bool
	var attempts int
	for ; attempts < 5; attempts++ {
		registered = utils.RegisterClient(ip, port, debug)
		if registered {
			break
		}
	}
	if !registered {
		fmt.Printf("Could not register client after %v attempts.\n", attempts)
	}
	utils.Loop = utils.LoopFuncs{pollEvents, update, draw}
}

func pollEvents() {
	for {
		if event := sdl.PollEvent(); event != nil {
			switch event.(type) {
			case *sdl.QuitEvent:
				utils.Running = false
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
