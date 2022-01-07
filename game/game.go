package game

import (
	"bufio"
	"fmt"
	"math"
	"os"

	// "time"
	"encoding/csv"
	"path/filepath"
	"strconv"
	"strings"
)

type Game struct {
	LevelChans   []chan *Level
	InputChan    chan *Input
	Levels       map[string]*Level
	CurrentLevel *Level
}

func NewGame(numWindows int) *Game {
	levelChans := make([]chan *Level, numWindows)
	for i := range levelChans {
		levelChans[i] = make(chan *Level)
	}
	inputChan := make(chan *Input)
	levels := LoadMapsFile()
	game := &Game{levelChans, inputChan, levels, nil}
	game.LoadWorldFile()
	game.CurrentLevel.LineOfSight()
	return game
}

type InputType int

const (
	None InputType = iota
	Up
	Down
	Left
	Right
	TakeAll
	TakeItem
	QuitGame
	CloseWindow
	DropItem
	Search //temporary

)

type Input struct {
	Typ          InputType
	Item         *Item
	LevelChannel chan *Level
}

type Tile struct {
	Rune        rune
	OverlayRune rune
	Visible     bool
	Seen        bool
}

const (
	Wall       rune = '#'
	Floor      rune = '.'
	ClosedDoor rune = '|'
	OpenDoor   rune = '/'
	Blank      rune = 0
	Pending    rune = -1
	UpStair    rune = 'u'
	DownStair  rune = 'd'
)

type Pos struct {
	X, Y int
}

type LevelPos struct {
	*Level
	Pos
}

type Entity struct {
	Pos
	Name string
	Rune rune
}

type Character struct {
	Entity
	ActionPoints float64
	SightRange   int
	Items        []*Item
}

type Player struct {
	Character
}

type Level struct {
	Organelle    [][]Tile
	Player       *Player
	Debug        map[Pos]bool
	Events       []string
	Descriptions []string
	EventPos     int
	DesPos       int
	Portals      map[Pos]*LevelPos
	Items        map[Pos][]*Item
}

func (level *Level) DropItem(itemTomove *Item, character *Character) {
	pos := character.Pos
	items := character.Items
	for i, item := range items {
		if item == itemTomove {
			character.Items = append(character.Items[:i], character.Items[i+1:]...)
			level.Items[pos] = append(level.Items[pos], item)
			level.AddEvent(character.Name + " dropped: " + item.Name)
			return
		}
	}
}

func (level *Level) MoveItem(movedItem *Item, character *Character) {
	pos := character.Pos
	items := level.Items[pos]
	for i, item := range items {
		if item == movedItem {
			items = append(items[:i], items[i+1:]...)
			level.Items[pos] = items
			character.Items = append(character.Items, item)
			level.AddEvent(character.Name + " picked up: " + item.Name)
			if item.Description != "" {
				level.AddDescriptionEvent(item.Description)
			}
			if item.Sequence != "" && item.Name != "Protein" {
				level.AddEvent("Current Sequence: " + item.Sequence)
			}
			if item.Name == "Intermdiate Info" {
				level.AddEvent("New Sequence: " + item.Sequence)
			} else if item.Name == "Protein" {
				level.AddEvent("Final Translated Protein Sequence: " + item.Sequence)
			}
			return
		}
	}
	panic("Tried to pick up something that you were not on top of.")
}

func (level *Level) AddEvent(event string) {
	level.Events[level.EventPos] = event
	level.EventPos++
	if level.EventPos == len(level.Events) {
		level.EventPos = 0
	}
}

func (level *Level) AddDescriptionEvent(description string) {
	level.Descriptions[level.DesPos] = description
	level.DesPos++
	if level.DesPos == len(level.Descriptions) {
		level.DesPos = 0
	}
}

func (level *Level) LineOfSight() {
	pos := level.Player.Pos
	dist := level.Player.SightRange

	for y := pos.Y - dist; y <= pos.Y+dist; y++ {
		for x := pos.X - dist; x <= pos.X+dist; x++ {
			xDelta := pos.X - x
			yDelta := pos.Y - y
			d := math.Sqrt(float64(xDelta*xDelta + yDelta*yDelta))
			if d <= float64(dist) {
				level.Bresenham(pos, Pos{x, y})
			}
		}
	}
}

