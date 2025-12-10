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

type LocationAreaData struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	GameIndex            int    `json:"game_index"`
	EncounterMethodRates []struct {
		EncounterMethod struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"encounter_method"`
		VersionDetails []struct {
			Rate    int `json:"rate"`
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
		} `json:"version_details"`
	} `json:"encounter_method_rates"`
	Location struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"location"`
	Names []struct {
		Name     string `json:"name"`
		Language struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"language"`
	} `json:"names"`
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			Version struct {
				Name string `json:"name"`
				URL  string `json:"url"`
			} `json:"version"`
			MaxChance        int `json:"max_chance"`
			EncounterDetails []struct {
				MinLevel        int   `json:"min_level"`
				MaxLevel        int   `json:"max_level"`
				ConditionValues []any `json:"condition_values"`
				Chance          int   `json:"chance"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
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

// This seems appropiate here, I guess
func FetchLocationArea(url string) (pokemon []string, err error) {

	var locationAreaData LocationAreaData
	var pokeList []string

	// if cache ok!
	if data, ok := cache.Get(url); ok {
		if err := json.Unmarshal(data, &locationAreaData); err != nil {
			return nil, err
		}

		for _, name := range locationAreaData.PokemonEncounters {
			pokeList = append(pokeList, name.Pokemon.Name)
		}

		return pokeList, nil

	}

	// no cache, then api call

	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// Add data to cache
	cache.Add(url, body)

	// Unmarshal and name retrieval

	if err := json.Unmarshal(body, &locationAreaData); err != nil {
		return nil, err
	}

	for _, name := range locationAreaData.PokemonEncounters {
		pokeList = append(pokeList, name.Pokemon.Name)
	}

	return pokeList, nil

}
