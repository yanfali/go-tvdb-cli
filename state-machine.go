package main

// Add a simple state machine function and context variable to simplify
// UI state management

import (
	"errors"

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
	seriesIndex   int              // current index into series
	episodeIndex  int              // current index into episodes
	totalEpisodes int              // total episodes in current series view
	allEpisodes   []tvdb.Episode   // all episodes for current series
	stateFnStack  []stateFn        // stack of functions
	searchText    string
	blink         bool
}

// push a state function onto a stack
// use this to decouple a transition from a state
func (my *termboxState) Push(fn stateFn) {
	my.stateFnStack = append(my.stateFnStack, fn)
}

// pop a state function from the stack
// the receiver of the popped function must
// invoke it and return the result to get
// the correct behavior.
func (my *termboxState) Pop() (stateFn, error) {
	fnCount := len(my.stateFnStack)
	if fnCount == 0 {
		return nil, errors.New("Stack is empty")
	}
	fn := my.stateFnStack[fnCount-1]
	my.stateFnStack = my.stateFnStack[0 : fnCount-1]
	return fn, nil
}
