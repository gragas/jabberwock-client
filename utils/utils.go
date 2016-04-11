package utils

import (
	"github.com/gragas/go-sdl2/sdl"	
)

var Delta uint32
var Running bool
var Window *sdl.Window
var Surface *sdl.Surface
var Loop LoopFuncs

type LoopFuncs struct {
	PollEvents func()
	Update func()
	Draw func()
}
