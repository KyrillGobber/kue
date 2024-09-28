package uiElements

import "github.com/gizak/termui/v3/widgets"

func GetFooter() *widgets.Paragraph {
	footer := widgets.NewParagraph()
	footer.Text = "hjkl: navigation | enter: scene select | ESC: back | t: quicktoggle on/off | q: quit"
	footer.SetRect(0, 31, 100, 34)
	footer.Border = true
	return footer
}
