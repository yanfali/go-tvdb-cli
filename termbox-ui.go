package main

// Use termbox-go to add a simple UI
// two views - main results and episodes
// TODO
// add a way to scroll through episodes
// add a way to look at more than 10 results.
// add a way to look at the details for a specific episode
import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/nsf/termbox-go"
	"github.com/yanfali/go-tvdb"
)

var (
	keyZero     = rune('0')
	keyOne      = rune('1')
	keyTwo      = rune('2')
	keyThree    = rune('3')
	keyFour     = rune('4')
	keyFive     = rune('5')
	keySix      = rune('6')
	keySeven    = rune('7')
	keyEight    = rune('8')
	keyNine     = rune('9')
	keyh        = rune('h')
	keyj        = rune('j')
	keyk        = rune('k')
	keyl        = rune('l')
	keyQuestion = rune('?')
	keySearch   = rune('/')
)

const (
	MAIN_TITLE   = "go-tvdb-cli edition"
	HELP_TITLE   = "go-tvdb-cli Help"
	DETAIL_TITLE = "go-tvdb-cli Details"
)

func drawClear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
}

func drawFlush() {
	termbox.Flush()
}

// decorator to wrap certain termbox boilerplate
func updateScreen(tx *termboxState, fn func(tx *termboxState)) {
	drawClear()
	fn(tx)
	printConsoleString(tx)
	drawFlush()
}

// track cursor position
type cursorMeta struct {
	text    string
	xOrigin int
	yOrigin int
	yOffset int
}

// render cursor
func drawCursor(curs cursorMeta) {
	width, _ := termbox.Size()
	s := fmt.Sprintf("%-"+strconv.Itoa(width)+"s", curs.text)
	for x, r := range s {
		termbox.SetCell(1+x, 4+curs.yOffset, r, termbox.ColorBlack, termbox.ColorYellow)
	}
}

func getViewPortSize(tx *termboxState) (width int, height int) {
	width, height = termbox.Size()
	if height < 20 {
		// render window no smaller than 20 rows
		height = 20
	}
	return width, int(math.Min(float64(len(tx.allEpisodes)), float64(height-6))) // make sure we can't overflow slice
}

func episodeMatch(episode *tvdb.Episode, searchText string) bool {
	return searchText != "" && strings.Contains(strings.ToLower(episode.EpisodeName), strings.ToLower(searchText))
}

// render episode view
func drawEpisode(tx *termboxState) {
	series := tx.results.Series[tx.index]
	var viewOffset = 0
	var width, viewPortHeight = getViewPortSize(tx)
	center := width / 2
	printTermboxString(center-len(MAIN_TITLE)/2, 0, MAIN_TITLE)
	if tx.episodeIndex+1 >= viewPortHeight {
		// if episode index has moved beyond view port change view offset and slice allEpisodes differently
		viewOffset = tx.episodeIndex + 1 - viewPortHeight
	}
	printTermboxString(1, 2, printSeries(&series))

	matchCount := 0
	if tx.searchText != "" {
		// count up searches across entire collection
		for _, episode := range tx.allEpisodes {
			if episodeMatch(&episode, tx.searchText) {
				matchCount++
			}
		}
		if matchCount > 0 {
			printTermboxString(0, 0, fmt.Sprintf("%d matches out of %d", matchCount, len(tx.allEpisodes)))
		}
	}

	for row, episode := range tx.allEpisodes[viewOffset : viewOffset+viewPortHeight] {
		episode.EpisodeName = ellipsisString(episode.EpisodeName, 70)
		fg := termbox.ColorWhite
		bg := termbox.ColorDefault
		if matchCount > 0 && episodeMatch(&episode, tx.searchText) {
			// highlight visible searches
			fg |= termbox.AttrReverse
		}
		printTermboxStringAttr(1, 4+row, fmt.Sprintf("%2v %2v %-70s", episode.SeasonNumber, episode.EpisodeNumber, episode.EpisodeName), fg, bg)
	}
	episode := tx.allEpisodes[tx.episodeIndex]
	s := fmt.Sprintf("%2v %2v %-70s", episode.SeasonNumber, episode.EpisodeNumber, episode.EpisodeName)
	curs := cursorMeta{text: s, xOrigin: 1, yOrigin: 4, yOffset: tx.episodeIndex - viewOffset}
	drawCursor(curs)
}

