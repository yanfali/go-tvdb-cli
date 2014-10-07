package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"
	"time"

	"github.com/codegangsta/cli"
	"github.com/nsf/termbox-go"
	"github.com/yanfali/go-tvdb"
)

var (
	RowTemplate = `{{.SeriesName | printf "%-40s"}} {{.Language | printf "%-10s"}} {{.FirstAired | printf "%-15s"}} {{.Genre}}`
	RowCompiled = template.Must(template.New("row").Parse(RowTemplate))
	config      tvdb.TvdbConfig
)

func ellipsisString(source string, length int) string {
	if len(source) > length {
		return fmt.Sprintf("%."+string(length-4)+"s...", source)
	}
	return source
}

func printSeries(series *tvdb.Series) string {

	w := bytes.NewBuffer([]byte{})
	series.SeriesName = ellipsisString(series.SeriesName, 40)
	err := RowCompiled.Execute(w, series)
	if err != nil {
		panic(err)
	}
	return w.String()
}

func init() {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.Name}} {{if .Flags}}[global options] {{end}} "search title"

VERSION:
   {{.Version}}{{if or .Author .Email}}

AUTHOR:{{if .Author}}
  {{.Author}}{{if .Email}} - <{{.Email}}>{{end}}{{else}}
  {{.Email}}{{end}}{{end}}

COMMANDS:
   {{range .Commands}}{{.Name}}{{with .ShortName}}, {{.}}{{end}}{{ "\t" }}{{.Usage}}
   {{end}}{{if .Flags}}
GLOBAL OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`
}

func main() {
	app := cli.NewApp()
	app.Name = "go-tvdb-cli"
	app.Usage = "make CLI queries against the thetvdb.com"
	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "max-results, m",
			Value: 10,
			Usage: "Maximum Number of Results to Show",
		},
		cli.StringFlag{
			Name:  "apikey, k",
			Value: "90CCCAB7A2B7509E",
			Usage: "thetvdb.com API key",
		},
		cli.StringFlag{
			Name:  "language, l",
			Value: "en",
			Usage: "Default Language to search using",
		},
	}

	app.Action = func(c *cli.Context) {
		if len(c.Args()) == 0 {
			fmt.Printf("Error: Not enough parameters\n")
			cli.ShowAppHelp(c)
			return
		}
		config = tvdb.NewDefaultTvdbConfig()
		if apiKey := c.String("apikey"); apiKey != "" {
			log.Printf("Using APIKEY %q\n", apiKey)
			config.ApiKey = apiKey
		}
		if lang := c.String("language"); lang != "en" {
			log.Printf("Using Language %q", lang)
			config.Language = lang
		}
		myTvdb := tvdb.NewTvdbWithConfig(config)
		results, err := myTvdb.SearchSeries(c.Args()[0], c.Int("max-results"))
		if err != nil {
			fmt.Errorf("error", err)
		}
		var tx = &termboxState{results: &results}
		tx.consoleFn = func(*termboxState) string {
			return fmt.Sprintf("Displaying %d results out of %d", len(results.Series), c.Int("max-results"))
		}

		if len(results.Series) == 0 {
			log.Printf("No results found for %q", c.Args()[0])
			return
		}
		if err := termbox.Init(); err != nil {
			panic(err)
		}
		defer termbox.Close()

		updateScreen(tx, drawSeries)

		currentState := SeriesEventHandler
		ch := getPollEventChan()
	loop:
		for {
			select {
			case tx.ev = <-ch:
				if currentState = currentState(tx); currentState == nil {
					break loop
				}
			case <-time.After(time.Second / 2):
				cursorBlink(tx)
				termbox.Flush()
			}
		}
	}
	app.Run(os.Args)
}

func getPollEventChan() chan termbox.Event {
	var ch = make(chan termbox.Event)
	go func() {
		for {
			ch <- termbox.PollEvent()
		}
	}()
	return ch
}
