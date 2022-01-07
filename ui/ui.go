package ui

import (
	"bufio"
	"final_project/game"
	"fmt"
	"image/png"
	"math/rand"
	"os"
	"strconv"
	"strings"

	"github.com/veandco/go-sdl2/mix"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const itemSizeratio = 0.033

type Mouse struct {
	leftButton  bool
	rightButton bool
	pos         game.Pos
}

func GetMouseState() *Mouse {
	mouseX, mouseY, mouseButtonState := sdl.GetMouseState()
	leftButton := mouseButtonState & sdl.ButtonLMask()
	rightButton := mouseButtonState & sdl.ButtonRMask()
	var userInput Mouse
	userInput.pos = game.Pos{int(mouseX), int(mouseY)}
	userInput.leftButton = !(leftButton == 0)
	userInput.rightButton = !(rightButton == 0)
	return &userInput
}

type uiState int

const (
	uiMain uiState = iota
	uiInventory
	uiCodon
)

type ui struct {
	state               uiState
	draggedItem         *game.Item
	winWidth            int
	winHeight           int
	renderer            *sdl.Renderer
	window              *sdl.Window
	textureAtlas        *sdl.Texture
	textureIndex        map[rune][]sdl.Rect
	prevKeyboardState   []uint8
	keyboardState       []uint8
	centerX             int
	centerY             int
	r                   *rand.Rand
	levelChan           chan *game.Level
	inputChan           chan *game.Input
	fontMedium          *ttf.Font
	fontSmall           *ttf.Font
	fontLarge           *ttf.Font
	str2Textsmall       map[string]*sdl.Texture
	str2Textmedium      map[string]*sdl.Texture
	str2Textlarge       map[string]*sdl.Texture
	eventBackground     *sdl.Texture
	inventoryBackground *sdl.Texture
	currentMouseState   *Mouse
	prevMouseState      *Mouse
}

func NewUI(inputChan chan *game.Input, levelChan chan *game.Level) *ui {
	ui := &ui{}
	ui.state = uiMain
	ui.str2Textsmall = make(map[string]*sdl.Texture)
	ui.str2Textmedium = make(map[string]*sdl.Texture)
	ui.str2Textlarge = make(map[string]*sdl.Texture)
	ui.inputChan = inputChan
	ui.levelChan = levelChan
	ui.r = rand.New(rand.NewSource(1))
	ui.winHeight = 650
	ui.winWidth = 1000
	window, err := sdl.CreateWindow("Cell RPG", 200, 200, int32(ui.winWidth), int32(ui.winHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		panic("Trouble creating the Window")
	}
	ui.window = window

	ui.renderer, err = sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic("Trouble creating renderer")
	}

	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	ui.textureAtlas = ui.ImgFromFileToTexture("ui/assets/game_tiles.png")
	ui.LoadTextureIndex()
	ui.keyboardState = sdl.GetKeyboardState()
	ui.prevKeyboardState = make([]uint8, len(ui.keyboardState))
	for i, v := range ui.keyboardState {
		ui.prevKeyboardState[i] = v
	}
	ui.centerX = -1
	ui.centerY = -1

	ui.fontSmall, err = ttf.OpenFont("ui/assets/Kingthings_Foundation.ttf", int(float64(ui.winWidth)*0.015))
	if err != nil {
		panic("Trouble opening fonts")
	}

	ui.fontMedium, err = ttf.OpenFont("ui/assets/Kingthings_Foundation.ttf", 32)
	if err != nil {
		panic("Trouble opening fonts")
	}

	ui.fontLarge, err = ttf.OpenFont("ui/assets/Kingthings_Foundation.ttf", 64)
	if err != nil {
		panic("Trouble opening fonts")
	}

	ui.eventBackground = ui.GetSinglePixelTex(sdl.Color{0, 0, 0, 128})
	ui.eventBackground.SetBlendMode(sdl.BLENDMODE_BLEND)

	ui.inventoryBackground = ui.GetSinglePixelTex(sdl.Color{0, 0, 255, 128})
	ui.inventoryBackground.SetBlendMode(sdl.BLENDMODE_BLEND)

	err = mix.OpenAudio(22050, mix.DEFAULT_FORMAT, 2, 4096)
	if err != nil {
		panic(err)
	}
	music, err := mix.LoadMUS("ui/assets/Castlevania_BGM.ogg")
	if err != nil {
		panic("Having trouble with the Castlevania tunes")
	}
	music.Play(-1)
	return ui
}

type fontSize int

const (
	fontSmall fontSize = iota
	fontMedium
	fontLarge
)

func (ui *ui) StringTexture(s string, color sdl.Color, size fontSize) *sdl.Texture {

	var font *ttf.Font
	switch size {
	case fontSmall:
		font = ui.fontSmall
		tex, exists := ui.str2Textsmall[s]
		if exists {
			return tex
		}
	case fontMedium:
		font = ui.fontMedium
		tex, exists := ui.str2Textmedium[s]
		if exists {
			return tex
		}
	case fontLarge:
		font = ui.fontLarge
		tex, exists := ui.str2Textlarge[s]
		if exists {
			return tex
		}
	}
	fontSurface, err := font.RenderUTF8Blended(s, color)
	if err != nil {
		panic("Trouble rendering font")
	}

	tex, err := ui.renderer.CreateTextureFromSurface(fontSurface)
	if err != nil {
		panic("Trouble creating font texture")
	}

	switch size {
	case fontSmall:
		ui.str2Textsmall[s] = tex
	case fontMedium:
		ui.str2Textmedium[s] = tex
	case fontLarge:
		ui.str2Textlarge[s] = tex
	}

	return tex
}

func (ui *ui) LoadTextureIndex() {
	ui.textureIndex = make(map[rune][]sdl.Rect)
	infile, err := os.Open("ui/assets/atlas_index.txt")
	if err != nil {
		panic("Trouble reading index file")
	}
	scanner := bufio.NewScanner(infile)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		tileRune := rune(line[0])
		xy := line[1:]
		splitXYC := strings.Split(xy, ",")
		x, err := strconv.ParseInt(strings.TrimSpace(splitXYC[0]), 10, 64)
		if err != nil {
			panic("Trouble getting x coordinates of tile")
		}
		y, err := strconv.ParseInt(strings.TrimSpace(splitXYC[1]), 10, 64)
		if err != nil {
			panic("Trouble getting y coordinates of tile")
		}
		variationCount, err := strconv.ParseInt(strings.TrimSpace(splitXYC[2]), 10, 64)
		if err != nil {
			panic("Trouble getting variation count")
		}
		var rects []sdl.Rect
		for i := int64(0); i < variationCount; i++ {
			rects = append(rects, sdl.Rect{int32(x * 32), int32(y * 32), 32, 32})
			x++
			if x > 62 {
				x = 0
				y++
			}
		}
		ui.textureIndex[tileRune] = rects

	}
}

