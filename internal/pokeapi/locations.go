package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/deimerin/pkdxcli/internal/pokecache"
)

type LocationResponse struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

// New Cache Definition
var cache = pokecache.NewCache(5 * time.Second)

func FetchLocations(url string) ([]string, string, string, error) {

	var locations LocationResponse
	var names []string

	// TRY CACHE FIRST
	if data, ok := cache.Get(url); ok {
		if err := json.Unmarshal(data, &locations); err != nil {
			return nil, "", "", err
		}

		for _, location := range locations.Results {
			names = append(names, location.Name)
		}

		if locations.Previous != nil {
			return names, locations.Next, locations.Previous.(string), nil
		}

		return names, locations.Next, "", nil

	}

	// API GET
	res, err := http.Get(url)
	if err != nil {
		return nil, "", "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, "", "", err
	}

	// ADD DATA TO THE CACHE
	cache.Add(url, body)

	// Unmarshal response
	if err = json.Unmarshal(body, &locations); err != nil {
		return nil, "", "", err
	}

	for _, location := range locations.Results {
		names = append(names, location.Name)
	}

	if locations.Previous != nil {
		return names, locations.Next, locations.Previous.(string), nil
	}

	return names, locations.Next, "", nil
}
