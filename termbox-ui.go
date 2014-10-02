package main

// Use termbox-go to add a simple UI
// two views - main results and episodes
// TODO
// add a way to scroll through episodes
// add a way to look at more than 10 results.
// add a way to look at the details for a specific episode
import (
	"fmt"

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

func drawEpisode(ch rune, results tvdb.SeriesList) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	width, _ := termbox.Size()
	center := width / 2
	printTermboxString(center-len(title)/2, 0, title)
	var index = ch - '0'
	printTermboxString(1, 2, printSeries(results.Series[index]))
	allEpisodes := []tvdb.Episode{}
	for _, episodes := range results.Series[index].Seasons {
		allEpisodes = append(allEpisodes, episodes...)
	}
	By(seasonAndEpisode).Sort(allEpisodes)
	for row, episode := range allEpisodes {
		printTermboxString(1, 4+row, fmt.Sprintf("%2v %2v %-70s", episode.SeasonNumber, episode.EpisodeNumber, episode.EpisodeName))
	}
	termbox.Flush()
}

func drawAll(results tvdb.SeriesList) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	width, _ := termbox.Size()
	center := width / 2
	printTermboxString(center-len(title)/2, 1, title)
	printTermboxString(1, 3, fmt.Sprintf("%-5s %-40s %-10s %-15s %s", "Index", "Title", "Language", "First Aired", "Genre"))
	for row, series := range results.Series {
		printTermboxString(1, 4+row, fmt.Sprintf("%-5d %s", row, printSeries(series)))
	}
	termbox.Flush()
}

func printTermboxString(x, y int, s string) {
	for _, r := range s {
		termbox.SetCell(x, y, r, termbox.ColorWhite, termbox.ColorDefault)
		x += 1
	}
}
