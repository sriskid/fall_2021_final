package main

import (
	"final_project/game"
	"final_project/ui"
)

const (
	screenWidth  = 600
	screenHeight = 800
)

func main() {
	StartingScreen()
	game := game.NewGame(1)
	go func() {
		game.Run()
	}()
	ui := ui.NewUI(game.InputChan, game.LevelChans[0])
	ui.Run()

}
