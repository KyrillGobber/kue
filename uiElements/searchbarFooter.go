package uiElements

import (
	"github.com/gizak/termui/v3/widgets"
)

func GetSearchBarFooter() *widgets.Paragraph {
	searchbarFooter := widgets.NewParagraph()
	searchbarFooter.SetRect(50, 28, 100, 29)
	searchbarFooter.Border = false
    searchbarFooter.Text = "Esc / Enter: apply filter | Ctl+l: clear filter"
	return searchbarFooter
}
