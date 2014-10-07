package main

// Episode UI state management

import (
	"fmt"
	"log"

	"github.com/nsf/termbox-go"
)

func returnToCaller(tx *termboxState) stateFn {
	var (
		fn  stateFn
		err error
	)
	if fn, err = tx.Pop(); err != nil {
		log.Fatalln(err)
	}
	return fn(tx)
}

// When inside the Episode UI this is the termbox event handler
func HelpEventHandler(tx *termboxState) stateFn {
	switch tx.ev.Type {
	case termbox.EventKey:
		switch tx.ev.Key {
		case termbox.KeyEsc, termbox.KeyArrowLeft:
			return returnToCaller(tx)
		}
		switch tx.ev.Ch {
		case keyh:
			return returnToCaller(tx)
		}
	case termbox.EventResize:
		updateScreen(tx, drawHelp)
	}
	return HelpEventHandler
}

type help struct {
	keys string
	text string
}

var helpData = []help{
	{"?", "Help"},
	{"Left Arrow, h, Esc", "Back to Previous Level or Exit"},
	{"Up Arrow, k", "Move cursor up"},
	{"Down Arrow, j", "Move cursor down"},
	{"Right Arrow, l", "Down to next level"},
}

func drawHelp(tx *termboxState) {

	for row, help := range helpData {
		printTermboxString(0, row+3, fmt.Sprintf("%-20s %s", help.keys, help.text))
	}
	tx.consoleFn = emptyConsoleFn
}
