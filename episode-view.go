package main

// Episode UI state management

import (
	"math"

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

func episodeCursorPgup(tx *termboxState) {
	_, height := getViewPortSize(tx)
	if tx.totalEpisodes > height {
		tx.episodeIndex = int(math.Max(float64(tx.episodeIndex-height), float64(0)))
	}
	updateScreen(tx, drawEpisode)
}

func episodeCursorPgdn(tx *termboxState) {
	_, height := getViewPortSize(tx)
	if tx.totalEpisodes > height {
		tx.episodeIndex = int(math.Min(float64(tx.episodeIndex+height), float64(tx.totalEpisodes-1)))
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
		case termbox.KeyPgup, termbox.KeyCtrlB:
			episodeCursorPgup(tx)
		case termbox.KeyPgdn, termbox.KeyCtrlF:
			episodeCursorPgdn(tx)

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

func transitionToSeriesState(tx *termboxState) stateFn {
	updateScreen(tx, drawSeries)
	return SeriesEventHandler
}
