package main

import (
	"fmt"
	"os"

	"github.com/garfunkel/go-tvdb"
)

func main() {
	results, err := tvdb.SearchSeries(os.Args[1], 10)
	if err != nil {
		fmt.Errorf("error", err)
	}
	for _, series := range results.Series {
		fmt.Println(series.SeriesName)
	}
}
