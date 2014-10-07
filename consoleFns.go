package main

import (
	"fmt"
	"strings"

	"github.com/nsf/termbox-go"
)

func getConsoleMsg(tx *termboxState) string {
	var s string
	if tx.consoleFn != nil {
		s = tx.consoleFn(tx)
	}
	return s
}

func emptyConsoleFn(*termboxState) string {
	return ""
}

func episodeConsoleFn(tx *termboxState) string {
	return fmt.Sprintf("%d Seasons. %d Episodes", len(tx.results.Series[tx.index].Seasons), tx.totalEpisodes)
}

func printConsoleString(tx *termboxState) {
	width, height := termbox.Size()
	s := getConsoleMsg(tx)
	pad := width - len(s)
	s += strings.Repeat(" ", pad)
	for x, r := range s {
		termbox.SetCell(x, height-1, r, termbox.ColorWhite, termbox.ColorBlue)
	}
}

func cursorBlink(tx *termboxState) {
	_, height := termbox.Size()
	x := len(getConsoleMsg(tx))
	color := termbox.ColorBlue
	if tx.blink {
		color = termbox.ColorWhite
	}
	tx.blink = !tx.blink
	termbox.SetCell(x, height-1, rune(' '), color, color)
}
