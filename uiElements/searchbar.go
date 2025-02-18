package uiElements

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func GetSearchBar(filter string) *widgets.Paragraph {
	searchbar := widgets.NewParagraph()
	searchbar.Title = "Filter:"
	searchbar.Text = filter
	searchbar.SetRect(50, 25, 100, 28)
	searchbar.Border = true
	searchbar.TitleStyle.Fg = ui.ColorGreen
	return searchbar
}
