package pokeapi

import (
	"encoding/json"
	"io"
	"net/http"
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

func FetchLocations(url string) ([]string, string, string, error) {

	var locations LocationResponse
	var names []string

	// API GET
	res, err := http.Get(url)
	if err != nil {
		return nil, "", "", err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return nil, "", "", err
	}

	// Unmarshal response
	err = json.Unmarshal(body, &locations)
	if err != nil {
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
