package debug

import (
	"github.com/gragas/go-sdl2/sdl"	
	"github.com/gragas/jabberwock-client/utils"
	"github.com/gragas/jabberwock-server/game"
)

func Init(ip string, port int, quiet bool, debug bool, serverDebug bool) {
	done := make(chan string)
	go game.StartGame(ip, port, quiet, serverDebug, done)
	<-done // make sure the server is setup, then register with it
	utils.RegisterClient(ip, port, debug)
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
