package debug

import (
	"github.com/gragas/go-sdl2/sdl"	
	"github.com/gragas/jabberwock-client/utils"
)

func PollEvents() {
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

func Update() {

}

func Draw() {

}
