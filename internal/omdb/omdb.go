package omdb

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

type OMDBFetcher struct {
	address string
	aPIkey  string
	client  http.Client
}

func NewOMDBFetcher(address, APIkey string) *OMDBFetcher {
	client := http.Client{
		Timeout: 10 * time.Second,
	}

	return &OMDBFetcher{
		address: address,
		aPIkey:  APIkey,
		client:  client,
	}
}

func (f *OMDBFetcher) FetchPlot(c context.Context, id string) (string, error) {
	url := f.address + "?apikey=" + f.aPIkey + "&i=" + id

	req, err := http.NewRequestWithContext(c, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := f.client.Do(req)
	if err != nil {
		return "", err
	}

	decoder := json.NewDecoder(resp.Body)
	defer resp.Body.Close()

	type respBody struct {
		Plot string `json:"Plot"`
	}
	var data respBody

	if err := decoder.Decode(&data); err != nil {
		return "", err
	}
	return data.Plot, nil
}
