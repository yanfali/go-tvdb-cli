package main

// Add a simple state machine function and context variable to simplify
// UI state management

import (
	"github.com/nsf/termbox-go"
	"github.com/yanfali/go-tvdb"
)

// Return the next state you want to be in
// This can be the current state you are already in
type stateFn func(tx *termboxState) stateFn

// Track the state machine context
type termboxState struct {
	ev         termbox.Event    // last recorded event
	results    *tvdb.SeriesList // current results from thetvdb.com
	index      int              // current index into the results
	lastInput  rune             // last character typed
	consoleMsg string           // current console message
}

// When inside the Episode UI this is the termbox event handler
func EpisodeEventHandler(tx *termboxState) stateFn {
	switch tx.ev.Type {
	case termbox.EventKey:
		switch tx.ev.Key {
		case termbox.KeyEsc:
			tx.consoleMsg = ""
			updateScreen(tx, drawAll)
			return SeriesEventHandler
		}
	case termbox.EventResize:
		updateScreen(tx, drawEpisode)
	}
	return EpisodeEventHandler
}

// When inside the Series UI this is the termbox event handler
func SeriesEventHandler(tx *termboxState) stateFn {
	switch tx.ev.Type {
	case termbox.EventKey:
		switch tx.ev.Key {
		case termbox.KeyEsc:
			return nil // returning a nil state means state machine has reached stop
		}
		switch tx.ev.Ch {
		case keyZero, keyOne, keyTwo, keyThree, keyFour, keyFive, keySix, keySeven, keyEight, keyNine:
			tx.lastInput = tx.ev.Ch
			tx.index = int(tx.ev.Ch - '0')
			if len(tx.results.Series[tx.index].Seasons) == 0 {
				if err := tx.results.Series[tx.index].GetDetail(config); err != nil {
					updateScreen(tx, drawAll)
					return SeriesEventHandler
				}
			}
			updateScreen(tx, drawEpisode)
			return EpisodeEventHandler
		}
	case termbox.EventResize:
		updateScreen(tx, drawAll)
	}
	return SeriesEventHandler
}
