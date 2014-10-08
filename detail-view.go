package main

// Episode UI state management

import "github.com/nsf/termbox-go"

func detailEpisodeCursorUp(tx *termboxState) {
	if tx.episodeIndex > 0 {
		tx.episodeIndex--
		updateScreen(tx, drawDetails)
	}
}

func detailEpisodeCursorDown(tx *termboxState) {
	if tx.episodeIndex < tx.totalEpisodes-1 {
		tx.episodeIndex++
		updateScreen(tx, drawDetails)
	}
}

// When inside the Episode UI this is the termbox event handler
func DetailEventHandler(tx *termboxState) stateFn {
	switch tx.ev.Type {
	case termbox.EventKey:
		switch tx.ev.Key {
		case termbox.KeyEsc, termbox.KeyArrowLeft:
			updateScreen(tx, drawEpisode)
			return EpisodeEventHandler
		}
		switch tx.ev.Ch {
		case keyj:
			detailEpisodeCursorDown(tx)
		case keyk:
			detailEpisodeCursorUp(tx)
		case keyh:
			updateScreen(tx, drawEpisode)
			return EpisodeEventHandler
		}
	case termbox.EventResize:
		updateScreen(tx, drawDetails)
	}
	return DetailEventHandler
}
