package main

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"os"
)

type Cell struct {
	tex *sdl.Texture
}

type Box struct {
	tex  *sdl.Texture
	x, y int32
}

func StartingScreen() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("help")
	}
	window, err := sdl.CreateWindow("final project", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, screenWidth, screenHeight,
		sdl.WINDOW_OPENGL)
	if err != nil {
		fmt.Println("Could not create windows")
		return
	}

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println("Trouble with renderer", err)
		return
	}

	cell, err := CreateCell(renderer)
	if err != nil {
		fmt.Println("Trouble creating cell:", err)
		return
	}
	nucbox, err := CreateBox("images/label_boxes/nucleus_box.bmp", renderer)
	if err != nil {
		fmt.Println("Trouble creating nucleus box:", err)
		return
	}
	mitobox, err := CreateBox("images/label_boxes/mitochondria_box.bmp", renderer)
	if err != nil {
		fmt.Println("Trouble creating mitochondria box:", err)
		return
	}
	rerbox, err := CreateBox("images/label_boxes/RER_box.bmp", renderer)
	if err != nil {
		fmt.Println("Trouble creating Rough ER box:", err)
		return
	}
	serbox, err := CreateBox("images/label_boxes/SER_box.bmp", renderer)
	if err != nil {
		fmt.Println("Trouble creating Smooth ER box:", err)
		return
	}
	plasmabox, err := CreateBox("images/label_boxes/plasma_mem_box.bmp", renderer)
	if err != nil {
		fmt.Println("Trouble creating Plasma membrane box:", err)
		return
	}
	golgibox, err := CreateBox("images/label_boxes/golgi_box.bmp", renderer)
	if err != nil {
		fmt.Println("Trouble creating Golgi box:", err)
		return
	}
	game_box, err := CreateBox("images/play_game.bmp", renderer)
	if err != nil {
		fmt.Println("Trouble creating Play Game Box:", err)
		return
	}
	exit_box, err := CreateBox("images/exit.bmp", renderer)
	if err != nil {
		fmt.Println("Trouble creating cell:", err)
		return
	}
	currentMouseState := GetMouseState()
	prevMouseState := currentMouseState
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		currentMouseState = GetMouseState()
		if !currentMouseState.leftButton && prevMouseState.leftButton {
			mouseX, mouseY, _ := sdl.GetMouseState()
			// fmt.Println("Left Click", mouseX, mouseY)
			if (275 <= mouseX) && (mouseX <= 354) && (140 <= mouseY) && (mouseY <= 176) {
				window.Destroy()
				DisplayInformationWindow("nucleus")
				return
			} else if (266 <= mouseX) && (mouseX <= 395) && (562 <= mouseY) && (mouseY <= 589) {
				window.Destroy()
				DisplayInformationWindow("mitochondria")
				return
			} else if (405 <= mouseX) && (mouseX <= 598) && (205 <= mouseY) && (mouseY <= 253) {
				window.Destroy()
				DisplayInformationWindow("roughER")
				return
			} else if (21 <= mouseX) && (mouseX <= 154) && (480 <= mouseY) && (mouseY <= 580) {
				window.Destroy()
				DisplayInformationWindow("smoothER")
				return
			} else if (481 <= mouseX) && (mouseX <= 600) && (277 <= mouseY) && (mouseY <= 342) {
				window.Destroy()
				DisplayInformationWindow("plasma")
				return
			} else if (371 <= mouseX) && (mouseX <= 482) && (512 <= mouseY) && (mouseY <= 559) {
				window.Destroy()
				DisplayInformationWindow("golgi")
				return
			} else if (18 <= mouseX) && (mouseX <= 216) && (60 <= mouseY) && (mouseY <= 139) {
				return
			} else if (443 <= mouseX) && (mouseX <= 596) && (653 <= mouseY) && (mouseY <= 706) {
				fmt.Println("Bye")
				os.Exit(0)
			}

		}

		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.Clear()
		cell.DrawCell(renderer)
		nucbox.DrawBox(275, 140, 79, 36, renderer)
		mitobox.DrawBox(265, 560, 132, 31, renderer)
		rerbox.DrawBox(404, 203, 197, 53, renderer)
		serbox.DrawBox(20, 488, 137, 95, renderer)
		plasmabox.DrawBox(481, 276, 152, 68, renderer)
		golgibox.DrawBox(371, 512, 113, 50, renderer)
		game_box.DrawBox(16, 59, 200, 81, renderer)
		exit_box.DrawBox(425, 630, 201, 100, renderer)
		renderer.Present()

		prevMouseState = currentMouseState
	}
	defer renderer.Destroy()
	defer window.Destroy()
}

func CreateCell(renderer *sdl.Renderer) (c Cell, err error) {

	img, err := sdl.LoadBMP("images/cell_image.bmp")
	if err != nil {
		fmt.Println("Had trouble loading the image", err)
		return
	}
	c.tex, err = renderer.CreateTextureFromSurface(img)
	if err != nil {
		fmt.Println("Had trouble texture", err)
		return
	}
	defer img.Free()
	return c, nil
}

func CreateBox(filename string, renderer *sdl.Renderer) (box Box, err error) {
	img, err := sdl.LoadBMP(filename)
	if err != nil {
		fmt.Println("Had trouble loading the box image", err)
		return
	}
	box.tex, err = renderer.CreateTextureFromSurface(img)
	if err != nil {
		fmt.Println("Had trouble box texture", err)
		return
	}
	defer img.Free()
	return box, nil
}

func (c *Cell) DrawCell(renderer *sdl.Renderer) {
	renderer.Copy(c.tex, &sdl.Rect{X: 0, Y: 0, W: 512, H: 465}, &sdl.Rect{X: 50, Y: 150, W: 512, H: 465})
}

func (box *Box) DrawBox(x, y, picw, pich int32, renderer *sdl.Renderer) {
	box.x = x
	box.y = y
	renderer.Copy(box.tex, &sdl.Rect{X: 0, Y: 0, W: picw, H: pich}, &sdl.Rect{X: box.x, Y: box.y, W: picw, H: pich})
}
