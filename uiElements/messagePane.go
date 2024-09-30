package uiElements

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func ShowMessage(msg string) {
	termWidth, termHeight := ui.TerminalDimensions()
	msgDialog := widgets.NewParagraph()
	msgDialog.Text = msg
	msgDialog.SetRect(termWidth/2-20, termHeight/2-3, termWidth/2+20, termHeight/2+3)
	msgDialog.Border = true
	msgDialog.TextStyle.Fg = ui.ColorRed

	ui.Render(msgDialog)

	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<Escape>", "<Enter>", "<C-c>":
			ui.Clear()
			return
		}
        ui.Render(msgDialog)
	}
}
