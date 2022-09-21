package csvreader

import (
	"context"
	"encoding/csv"
	"io"
	"log"

	"github.com/piotrstrzalka/omdbmovie/internal/movie"
)

func ReadCSV(c context.Context, reader io.Reader) <-chan movie.MovieRecord {
	output := make(chan movie.MovieRecord)
	go func() {
		cReader := csv.NewReader(reader)
		cReader.Comma = '\t'

		defer func() {
			close(output)
		}()
		for {
			record, err := cReader.Read()
			if err != nil {
				if err == io.EOF {
					// log.Println("EOF detected")
					break
				}
				//todo increase error counter
				// log.Printf("Error while parsing csv: %v", err)
				continue
			}

			select {
			case output <- movie.MovieRecord(record):
			case <-c.Done():
				log.Println("Context cancelled")
				return
			}
		}
	}()
	return output
}
