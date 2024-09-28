package uiElements

import (
	"github.com/gizak/termui/v3/widgets"
	ui "github.com/gizak/termui/v3"
)

func GetHeader() *widgets.Paragraph {
	header := widgets.NewParagraph()
	header.Text = "Kue, your CLI Hue controller"
	header.SetRect(0, 0, 100, 3)
	header.Border = true
	header.TextStyle.Fg = ui.ColorGreen
	return header
}