func printTitleString(title string) {
	width, _ := termbox.Size()
	center := width / 2
	printTermboxStringAttr(center-len(title)/2, 1, title, termbox.ColorWhite|termbox.AttrUnderline, termbox.ColorDefault)
}

// render series view
func drawSeries(tx *termboxState) {
	printTitleString(MAIN_TITLE)
	printTermboxString(1, 3, fmt.Sprintf("%-5s %-40s %-10s %-15s %s", "Index", "Title", "Language", "First Aired", "Genre"))
	for row, series := range tx.results.Series {
		printTermboxString(1, 4+row, fmt.Sprintf("%-5d %s", row, printSeries(&series)))
	}

	s := fmt.Sprintf("%-5d %s", tx.seriesIndex, printSeries(&tx.results.Series[tx.seriesIndex]))
	curs := cursorMeta{text: s, xOrigin: 1, yOrigin: 5, yOffset: tx.seriesIndex}
	drawCursor(curs)
}

// Return a function that returns the supplied value once per invocatin
// repeated calls to the function return an empty string
// inspired by underscore.js
func once(value string) func() string {
	notDone := true
	return func() string {
		if notDone {
			notDone = false
			return value
		}
		return ""
	}
}

func printDetailAttributePair(key, value string, row int) {
	width, _ := termbox.Size()
	values := []string{}
	width -= 30
	if len(value) > width {
		// format longer strings into chunks of width length.
		// not smart about content and doesn't try to look for a break
		valLen := float64(len(value))
		for start, end := 0, width; start < int(valLen); start += width {
			values = append(values, value[start:end])
			end = int(math.Min(float64(end+width), valLen))
		}
		keyFn := once(key)
		for _, val := range values {
			printTermboxStringAttr(5, row, fmt.Sprintf("%-20s", keyFn()), termbox.ColorWhite, termbox.ColorBlue)
			printTermboxStringAttr(25, row, fmt.Sprintf("%-"+strconv.Itoa(width)+"s", val), termbox.ColorBlue, termbox.ColorWhite)
			row++
		}
	} else {
		printTermboxStringAttr(5, row, fmt.Sprintf("%-20s", key), termbox.ColorWhite, termbox.ColorBlue)
		printTermboxStringAttr(25, row, fmt.Sprintf("%-"+strconv.Itoa(width)+"s", value), termbox.ColorBlue, termbox.ColorWhite)
	}
}

type pairs struct {
	Key   string
	Value string
}

func drawDetails(tx *termboxState) {
	printTitleString(DETAIL_TITLE)
	rowStart := 4
	episode := tx.allEpisodes[tx.episodeIndex]
	pairs := []pairs{
		{"Id", strconv.FormatInt(int64(episode.Id), 10)},
		{"Name", episode.EpisodeName},
		{"First Aired", episode.FirstAired},
		{"Season", strconv.FormatInt(int64(episode.SeasonNumber), 10)},
		{"Episode", strconv.FormatInt(int64(episode.EpisodeNumber), 10)},
		{"Overview", episode.Overview},
	}
	for index, pair := range pairs {
		printDetailAttributePair(pair.Key, pair.Value, rowStart+index)
	}
}

func printTermboxStringAttr(x, y int, s string, fg termbox.Attribute, bg termbox.Attribute) {
	for _, r := range s {
		termbox.SetCell(x, y, r, fg, bg)
		x += 1
	}
}

func printTermboxString(x, y int, s string) {
	printTermboxStringAttr(x, y, s, termbox.ColorWhite, termbox.ColorDefault)
}
