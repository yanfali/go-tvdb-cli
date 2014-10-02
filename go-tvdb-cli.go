package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/nsf/termbox-go"
	"github.com/yanfali/go-tvdb"
)

var RowTemplate = `{{.SeriesName | printf "%-40s"}} {{.Language | printf "%-10s"}} {{.FirstAired | printf "%-15s"}} {{.Genre}}`
var RowCompiled = template.Must(template.New("row").Parse(RowTemplate))

var config tvdb.TvdbConfig

var inEpisodes = false

func printSeries(data interface{}) string {

	w := bytes.NewBuffer([]byte{})
	err := RowCompiled.Execute(w, data)
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
			Value: "",
			Usage: "thetvdb.com [uses default from go-tvdb]",
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

		if len(results.Series) == 0 {
			log.Printf("No results found for %q", c.Args()[0])
			return
		}
		if err := termbox.Init(); err != nil {
			panic(err)
		}
		defer termbox.Close()

		drawAll(results)
		var lastKey rune
	loop:
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc:
					if inEpisodes {
						inEpisodes = false
						drawAll(results)
						break
					}
					break loop
				}
				switch ev.Ch {
				case keyZero, keyOne, keyTwo, keyThree, keyFour, keyFive, keySix, keySeven, keyEight, keyNine:
					inEpisodes = true
					lastKey = ev.Ch
					if len(results.Series[lastKey-'0'].Seasons) == 0 {
						if err := results.Series[lastKey-'0'].GetDetail(config); err != nil {
							drawAll(results)
							break
						}
					}
					drawEpisode(lastKey, results)
				}
			case termbox.EventResize:
				if inEpisodes {
					drawEpisode(lastKey, results)
				} else {
					drawAll(results)
				}
			}
		}
	}
	app.Run(os.Args)
}
