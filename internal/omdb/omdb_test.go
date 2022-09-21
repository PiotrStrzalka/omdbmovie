package omdb_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/piotrstrzalka/omdbmovie/internal/omdb"
)

func TestFetchPlot(t *testing.T) {
	fetcher := omdb.NewOMDBFetcher("http://www.omdbapi.com/", "")

	plot, err := fetcher.FetchPlot(context.Background(), "tt0011108")
	if err != nil {
		t.FailNow()
	}

	fmt.Println("Plot: ", plot)
}
