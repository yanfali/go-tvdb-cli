package main

// Episode UI state management

import (
	"fmt"

	"github.com/nsf/termbox-go"
)

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
			return transitionToSeriesState(tx)
		case termbox.KeyArrowUp:
			episodeCursorUp(tx)
		case termbox.KeyArrowDown:
			episodeCursorDown(tx)
		}
		switch tx.ev.Ch {
		case keyh:
			return transitionToSeriesState(tx)
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

func transitionToEpisodeState(tx *termboxState) stateFn {
	tx.index = tx.seriesIndex
	if err := fetchSeasons(tx); err != nil {
		tx.consoleMsg = fmt.Sprintf("%v", err)
		updateScreen(tx, drawSeries)
		return SeriesEventHandler
	}
	updateScreen(tx, drawEpisode)
	return EpisodeEventHandler
}
