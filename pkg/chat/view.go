package chat

import "github.com/jroimartin/gocui"

// TerminalUI creates the chat UI in our terminal
func TerminalUI(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	g.Cursor = true

	if messages, err := g.SetView("messages", 0, 0, maxX-20, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		messages.Title = " Received Messages: "
		messages.Autoscroll = true
		messages.Wrap = true
	}

	if input, err := g.SetView("input", 0, maxY-10, maxX-20, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		input.Title = " Send Message (begin with 'all' or list of IDs like '1,2,3') "
		input.Autoscroll = false
		input.Wrap = true
		input.Editable = true
	}

	if clientIDList, err := g.SetView("clientIDList", maxX-20, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		clientIDList.Title = " Client ID List: "
		clientIDList.Autoscroll = false
		clientIDList.Wrap = true
	}

	if name, err := g.SetView("name", maxX/2-10, maxY/2-1, maxX/2+10, maxY/2+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		g.SetCurrentView("name")
		name.Title = " name: "
		name.Autoscroll = false
		name.Wrap = true
		name.Editable = true
	}

	return nil
}
