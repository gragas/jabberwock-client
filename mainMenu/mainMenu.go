package mainMenu

import (
	"fmt"
	"github.com/gragas/go-sdl2/sdl"	
	"github.com/gragas/jabberwock-client/utils"
)

var receiverToHandler chan string

func Init(ip string, port int, initialized chan bool, quiet bool, debug bool, serverDebug bool) {
	utils.Loop = utils.LoopFuncs{pollEvents, draw}
	initialized <- true
}

func quit() {
	utils.Running = false
}

func pollEvents() {
	for {
		if event := sdl.PollEvent(); event != nil {
			switch event.(type) {
			case *sdl.QuitEvent:
				utils.Running = false
			case *sdl.KeyDownEvent:
				sym := event.(*sdl.KeyDownEvent).Keysym.Sym
				switch sym {
				case sdl.K_ESCAPE:
					quit()
				default:
					fmt.Println(sym)
				}
			}
		} else {
			break
		}
	}
}

func update(msg string, debug bool) {
	fmt.Println(msg)
}

func receiver(debug bool) {

}

func handleMessage(debug bool) {
	select {
	case msg := <-receiverToHandler:
		update(msg, debug)
	default:
		/* if debug {
		     	 *	fmt.Printf("Nothing received!\n")
			     * }
		*/
	}
}

func draw(dest *sdl.Surface) {
	
}
