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
	"strings"

	"github.com/nsf/termbox-go"
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
	title = "go-tvdb-cli edition"
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
	sWidth := len(curs.text)
	pad := width - sWidth
	s := curs.text + strings.Repeat(" ", pad)
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

// render episode view
func drawEpisode(tx *termboxState) {
	series := tx.results.Series[tx.index]
	var viewOffset = 0
	var width, viewPortHeight = getViewPortSize(tx)
	center := width / 2
	printTermboxString(center-len(title)/2, 0, title)
	if tx.episodeIndex+1 >= viewPortHeight {
		// if episode index has moved beyond view port change view offset and slice allEpisodes differently
		viewOffset = tx.episodeIndex + 1 - viewPortHeight
	}
	printTermboxString(1, 2, printSeries(&series))
	tx.consoleMsg = fmt.Sprintf("%d Seasons. %d Episodes", len(series.Seasons), tx.totalEpisodes)

	for row, episode := range tx.allEpisodes[viewOffset : viewOffset+viewPortHeight] {
		episode.EpisodeName = ellipsisString(episode.EpisodeName, 70)
		printTermboxString(1, 4+row, fmt.Sprintf("%2v %2v %-70s", episode.SeasonNumber, episode.EpisodeNumber, episode.EpisodeName))
	}
	episode := tx.allEpisodes[tx.episodeIndex]
	s := fmt.Sprintf("%2v %2v %-70s", episode.SeasonNumber, episode.EpisodeNumber, episode.EpisodeName)
	curs := cursorMeta{text: s, xOrigin: 1, yOrigin: 4, yOffset: tx.episodeIndex - viewOffset}
	//tx.consoleMsg = fmt.Sprintf("vpo: %d vph: %d vpei: %d %v", viewOffset, viewPortHeight, tx.episodeIndex, curs)
	drawCursor(curs)
}

// render series view
func drawSeries(tx *termboxState) {
	width, _ := termbox.Size()
	center := width / 2
	printTermboxString(center-len(title)/2, 1, title)
	printTermboxString(1, 3, fmt.Sprintf("%-5s %-40s %-10s %-15s %s", "Index", "Title", "Language", "First Aired", "Genre"))
	for row, series := range tx.results.Series {
		printTermboxString(1, 4+row, fmt.Sprintf("%-5d %s", row, printSeries(&series)))
	}

	s := fmt.Sprintf("%-5d %s", tx.seriesIndex, printSeries(&tx.results.Series[tx.seriesIndex]))
	curs := cursorMeta{text: s, xOrigin: 1, yOrigin: 5, yOffset: tx.seriesIndex}
	drawCursor(curs)
}

func printTermboxString(x, y int, s string) {
	for _, r := range s {
		termbox.SetCell(x, y, r, termbox.ColorWhite, termbox.ColorDefault)
		x += 1
	}
}

func printConsoleString(tx *termboxState) {
	width, height := termbox.Size()
	sWidth := len(tx.consoleMsg)
	pad := width - sWidth
	s := tx.consoleMsg + strings.Repeat(" ", pad)
	for x, r := range s {
		termbox.SetCell(x, height-1, r, termbox.ColorWhite, termbox.ColorBlue)
	}
}

func cursorBlink(tx *termboxState) {
	_, height := termbox.Size()
	x := len(tx.consoleMsg)
	color := termbox.ColorBlue
	if tx.blink {
		color = termbox.ColorWhite
	}
	tx.blink = !tx.blink
	termbox.SetCell(x, height-1, rune(' '), color, color)
}
