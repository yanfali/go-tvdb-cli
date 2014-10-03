package main

// Series View UI state management

import (
	"errors"

	"github.com/nsf/termbox-go"
	"github.com/yanfali/go-tvdb"
)

func fetchSeasons(tx *termboxState) error {
	if len(tx.results.Series[tx.index].Seasons) == 0 {
		if err := tx.results.Series[tx.index].GetDetail(config); err != nil {
			updateScreen(tx, drawSeries)
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
	updateScreen(tx, drawSeries)
}

func seriesCursorDown(tx *termboxState) {
	if tx.seriesIndex < len(tx.results.Series)-1 {
		tx.seriesIndex++
	}
	updateScreen(tx, drawSeries)
}

func transitionToSeriesState(tx *termboxState) stateFn {
	updateScreen(tx, drawSeries)
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
			return transitionToEpisodeState(tx)
		}
		switch tx.ev.Ch {
		case keyj:
			seriesCursorDown(tx)
		case keyk:
			seriesCursorUp(tx)
		case keyl:
			return transitionToEpisodeState(tx)
		}
	case termbox.EventResize:
		updateScreen(tx, drawSeries)
	}
	return SeriesEventHandler
}
