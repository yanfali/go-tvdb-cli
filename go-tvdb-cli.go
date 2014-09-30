package main

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
	"text/template"

	"github.com/codegangsta/cli"
	"github.com/yanfali/go-tvdb"
)

var RowTemplate = `"{{.SeriesName}}", "{{.Genre}}", "{{.Language}}", "{{.FirstAired}}"
`
var RowCompiled = template.Must(template.New("row").Parse(RowTemplate))

func printSeries(data interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', 0)
	err := RowCompiled.Execute(w, data)
	if err != nil {
		panic(err)
	}
	w.Flush()
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
		config := tvdb.NewDefaultTvdbConfig()
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
		fmt.Printf("%q, %q, %q, %q\n", "Title", "Genre", "Language", "First Aired")
		for _, series := range results.Series {
			printSeries(series)
		}
	}
	app.Run(os.Args)
}
