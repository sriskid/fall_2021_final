package ui

import (
	"final_project/game"

	"github.com/veandco/go-sdl2/sdl"
)

func (ui *ui) DrawInventory(level *game.Level) {
	playerSrcRect := ui.textureIndex[level.Player.Rune][0]
	invRect := ui.InventoryRect()
	ui.renderer.Copy(ui.textureAtlas, &playerSrcRect, &sdl.Rect{invRect.X + invRect.X/4, invRect.Y, invRect.W / 2, invRect.H / 2})
	ui.renderer.Copy(ui.inventoryBackground, nil, invRect)
	for i, item := range level.Player.Items {
		itemSize := int32(itemSizeratio * float32(ui.winWidth))
		itemSrcrect := ui.textureIndex[item.Rune][0]
		if item == ui.draggedItem {
			ui.renderer.Copy(ui.textureAtlas, &itemSrcrect, &sdl.Rect{int32(ui.currentMouseState.pos.X), int32(ui.currentMouseState.pos.Y), itemSize, itemSize})
		} else {
			ui.renderer.Copy(ui.textureAtlas, &itemSrcrect, ui.InventoryItemRect(i))
		}

	}
}

func (ui *ui) InventoryRect() *sdl.Rect {
	invWidth := int32(float32(ui.winWidth) * 0.40)
	invHeight := int32(float32(ui.winHeight) * 0.75)
	offSetx := (int32(ui.winWidth) - invWidth) / 2
	offSety := (int32(ui.winHeight) - invHeight) / 2
	return &sdl.Rect{offSetx, offSety, invWidth, invHeight}
}

func (ui *ui) InventoryItemRect(i int) *sdl.Rect {
	invRect := ui.InventoryRect()
	itemSize := int32(itemSizeratio * float32(ui.winWidth))
	return &sdl.Rect{invRect.X + int32(i)*itemSize, invRect.Y + invRect.H - itemSize, itemSize, itemSize}
}

func (ui *ui) CheckDroppedItem(level *game.Level) *game.Item {
	invRect := ui.InventoryRect()
	mousePos := ui.currentMouseState.pos
	if invRect.HasIntersection(&sdl.Rect{int32(mousePos.X), int32(mousePos.Y), 1, 1}) {
		return nil
	}
	return ui.draggedItem
}

func (ui *ui) CheckInventoryItems(level *game.Level) *game.Item {
	if ui.currentMouseState.leftButton {
		mousePos := ui.currentMouseState.pos
		for i, item := range level.Player.Items {
			itemRect := ui.InventoryItemRect(i)
			if itemRect.HasIntersection(&sdl.Rect{int32(mousePos.X), int32(mousePos.Y), 1, 1}) {
				return item
			}
		}
	}
	return nil
}

func (ui *ui) CheckGroundItems(level *game.Level) *game.Item {
	if !ui.currentMouseState.leftButton && ui.prevMouseState.leftButton {
		mousePos := ui.currentMouseState.pos
		items := level.Items[level.Player.Pos]
		for i, item := range items {
			itemRect := ui.GroundItemRect(i)
			if itemRect.HasIntersection(&sdl.Rect{int32(mousePos.X), int32(mousePos.Y), 1, 1}) {
				return item
			}
		}
	}
	return nil
}
