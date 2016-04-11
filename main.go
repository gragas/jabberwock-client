package main

import (
	"flag"
	"github.com/gragas/go-sdl2/sdl"
	"github.com/gragas/jabberwock-lib/consts"
	"github.com/gragas/jabberwock-client/debug"
	"github.com/gragas/jabberwock-client/mainMenu"
	"github.com/gragas/jabberwock-client/utils"
	"time"
)

const (
	wString = "specifies the window width"
	hString = "specifies the window height"
	titleString = "specifies the window title"
	quietString = "specifies whether to be quiet"
	debugString = "specifies whether debug mode is enabled"
	serverDebugString = "specifies whether server debug mode is enabled"
	ipString    = "specifies the IP address this jabberwock server will bind to"
	portString  = "specifies the port this jabberwock server will bind to"
)

var ip, windowTitle string
var windowWidth, windowHeight, port int
var debugMode, quietMode, serverDebugMode bool

func main() {
	parseFlags()
	utils.Window, utils.Surface = initialize()
	defer utils.Window.Destroy()

	ticks := time.Now()
	for utils.Running {
		utils.Loop.PollEvents()
		utils.Loop.Update()
		utils.Loop.Draw()
		
		utils.Delta = time.Now().Sub(ticks).Nanoseconds()
		ticks = time.Now()
		if utils.Delta < consts.TicksPerFrame {
			time.Sleep(time.Duration(consts.TicksPerFrame - utils.Delta))
		}
		utils.Window.UpdateSurface()
	}
	sdl.Quit()
}

func parseFlags() {
	flag.IntVar(&windowWidth, "w", 800, wString)
	flag.IntVar(&windowHeight, "h", 600, hString)
	flag.StringVar(&windowTitle, "title", "Beware the Jabberwock!", titleString)
	flag.BoolVar(&quietMode, "quiet", false, quietString)
	flag.BoolVar(&debugMode, "debug", false, debugString)
	flag.BoolVar(&serverDebugMode, "sdebug", false, serverDebugString)
	flag.StringVar(&ip, "ip", "127.0.0.1", ipString)
	flag.IntVar(&port, "port", 5000, portString)
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
		debug.Init(ip, port, quietMode, debugMode, serverDebugMode)
	} else {
		mainMenu.Init(ip, port, quietMode, debugMode, serverDebugMode)
	}

	return window, surface
}
