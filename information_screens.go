package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

type Screen struct {
	tex *sdl.Texture
}

func DisplayInformationWindow(window_name string) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		fmt.Println("help")
	}
	window, err := sdl.CreateWindow(window_name, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, screenWidth, screenHeight,
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
	defer renderer.Destroy()
	defer window.Destroy()
	if window_name == "nucleus" {
		nucInfo, err := CreateScreen("images/label_boxes/nucleus_description.bmp", renderer)
		if err != nil {
			fmt.Println("Trouble creating cell:", err)
			return
		}
		WindowsScreenInfo(nucInfo, renderer, window, 600, 398, 600, 398)
	} else if window_name == "mitochondria" {
		mitoInfo, err := CreateScreen("images/label_boxes/mitochondria_description.bmp", renderer)
		if err != nil {
			fmt.Println("Trouble creating cell:", err)
			return
		}
		WindowsScreenInfo(mitoInfo, renderer, window, 696, 536, 600, 500)
	} else if window_name == "roughER" {
		rerInfo, err := CreateScreen("images/label_boxes/rough_er_description.bmp", renderer)
		if err != nil {
			fmt.Println("Trouble creating cell:", err)
			return
		}
		WindowsScreenInfo(rerInfo, renderer, window, 675, 441, 600, 441)
	} else if window_name == "smoothER" {
		serInfo, err := CreateScreen("images/label_boxes/smooth_er_description.bmp", renderer)
		if err != nil {
			fmt.Println("Trouble creating cell:", err)
			return
		}
		WindowsScreenInfo(serInfo, renderer, window, 675, 441, 600, 441)
	} else if window_name == "plasma" {
		plasmaInfo, err := CreateScreen("images/label_boxes/plasma_membrane_description.bmp", renderer)
		if err != nil {
			fmt.Println("Trouble creating cell:", err)
			return
		}
		WindowsScreenInfo(plasmaInfo, renderer, window, 675, 292, 550, 292)
	} else if window_name == "golgi" {
		rerInfo, err := CreateScreen("images/label_boxes/golgi_description.bmp", renderer)
		if err != nil {
			fmt.Println("Trouble creating cell:", err)
			return
		}
		WindowsScreenInfo(rerInfo, renderer, window, 675, 441, 600, 441)
	}

}

func WindowsScreenInfo(boxInfo Screen, renderer *sdl.Renderer, window *sdl.Window, W, H, picScaleW, picScaleH int32) {
	currentMouseState := GetMouseState()
	prevMouseState := currentMouseState
	return_box, err := CreateBox("images/return_arrow.bmp", renderer)
	if err != nil {
		fmt.Println("Trouble creating cell:", err)
		return
	}
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
			fmt.Println("Left Click", mouseX, mouseY)
			if (465 <= mouseX) && (mouseX <= 600) && (50 <= mouseY) && (mouseY <= 150) {
				window.Destroy()
				StartingScreen()
				return
			}

		}
		renderer.SetDrawColor(255, 255, 255, 255)
		renderer.Clear()
		boxInfo.DrawScreen(renderer, W, H, picScaleW, picScaleH)
		return_box.DrawBox(465, 50, 135, 100, renderer)
		renderer.Present()
		prevMouseState = currentMouseState
	}
}

func CreateScreen(filename string, renderer *sdl.Renderer) (s Screen, err error) {
	img, err := sdl.LoadBMP(filename)
	if err != nil {
		fmt.Println("Had trouble loading the image", err)
		return
	}
	s.tex, err = renderer.CreateTextureFromSurface(img)
	if err != nil {
		fmt.Println("Had trouble texture", err)
		return
	}
	defer img.Free()
	return s, nil
}

func (s *Screen) DrawScreen(renderer *sdl.Renderer, picW, picH, picScaleW, picScaleH int32) {
	renderer.Copy(s.tex, &sdl.Rect{X: 0, Y: 0, W: picW, H: picH}, &sdl.Rect{X: 0, Y: 200, W: picScaleW, H: picScaleH})
}
