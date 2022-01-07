package main

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Mouse struct {
	leftButton  bool
	rightButton bool
	x, y        int
}

func GetMouseState() Mouse {
	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()
	leftButton := mouseButtonState & sdl.ButtonLMask()
	rightButton := mouseButtonState & sdl.ButtonRMask()
	var userInput Mouse
	userInput.x = int(mouseX)
	userInput.y = int(mouseY)
	userInput.leftButton = !(leftButton == 0)
	userInput.rightButton = !(rightButton == 0)
	return userInput
}
