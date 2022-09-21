package movie_test

import (
	"testing"

	"github.com/piotrstrzalka/omdbmovie/internal/movie"
)

func TestFilterPassAll(t *testing.T) {
	f := movie.NewFilter()

	data := movie.MovieRecord([]string{"Id", "TitleType", "Title"})
	if f.Filter(&data) != true {
		t.Fail()
	}
}

func TestFilterTitle(t *testing.T) {
	f := movie.NewFilter().AddStage(movie.PrimaryTitleEquals("Title1"))

	data := movie.MovieRecord([]string{"Id", "TitleType", "Title"})
	if f.Filter(&data) == true {
		t.Fail()
	}
}
