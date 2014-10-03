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
