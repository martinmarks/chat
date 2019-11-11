package main

import (
	"log"

	"github.com/jroimartin/gocui"
	"github.com/martinmarks/chat/pkg/chat"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Fatal(err)
	}
	defer g.Close()

	g.SetManagerFunc(chat.TerminalUI)
	g.SetKeybinding("name", gocui.KeyEnter, gocui.ModNone, chat.RunApplication)
	g.SetKeybinding("input", gocui.KeyEnter, gocui.ModNone, chat.SendMessage)
	g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, chat.Quit)
	g.MainLoop()
}
