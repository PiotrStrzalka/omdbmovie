package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/piotrstrzalka/omdbmovie/internal/concurrency"
	"github.com/piotrstrzalka/omdbmovie/internal/csvreader"
	"github.com/piotrstrzalka/omdbmovie/internal/movie"
	"github.com/piotrstrzalka/omdbmovie/internal/omdb"
)

const (
	filteringWorkers = 5
	fetchingWorkers  = 5
	regexWorkers     = 5
)

func main() {

	maxTime := flag.Duration("maxRunTime", time.Hour*24*365, "maximum run time of the application")
	filePath := flag.String("filePath", "", "absolute path to the inflated title.basics.tsv.gz file")
	apikey := flag.String("apikey", "", "API key for omdb")
	plotFilter := flag.String("plotFilter", "(.*?)", "regex for filtering plot")

	flag.Parse()

	c, cancel := configuraExitOptions(context.Background(), *maxTime, syscall.SIGTERM)
	defer cancel()

	//CONFIGURE workers
	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatalf("Cannot open file: %v", err)
	}
	defer file.Close()

	movieFilter := movie.NewFilter()
	movieFilter.AddFlagFilters()

	if movieFilter.Stages() == 0 {
		log.Fatal("Some filters are welcome")
	}

	fetcher := omdb.NewOMDBFetcher("http://www.omdbapi.com/", *apikey)

	regexer, err := regexp.Compile(*plotFilter)
	if err != nil {
		log.Fatalf("Cannot parse regex: %v", err)
	}

	//PIPELINE chaining
	readChannel := csvreader.ReadCSV(c, file)

	filterChannel := concurrency.ProccessConcurrently(c, readChannel, func(c context.Context, data movie.MovieRecord) (movie.MovieRecord, error) {
		if movieFilter.Filter(&data) {
			return data, nil
		}
		return movie.MovieRecord{}, errors.New("")
	}, filteringWorkers)

	plotChannel := concurrency.ProccessConcurrently(c, filterChannel, func(c context.Context, data movie.MovieRecord) (movie.MovieResult, error) {
		plot, err := fetcher.FetchPlot(c, *data.Id())
		if err != nil {
			return movie.MovieResult{}, err
		}
		return movie.MovieResult{Id: *data.Id(), Title: *data.PrimaryTitle(), Plot: plot}, nil
	}, fetchingWorkers)

	regexChannel := concurrency.ProccessConcurrently(c, plotChannel, func(c context.Context, data movie.MovieResult) (movie.MovieResult, error) {
		if regexer.MatchString(data.Plot) {
			return data, nil
		}
		return movie.MovieResult{}, errors.New("data not match")
	}, regexWorkers)

	for res := range regexChannel {
		fmt.Println(res)
	}

	log.Println("Thanks for using")
}

func configuraExitOptions(c context.Context, maxTime time.Duration, signals ...os.Signal) (context.Context, func()) {
	c, stopFn := signal.NotifyContext(c, signals...)
	c, cancelFn := context.WithDeadline(c, time.Now().Add(maxTime))
	return c, func() {
		stopFn()
		cancelFn()
	}
}