func (level *Level) Bresenham(start, end Pos) {
	steep := math.Abs(float64(end.Y-start.Y)) > math.Abs(float64(end.X-start.X))
	if steep {
		start.X, start.Y = start.Y, start.X
		end.X, end.Y = end.Y, end.X
	}
	deltaY := int(math.Abs(float64(end.Y - start.Y)))
	err := 0
	y := start.Y
	ystep := 1
	if start.Y >= end.Y {
		ystep = -1
	}
	if start.X > end.X {
		deltaX := start.X - end.X
		for x := start.X; x > end.X; x-- {
			var pos Pos
			if steep {
				pos = Pos{y, x}
			} else {
				pos = Pos{x, y}
			}
			level.Organelle[pos.Y][pos.X].Visible = true
			level.Organelle[pos.Y][pos.X].Seen = true
			if !CanSeeThrough(level, pos) {
				return
			}
			err += deltaY
			if 2*err >= deltaX {
				y += ystep
				err -= deltaX
			}
		}
	} else {
		deltaX := end.X - start.X
		for x := start.X; x < end.X; x++ {
			var pos Pos
			if steep {
				pos = Pos{y, x}
			} else {
				pos = Pos{x, y}
			}
			level.Organelle[pos.Y][pos.X].Visible = true
			level.Organelle[pos.Y][pos.X].Seen = true
			if !CanSeeThrough(level, pos) {
				return
			}
			err += deltaY
			if 2*err >= deltaX {
				y += ystep
				err -= deltaX
			}
		}
	}
}

func (game *Game) LoadWorldFile() {
	file, err := os.Open("game/maps/world.txt")
	if err != nil {
		panic("Could not open world text")
	}
	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = -1
	csvReader.TrimLeadingSpace = true
	rows, err := csvReader.ReadAll()
	if err != nil {
		panic("Trouble reading the CSV file")
	}
	for rowi, row := range rows {
		if rowi == 0 {
			game.CurrentLevel = game.Levels[row[0]]
			if (game.CurrentLevel) == nil {
				fmt.Println("Could not find current level")
				panic(nil)
			}
			continue
		}
		levelWithportal := game.Levels[row[0]]
		if (levelWithportal) == nil {
			fmt.Println("Could not find portal level names")
			panic(nil)
		}
		x, _ := strconv.ParseInt(row[1], 10, 64)
		y, _ := strconv.ParseInt(row[2], 10, 64)
		pos := Pos{int(x), int(y)}
		levelToteleport := game.Levels[row[3]]
		if (levelToteleport) == nil {
			fmt.Println("Could not find teleport level names")
			panic(nil)
		}
		x, _ = strconv.ParseInt(row[4], 10, 64)
		y, _ = strconv.ParseInt(row[5], 10, 64)
		posToteleport := Pos{int(x), int(y)}
		levelWithportal.Portals[pos] = &LevelPos{levelToteleport, posToteleport}
	}

}

