package menu

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func GetItemMenu(rooms []string, coords Coords) *widgets.List {

	itemList := widgets.NewList()
	itemList.Rows = rooms
	itemList.TextStyle = ui.NewStyle(ui.ColorWhite)
    itemList.SelectedRowStyle.Fg = ui.ColorYellow
	itemList.WrapText = false
	itemList.SetRect(coords.X1, coords.Y1, coords.X2, coords.Y2)

    return itemList
}