func (ui *ui) ImgFromFileToTexture(filename string) *sdl.Texture {
	inFile, err := os.Open(filename)
	if err != nil {
		panic("Trouble opening the file")
	}
	defer inFile.Close()

	img, err := png.Decode(inFile)
	if err != nil {
		panic("Trouble opening PNG")
	}

	w := img.Bounds().Max.X
	h := img.Bounds().Max.Y

	pixels := make([]byte, w*h*4)
	bIndex := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			pixels[bIndex] = byte(r / 256)
			bIndex++
			pixels[bIndex] = byte(g / 256)
			bIndex++
			pixels[bIndex] = byte(b / 256)
			bIndex++
			pixels[bIndex] = byte(a / 256)
			bIndex++
		}
	}
	tex, err := ui.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, int32(w), int32(h))
	if err != nil {
		panic("Trouble with texture")
	}
	tex.Update(nil, pixels, w*4)

	err = tex.SetBlendMode(sdl.BLENDMODE_BLEND)
	if err != nil {
		panic("Having trouble with alpha blending")
	}
	return tex
}

func init() {

	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println("We have a problem here")
		return
	}
	err = ttf.Init()
	if err != nil {
		panic("We have a big problem initializing the font")
	}

	err = mix.Init(mix.INIT_OGG)
	//Bug with SDL. ignore the error
	/* if err != nil {
		panic("Trouble getting in some tunes")
	}*/

}

