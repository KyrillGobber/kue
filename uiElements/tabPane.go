package uiElements

import "github.com/gizak/termui/v3/widgets"

func GetTabs() *widgets.TabPane {
	tabpane := widgets.NewTabPane("Rooms", "Zones")
	tabpane.SetRect(5, 3, 100, 6)
	tabpane.Border = true
	return tabpane
}