func LoadMapsFile() map[string]*Level {
	player := &Player{}
	player.Name = "Sriram"
	player.Rune = '@'
	player.ActionPoints = 0
	player.SightRange = 7

	levels := make(map[string]*Level)

	filenames, err := filepath.Glob("game/maps/*.map")
	if err != nil {
		panic("Could not load levels")
	}
	for _, filename := range filenames {
		extenInd := strings.LastIndex(filename, ".map")
		lastSlashindex := strings.LastIndex(filename, "\\")
		levelName := filename[lastSlashindex+1 : extenInd]
		file, err := os.Open(filename)
		if err != nil {
			panic("Could not open file")
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		levelLines := make([]string, 0)
		longestRow := 0
		index := 0
		for scanner.Scan() {
			levelLines = append(levelLines, scanner.Text())
			if len(levelLines[index]) > longestRow {
				longestRow = len(levelLines[index])
			}
			index++
		}
		level := &Level{}
		level.Debug = make(map[Pos]bool)
		level.Events = make([]string, 9)
		level.Descriptions = make([]string, 1)
		level.Player = player

		level.Organelle = make([][]Tile, len(levelLines))
		level.Portals = make(map[Pos]*LevelPos)
		level.Items = make(map[Pos][]*Item)
		for i := range level.Organelle {
			level.Organelle[i] = make([]Tile, longestRow)
		}
		for y := 0; y < len(level.Organelle); y++ {
			line := levelLines[y]

			for x, c := range line {
				pos := Pos{x, y}
				var t Tile
				t.OverlayRune = Blank
				switch c {
				case ' ', '\t', '\n', '\r':
					t.Rune = Blank
				case '#':
					t.Rune = Wall
				case '|':
					t.OverlayRune = ClosedDoor
					t.Rune = Pending
				case '/':
					t.OverlayRune = OpenDoor
					t.Rune = Pending
				case '.':
					t.Rune = Floor
				case 'u':
					t.OverlayRune = UpStair
					t.Rune = Pending
				case 'd':
					t.OverlayRune = DownStair
					t.Rune = Pending
				case 'P':
					level.Items[pos] = append(level.Items[pos], Protein(pos))
					t.Rune = Pending
				case 'E':
					level.Items[pos] = append(level.Items[pos], RNAPol(pos))
					t.Rune = Pending
				case 'I':
					level.Items[pos] = append(level.Items[pos], Instructions(pos, level))
					t.Rune = Pending
				case 'Q':
					level.Items[pos] = append(level.Items[pos], FirstDirection(pos))
					t.Rune = Pending
				case 'S':
					level.Items[pos] = append(level.Items[pos], GeneSequence(pos))
					t.Rune = Pending
				case 'T':
					level.Items[pos] = append(level.Items[pos], TranscriptionFactor(pos))
					t.Rune = Pending
				case 'N':
					level.Items[pos] = append(level.Items[pos], InterInfo(pos))
					t.Rune = Pending
				case 'R':
					level.Items[pos] = append(level.Items[pos], TranslationInfo(pos))
					t.Rune = Pending
				case 'F':
					level.Items[pos] = append(level.Items[pos], FinalNote(pos))
					t.Rune = Pending
				case '@':
					level.Player.X = x
					level.Player.Y = y
					t.Rune = Pending
				default:
					panic("Not a valid character in map")
				}
				level.Organelle[y][x] = t
			}
		}
		for y, row := range level.Organelle {
			for x, tile := range row {
				if tile.Rune == Pending {
					level.Organelle[y][x].Rune = level.BreadthFirstSearchFloor(Pos{x, y})
				}
			}
		}
		levels[levelName] = level
	}
	return levels
}

func InRange(level *Level, pos Pos) bool {
	return pos.X < len(level.Organelle[0]) && pos.Y < len(level.Organelle) && pos.X >= 0 && pos.Y >= 0
}

func CanWalk(level *Level, pos Pos) bool {
	if InRange(level, pos) {
		tile := level.Organelle[pos.Y][pos.X]
		switch tile.Rune {
		case Wall, Blank:
			return false
		}
		switch tile.OverlayRune {
		case ClosedDoor:
			return false
		}
		return true
	}
	return false
}

func CheckDoor(level *Level, pos Pos) {
	tile := level.Organelle[pos.Y][pos.X]
	if tile.OverlayRune == ClosedDoor {
		level.Organelle[pos.Y][pos.X].OverlayRune = OpenDoor
		level.LineOfSight()
	}
}

func (game *Game) Move(to Pos) {
	level := game.CurrentLevel
	player := level.Player
	levelAndpos := level.Portals[to]
	if levelAndpos != nil {
		game.CurrentLevel = levelAndpos.Level
		game.CurrentLevel.Player.Pos = levelAndpos.Pos
		game.CurrentLevel.LineOfSight()
	} else {
		player.Pos = to
		for y, row := range level.Organelle {
			for x := range row {
				level.Organelle[y][x].Visible = false
			}
		}
		level.LineOfSight()
	}
}

func (game *Game) ResolveMovement(pos Pos) {
	level := game.CurrentLevel
	if CanWalk(level, pos) {
		game.Move(pos)
	} else {
		CheckDoor(level, pos)
	}
}

func CanSeeThrough(level *Level, pos Pos) bool {
	if InRange(level, pos) {
		tile := level.Organelle[pos.Y][pos.X]
		switch tile.Rune {
		case Wall, ClosedDoor, Blank:
			return false
		default:
			return true
		}
	}
	return false
}

func (game *Game) HandleInput(input *Input) {
	level := game.CurrentLevel
	p := level.Player
	switch input.Typ {
	case Up:
		newPos := Pos{p.X, p.Y - 1}
		game.ResolveMovement(newPos)
	case Down:
		newPos := Pos{p.X, p.Y + 1}
		game.ResolveMovement(newPos)
	case Left:
		newPos := Pos{p.X - 1, p.Y}
		game.ResolveMovement(newPos)
	case Right:
		newPos := Pos{p.X + 1, p.Y}
		game.ResolveMovement(newPos)
	case Search:
		level.Astar(level.Player.Pos, Pos{3, 3})
	case TakeAll:
		for _, item := range level.Items[p.Pos] {
			level.MoveItem(item, &level.Player.Character)
		}
	case TakeItem:
		level.MoveItem(input.Item, &level.Player.Character)
	case DropItem:
		level.DropItem(input.Item, &level.Player.Character)
	case CloseWindow:
		close(input.LevelChannel)
		chanIndex := 0
		for i, c := range game.LevelChans {
			if c == input.LevelChannel {
				chanIndex = i
				break
			}
		}
		game.LevelChans = append(game.LevelChans[:chanIndex], game.LevelChans[chanIndex+1:]...)
	}
}

func GetNeighbors(level *Level, pos Pos) []Pos {
	neighbors := make([]Pos, 0, 4)
	left := Pos{pos.X - 1, pos.Y}
	right := Pos{pos.X + 1, pos.Y}
	up := Pos{pos.X, pos.Y - 1}
	down := Pos{pos.X, pos.Y + 1}

	if CanWalk(level, left) {
		neighbors = append(neighbors, left)
	}
	if CanWalk(level, right) {
		neighbors = append(neighbors, right)
	}
	if CanWalk(level, up) {
		neighbors = append(neighbors, up)
	}
	if CanWalk(level, down) {
		neighbors = append(neighbors, down)
	}
	return neighbors
}

func (level *Level) BreadthFirstSearchFloor(start Pos) rune {
	frontier := make([]Pos, 0, 8)
	frontier = append(frontier, start)
	visited := make(map[Pos]bool)
	visited[start] = true

	for len(frontier) > 0 {
		current := frontier[0]
		currentTile := level.Organelle[current.Y][current.X]
		switch currentTile.Rune {
		case Floor:
			return Floor
		default:
		}
		frontier = frontier[1:]
		for _, next := range GetNeighbors(level, current) {
			if !visited[next] {
				frontier = append(frontier, next)
				visited[next] = true
			}
		}
	}
	return Floor
}

func (level *Level) Astar(start Pos, goal Pos) []Pos {
	frontier := make(pQueue, 0, 8)
	frontier = frontier.Push(start, 1)
	cameFrom := make(map[Pos]Pos)
	cameFrom[start] = start
	costSoFar := make(map[Pos]int)
	costSoFar[start] = 0

	var current Pos
	for len(frontier) > 0 {
		frontier, current = frontier.Pop()
		if current == goal {
			path := make([]Pos, 0)
			p := current
			for p != start {
				path = append(path, p)
				p = cameFrom[p]
			}
			path = append(path, p)
			fmt.Println("Done!")
			for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
				path[i], path[j] = path[j], path[i]
			}

			return path
		}

		for _, next := range GetNeighbors(level, current) {
			newCost := costSoFar[current] + 1
			_, exists := costSoFar[next]
			if !exists || newCost < costSoFar[next] {
				costSoFar[next] = newCost
				xDist := int(math.Abs(float64(goal.X - next.X)))
				yDist := int(math.Abs(float64(goal.Y - next.Y)))
				priority := newCost + xDist + yDist
				frontier = frontier.Push(next, priority)
				cameFrom[next] = current
			}
		}
	}
	return nil
}

func (game *Game) Run() {
	fmt.Println("Running Game")
	for _, lchan := range game.LevelChans {
		lchan <- game.CurrentLevel
	}
	for input := range game.InputChan {
		if input.Typ == QuitGame {
			return
		}
		game.HandleInput(input)

		if len(game.LevelChans) == 0 {
			return
		}
		for _, lchan := range game.LevelChans {
			lchan <- game.CurrentLevel
		}
	}

}
