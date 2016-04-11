package main

import (
	"flag"
	"github.com/gragas/go-sdl2/sdl"
	"github.com/gragas/jabberwock-lib/consts"
	"github.com/gragas/jabberwock-client/debug"
	"github.com/gragas/jabberwock-client/mainMenu"
	"github.com/gragas/jabberwock-client/utils"
)

var windowTitle string
var windowWidth, windowHeight int
var debugMode bool

func main() {
	parseFlags()
	utils.Window, utils.Surface = initialize()
	defer utils.Window.Destroy()

	ticks := uint32(0)

	for utils.Running {
		utils.Loop.PollEvents()
		utils.Loop.Update()
		utils.Loop.Draw()

		utils.Delta = sdl.GetTicks() - ticks
		ticks = sdl.GetTicks()
		if utils.Delta < consts.TicksPerFrame {
			sdl.Delay(consts.TicksPerFrame - utils.Delta)
		}
		utils.Window.UpdateSurface()
	}
	sdl.Quit()
}

func parseFlags() {
	flag.IntVar(&windowWidth, "w", 800, "specifies the window width")
	flag.IntVar(&windowHeight, "h", 600, "specifies the window height")
	flag.StringVar(&windowTitle, "t", "jabberwock", "specifies the window title")
	flag.BoolVar(&debugMode, "d", true, "specifies whether debug mode is enabled")
	flag.Parse()
}

func initialize() (*sdl.Window, *sdl.Surface) {
	sdl.Init(sdl.INIT_EVERYTHING)
	window, err := sdl.CreateWindow(windowTitle,
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	utils.Running = true
	if debugMode {
		utils.Loop = utils.LoopFuncs{
			debug.PollEvents,
			debug.Update,
			debug.Draw}
	} else {
		utils.Loop = utils.LoopFuncs{
			mainMenu.PollEvents,
			mainMenu.Update,
			mainMenu.Draw}
	}

	return window, surface
}
