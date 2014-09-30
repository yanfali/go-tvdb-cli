package main

import (
	"fmt"
	"os"

	"github.com/yanfali/go-tvdb"
)

func main() {
	myTvdb := tvdb.NewTvdb()
	results, err := myTvdb.SearchSeries(os.Args[1], 10)
	if err != nil {
		fmt.Errorf("error", err)
	}
	for _, series := range results.Series {
		fmt.Println(series.SeriesName)
	}
}