//Depending on your preferred camera style you can comment out if statement
//and "camera follow" if statements if you want the camera to always follow you
//or do not comment those out if you want the camera to follow you after a certain
//distance
func (ui *ui) Draw(level *game.Level) {
	if ui.centerX == -1 && ui.centerY == -1 {
		ui.centerX = level.Player.X
		ui.centerY = level.Player.Y
	}

	//Camera follow
	limit := 5
	if level.Player.X > ui.centerX+limit {
		diff := level.Player.X - (ui.centerX + limit)
		ui.centerX += diff
	} else if level.Player.X < ui.centerX-limit {
		diff := (ui.centerX - limit) - level.Player.X
		ui.centerX -= diff
	} else if level.Player.Y > ui.centerY+limit {
		diff := level.Player.Y - (ui.centerY + limit)
		ui.centerY += diff
	} else if level.Player.Y < ui.centerY-limit {
		diff := (ui.centerY - limit) - level.Player.Y
		ui.centerY -= diff
	}
	offsetX := int32((ui.winWidth / 2) - ui.centerX*32)
	offsetY := int32((ui.winHeight / 2) - ui.centerY*32)

	ui.renderer.Clear()
	ui.r.Seed(1)
	for y, row := range level.Organelle {
		for x, tile := range row {
			if tile.Rune != game.Blank {
				srcRects := ui.textureIndex[tile.Rune]
				srcRect := srcRects[ui.r.Intn(len(srcRects))]
				if tile.Visible || tile.Seen {
					dstRect := sdl.Rect{int32(x*32) + offsetX, int32(y*32) + offsetY, 32, 32}
					pos := game.Pos{x, y}
					if level.Debug[pos] {
						ui.textureAtlas.SetColorMod(128, 0, 0)
					} else if tile.Seen && !tile.Visible {
						ui.textureAtlas.SetColorMod(128, 128, 128)
					} else {
						ui.textureAtlas.SetColorMod(255, 255, 255)
					}
					ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)

					if tile.OverlayRune != game.Blank {
						srcRect := ui.textureIndex[tile.OverlayRune][0]
						ui.renderer.Copy(ui.textureAtlas, &srcRect, &dstRect)
					}
				}
			}
		}

	}
	for pos, items := range level.Items {
		if level.Organelle[pos.Y][pos.X].Visible {
			for _, item := range items {
				itemSrcrect := ui.textureIndex[item.Rune][0]
				ui.renderer.Copy(ui.textureAtlas, &itemSrcrect, &sdl.Rect{int32(pos.X)*32 + offsetX, int32(pos.Y)*32 + offsetY, 32, 32})
			}
		}
	}
	//Player Sprite coordinates on PNG {21,59}
	playerSrcRect := ui.textureIndex['@'][0]
	ui.renderer.Copy(ui.textureAtlas, &playerSrcRect, &sdl.Rect{int32(level.Player.X)*32 + offsetX, int32(level.Player.Y)*32 + offsetY, 32, 32})
	text := ui.StringTexture("The Cell", sdl.Color{0, 255, 0, 0}, fontMedium)
	_, _, w, h, _ := text.Query()
	ui.renderer.Copy(text, nil, &sdl.Rect{0, 0, w, h})
	textStartY := int32(float32(ui.winHeight) * .68)
	textWidth := int32(float32(ui.winWidth) * .25)
	ui.renderer.Copy(ui.eventBackground, nil, &sdl.Rect{0, textStartY, textWidth, int32(ui.winHeight) - textStartY})
	i := level.EventPos
	count := 0
	_, fontSizeY, _ := ui.fontSmall.SizeUTF8("A")
	for {
		event := level.Events[i]
		if event != "" {
			tex := ui.StringTexture(event, sdl.Color{255, 0, 0, 0}, fontSmall)
			_, _, w, h, _ := tex.Query()
			ui.renderer.Copy(tex, nil, &sdl.Rect{5, int32(count*fontSizeY) + textStartY, w, h})
		}
		i = (i + 1) % (len(level.Events))
		count++
		if i == level.EventPos {
			break
		}
	}
	missionTextY := int32(float32(ui.winHeight) * .10)
	// missionTextendY := int32(float32(ui.winHeight) * .50)
	// missionTextWidth := int32(float32(ui.winWidth) * .25)

	x := level.DesPos
	// countDes := 0
	for {
		description := level.Descriptions[x]
		if description != "" {
			brokenString := strings.Split(description, ". ")
			for i, sentence := range brokenString {
				tex := ui.StringTexture(sentence, sdl.Color{0, 0, 255, 0}, fontSmall)
				_, _, w, h, _ := tex.Query()
				ui.renderer.Copy(tex, nil, &sdl.Rect{5, int32(i*fontSizeY) + missionTextY, w, h})
			}
		}
		x = (x + 1) % (len(level.Descriptions))
		// countDes++
		if x == level.DesPos {
			break
		}
	}

	inventoryStart := int32(float32(ui.winWidth) * 0.9)
	inventoryWidth := int32(ui.winWidth) - inventoryStart
	itemSize := int32(itemSizeratio * float32(ui.winWidth))
	ui.renderer.Copy(ui.inventoryBackground, nil, &sdl.Rect{inventoryStart, int32(ui.winHeight) - itemSize, inventoryWidth, itemSize})
	//Inventory UI
	items := level.Items[level.Player.Pos]
	for i, item := range items {
		itemSrcrect := ui.textureIndex[item.Rune][0]
		ui.renderer.Copy(ui.textureAtlas, &itemSrcrect, ui.GroundItemRect(i))
	}
}

