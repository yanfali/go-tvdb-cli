package main

// Use termbox-go to add a simple UI
// two views - main results and episodes
// TODO
// add a way to scroll through episodes
// add a way to look at more than 10 results.
// add a way to look at the details for a specific episode
import (
	"fmt"
	"strings"

	"github.com/nsf/termbox-go"
	"github.com/yanfali/go-tvdb"
)

var (
	keyZero  = rune('0')
	keyOne   = rune('1')
	keyTwo   = rune('2')
	keyThree = rune('3')
	keyFour  = rune('4')
	keyFive  = rune('5')
	keySix   = rune('6')
	keySeven = rune('7')
	keyEight = rune('8')
	keyNine  = rune('9')
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

func updateScreen(tx *termboxState, fn func(tx *termboxState)) {
	drawClear()
	fn(tx)
	printConsoleString(tx)
	drawFlush()
}

func drawEpisode(tx *termboxState) {
	width, _ := termbox.Size()
	center := width / 2
	printTermboxString(center-len(title)/2, 0, title)
	series := tx.results.Series[tx.index]

	printTermboxString(1, 2, printSeries(series))
	allEpisodes := []tvdb.Episode{}
	for _, episodes := range series.Seasons {
		allEpisodes = append(allEpisodes, episodes...)
	}
	By(seasonAndEpisode).Sort(allEpisodes)
	tx.consoleMsg = fmt.Sprintf("%d Seasons. %d Episodes", len(series.Seasons), len(allEpisodes))
	for row, episode := range allEpisodes {
		printTermboxString(1, 4+row, fmt.Sprintf("%2v %2v %-70s", episode.SeasonNumber, episode.EpisodeNumber, episode.EpisodeName))
	}
}

func drawAll(tx *termboxState) {
	width, _ := termbox.Size()
	center := width / 2
	printTermboxString(center-len(title)/2, 1, title)
	printTermboxString(1, 3, fmt.Sprintf("%-5s %-40s %-10s %-15s %s", "Index", "Title", "Language", "First Aired", "Genre"))
	for row, series := range tx.results.Series {
		printTermboxString(1, 4+row, fmt.Sprintf("%-5d %s", row, printSeries(series)))
	}
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
