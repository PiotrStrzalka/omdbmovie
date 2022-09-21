package movie

import (
	"flag"
	"reflect"
	"sort"
	"strings"
)

var titleType = flag.String("titleType", "", "filter on titleType column")
var primaryTitle = flag.String("primaryTitle", "", "filter on primaryTitle column")
var originalTitle = flag.String("originalTitle", "", "filter on originalTitle column")
var genre = flag.String("genre", "", "filter on genre column")
var startYear = flag.Int("startYear", 0, "filter on startYear column")
var endYear = flag.Int("endYear", 0, "filter on endYear column")
var runtimeMinutes = flag.Int("runtimeMinutes", 0, "filter on runtimeMinutes column")
var genres = flag.String("genres", "", "filter on genres column")

func (m *MovieFilter) AddFlagFilters() {
	if isFlagPassed("titleType") {
		m.AddStage(TitleTypeEquals(*titleType))
	}
	if isFlagPassed("primaryTitle") {
		m.AddStage(PrimaryTitleEquals(*primaryTitle))
	}
	if isFlagPassed("originalTitle") {
		m.AddStage(OriginalTitleEquals(*originalTitle))
	}
	if isFlagPassed("genre") {
		m.AddStage(GenreContains(*genre))
	}
	if isFlagPassed("startYear") {
		m.AddStage(StartYearEquals(*startYear))
	}
	if isFlagPassed("endYear") {
		m.AddStage(EndYearEquals(*endYear))
	}
	if isFlagPassed("genres") {
		m.AddStage(GenreEquals(strings.Split(*genres, ",")))
	}
	if isFlagPassed("runtimeMinutes") {
		m.AddStage(RuntimeMinutesEquals(*runtimeMinutes))
	}
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

type filterFunc func(m *MovieRecord) bool

type MovieFilter struct {
	filterFn []filterFunc
}

func NewFilter() *MovieFilter {
	return &MovieFilter{}
}

func (m *MovieFilter) AddStage(f filterFunc) *MovieFilter {
	m.filterFn = append(m.filterFn, f)
	return m
}

func (m *MovieFilter) Stages() int {
	return len(m.filterFn)
}

func PrimaryTitleEquals(value string) filterFunc {
	return func(mr *MovieRecord) bool {
		return *mr.PrimaryTitle() == value
	}
}

func OriginalTitleEquals(value string) filterFunc {
	return func(mr *MovieRecord) bool {
		return *mr.OriginalTitle() == value
	}
}

func TitleTypeEquals(value string) filterFunc {
	return func(mr *MovieRecord) bool {
		return *mr.TitleType() == value
	}
}

func GenreEquals(value []string) filterFunc {
	return func(mr *MovieRecord) bool {
		sort.Strings(value)
		genres := mr.Genres()
		sort.Strings(genres)
		return reflect.DeepEqual(value, genres)
	}
}

func StartYearEquals(value int) filterFunc {
	return func(mr *MovieRecord) bool {
		v, err := mr.StartYear()
		if err != nil {
			return false
		}
		return v == value
	}
}

func EndYearEquals(value int) filterFunc {
	return func(mr *MovieRecord) bool {
		v, err := mr.StartYear()
		if err != nil {
			return false
		}
		return v == value
	}
}

func GenreContains(value string) filterFunc {
	return func(mr *MovieRecord) bool {
		for _, g := range mr.Genres() {
			if g == value {
				return true
			}
		}
		return false
	}
}

func RuntimeMinutesEquals(value int) filterFunc {
	return func(mr *MovieRecord) bool {
		v, err := mr.RuntimeMinutes()
		if err != nil {
			return false
		}
		return v == value
	}
}

func (m *MovieFilter) Filter(mr *MovieRecord) bool {
	for _, f := range m.filterFn {
		if !f(mr) {
			return false
		}
	}
	return true
}