func (ui *ui) GroundItemRect(i int) *sdl.Rect {
	itemSize := int32(itemSizeratio * float32(ui.winWidth))
	return &sdl.Rect{int32(ui.winWidth) - itemSize - int32(i)*itemSize, int32(ui.winHeight) - itemSize, itemSize, itemSize}
}

func (ui *ui) GetSinglePixelTex(color sdl.Color) *sdl.Texture {
	tex, err := ui.renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STATIC, 1, 1)
	if err != nil {
		panic(err)
	}
	pixels := make([]byte, 4)
	pixels[0] = color.R
	pixels[1] = color.G
	pixels[2] = color.B
	pixels[3] = color.A
	tex.Update(nil, pixels, 4)
	return tex
}

func (ui *ui) Run() {
	var newLevel *game.Level
	ui.prevMouseState = GetMouseState()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch e := event.(type) {
			case *sdl.QuitEvent:
				ui.inputChan <- &game.Input{Typ: game.QuitGame}
			case *sdl.WindowEvent:
				if e.Event == sdl.WINDOWEVENT_CLOSE {
					ui.inputChan <- &game.Input{Typ: game.CloseWindow, LevelChannel: ui.levelChan}
				}
			}
		}
		ui.currentMouseState = GetMouseState()
		select {
		case newLevel = <-ui.levelChan:
		default:
		}
		ui.Draw(newLevel)
		var input game.Input
		if ui.state == uiInventory {
			if ui.draggedItem != nil && !ui.currentMouseState.leftButton && ui.prevMouseState.leftButton {
				item := ui.CheckDroppedItem(newLevel)
				if item != nil {
					input.Typ = game.DropItem
					input.Item = ui.draggedItem
					ui.draggedItem = nil
				}
			}
			if !ui.currentMouseState.leftButton || ui.draggedItem == nil {
				ui.draggedItem = ui.CheckInventoryItems(newLevel)
			}
			ui.DrawInventory(newLevel)
		} else if ui.state == uiCodon {
			codonStartX := int32(float32(ui.winWidth) * 0.3)
			codonStartY := int32(float32(ui.winHeight) * 0.2)
			img := ui.ImgFromFileToTexture("ui/assets/codon_table.png")
			ui.renderer.Copy(img, nil, &sdl.Rect{codonStartX, codonStartY, 552, 473})
		}
		ui.renderer.Present()
		item := ui.CheckGroundItems(newLevel)
		if item != nil {
			input.Typ = game.TakeItem
			input.Item = item
		}
		if sdl.GetKeyboardFocus() == ui.window || sdl.GetMouseFocus() == ui.window {
			if ui.keyboardState[sdl.SCANCODE_UP] == 1 && ui.prevKeyboardState[sdl.SCANCODE_UP] == 0 {
				input.Typ = game.Up
			} else if ui.keyboardState[sdl.SCANCODE_DOWN] == 1 && ui.prevKeyboardState[sdl.SCANCODE_DOWN] == 0 {
				input.Typ = game.Down
			} else if ui.keyboardState[sdl.SCANCODE_LEFT] == 1 && ui.prevKeyboardState[sdl.SCANCODE_LEFT] == 0 {
				input.Typ = game.Left
			} else if ui.keyboardState[sdl.SCANCODE_RIGHT] == 1 && ui.prevKeyboardState[sdl.SCANCODE_RIGHT] == 0 {
				input.Typ = game.Right
			} else if ui.keyboardState[sdl.SCANCODE_S] == 0 && ui.prevKeyboardState[sdl.SCANCODE_S] != 0 {
				input.Typ = game.Search
			} else if ui.keyboardState[sdl.SCANCODE_T] == 0 && ui.prevKeyboardState[sdl.SCANCODE_T] != 0 {
				input.Typ = game.TakeAll
			} else if ui.keyboardState[sdl.SCANCODE_I] == 0 && ui.prevKeyboardState[sdl.SCANCODE_I] != 0 {
				if ui.state == uiMain {
					ui.state = uiInventory
				} else {
					ui.state = uiMain
				}
			} else if ui.keyboardState[sdl.SCANCODE_C] == 0 && ui.prevKeyboardState[sdl.SCANCODE_C] != 0 {
				if ui.state == uiMain {
					ui.state = uiCodon
				} else {
					ui.state = uiMain
				}
			}
			for i, v := range ui.keyboardState {
				ui.prevKeyboardState[i] = v
			}
			if input.Typ != game.None {
				ui.inputChan <- &input
			}
		}
		ui.prevMouseState = ui.currentMouseState
		sdl.Delay(10)
	}
}
