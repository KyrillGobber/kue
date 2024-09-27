package menu

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func GetSceneMenu(scenes []string, coords Coords) *widgets.List {

	sceneList := widgets.NewList()
	sceneList.Rows = scenes
	sceneList.TextStyle = ui.NewStyle(ui.ColorWhite)
    sceneList.SelectedRowStyle.Fg = ui.ColorYellow
	sceneList.WrapText = false
	sceneList.SetRect(coords.X1, coords.Y1, coords.X2, coords.Y2)

    return sceneList
}
