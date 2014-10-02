package main

// Episodes are returned in XML do not appear to guarentee any sort of sort order.
// Use go's built in sort infrastructure to sort by Season and Episode for presentation
import (
	"sort"

	"github.com/yanfali/go-tvdb"
)

func (by By) Sort(episodes []tvdb.Episode) {
	es := &episodeSorter{
		episodes: episodes,
		by:       by,
	}
	sort.Sort(es)
}

type By func(e1, e2 *tvdb.Episode) bool

type episodeSorter struct {
	episodes []tvdb.Episode
	by       func(e1, e2 *tvdb.Episode) bool
}

func (s *episodeSorter) Len() int {
	return len(s.episodes)
}

func (s *episodeSorter) Swap(i, j int) {
	s.episodes[i], s.episodes[j] = s.episodes[j], s.episodes[i]
}

func (s *episodeSorter) Less(i, j int) bool {
	return s.by(&s.episodes[i], &s.episodes[j])
}

func seasonAndEpisode(e1, e2 *tvdb.Episode) bool {
	if e1.SeasonNumber == e2.SeasonNumber {
		// use Episode Number as tie breaker
		return e1.EpisodeNumber < e2.EpisodeNumber
	}
	return e1.SeasonNumber < e2.SeasonNumber
}
