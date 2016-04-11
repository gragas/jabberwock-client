package mainMenu

import (
	"github.com/gragas/go-sdl2/sdl"	
	"github.com/gragas/jabberwock-client/utils"
)

func Init(ip string, port int, quiet bool, debug bool, serverDebug bool) {
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
