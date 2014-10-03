package main

// Add a simple state machine function and context variable to simplify
// UI state management

import (
	"errors"
	"fmt"

	"github.com/nsf/termbox-go"
	"github.com/yanfali/go-tvdb"
)

// Return the next state you want to be in
// This can be the current state you are already in
type stateFn func(tx *termboxState) stateFn

// Track the state machine context
type termboxState struct {
	ev            termbox.Event    // last recorded event
	results       *tvdb.SeriesList // current results from thetvdb.com
	index         int              // current index into the results
	lastInput     rune             // last character typed
	consoleMsg    string           // current console message
	seriesIndex   int
	episodeIndex  int
	totalEpisodes int
	allEpisodes   []tvdb.Episode
}

func episodeCursorUp(tx *termboxState) {
	if tx.episodeIndex > 0 {
		tx.episodeIndex--
	}
	updateScreen(tx, drawEpisode)
}

func episodeCursorDown(tx *termboxState) {
	if tx.episodeIndex < tx.totalEpisodes-1 {
		tx.episodeIndex++
	}
	updateScreen(tx, drawEpisode)
}

// When inside the Episode UI this is the termbox event handler
func EpisodeEventHandler(tx *termboxState) stateFn {
	switch tx.ev.Type {
	case termbox.EventKey:
		switch tx.ev.Key {
		case termbox.KeyEsc, termbox.KeyArrowLeft:
			return eventSeries(tx)
		case termbox.KeyArrowUp:
			episodeCursorUp(tx)
		case termbox.KeyArrowDown:
			episodeCursorDown(tx)
		}
		switch tx.ev.Ch {
		case keyh:
			return eventSeries(tx)
		case keyj:
			episodeCursorDown(tx)
		case keyk:
			episodeCursorUp(tx)
		}
	case termbox.EventResize:
		updateScreen(tx, drawEpisode)
	}
	return EpisodeEventHandler
}

func fetchSeasons(tx *termboxState) error {
	if len(tx.results.Series[tx.index].Seasons) == 0 {
		if err := tx.results.Series[tx.index].GetDetail(config); err != nil {
			updateScreen(tx, drawAll)
			return errors.New("error requesting series details")
		}
	}
	tx.allEpisodes = []tvdb.Episode{}
	for _, episodes := range tx.results.Series[tx.index].Seasons {
		tx.allEpisodes = append(tx.allEpisodes, episodes...)
	}
	By(seasonAndEpisode).Sort(tx.allEpisodes)
	tx.totalEpisodes = len(tx.allEpisodes)
	tx.episodeIndex = 0
	return nil
}

func seriesCursorUp(tx *termboxState) {
	if tx.seriesIndex > 0 {
		tx.seriesIndex--
	}
	updateScreen(tx, drawAll)
}

func seriesCursorDown(tx *termboxState) {
	if tx.seriesIndex < len(tx.results.Series)-1 {
		tx.seriesIndex++
	}
	updateScreen(tx, drawAll)
}

func eventEpisode(tx *termboxState) stateFn {
	tx.index = tx.seriesIndex
	if err := fetchSeasons(tx); err != nil {
		tx.consoleMsg = fmt.Sprintf("%v", err)
		updateScreen(tx, drawAll)
		return SeriesEventHandler
	}
	updateScreen(tx, drawEpisode)
	return EpisodeEventHandler
}

func eventSeries(tx *termboxState) stateFn {
	updateScreen(tx, drawAll)
	return SeriesEventHandler
}

// When inside the Series UI this is the termbox event handler
func SeriesEventHandler(tx *termboxState) stateFn {
	switch tx.ev.Type {
	case termbox.EventKey:
		switch tx.ev.Key {
		case termbox.KeyEsc:
			return nil // returning a nil state means state machine has reached stop
		case termbox.KeyArrowUp:
			seriesCursorUp(tx)
		case termbox.KeyArrowDown:
			seriesCursorDown(tx)
		case termbox.KeyEnter, termbox.KeyArrowRight:
			return eventEpisode(tx)
		}
		switch tx.ev.Ch {
		case keyj:
			seriesCursorDown(tx)
		case keyk:
			seriesCursorUp(tx)
		case keyl:
			return eventEpisode(tx)
		}
	case termbox.EventResize:
		updateScreen(tx, drawAll)
	}
	return SeriesEventHandler
}
