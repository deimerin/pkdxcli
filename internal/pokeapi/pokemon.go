package pokeapi

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
)

var basePokemonURL string = "https://pokeapi.co/api/v2/pokemon/"

type Pokemon struct {
	Name           string `json:"name"`
	BaseExperience int    `json:"base_experience"`
	Height         int    `json:"height"`
	Weight         int    `json:"weight"`
	Stats          []struct {
		BaseStat int `json:"base_stat"`
		Effort   int `json:"effort"`
		Stat     struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"stat"`
	} `json:"stats"`
	Types []struct {
		Slot int `json:"slot"`
		Type struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"type"`
	} `json:"types"`
}

func (p Pokemon) Catch() bool {

	difficulty := p.BaseExperience
	roll := rand.Intn(100)

	if difficulty < 10 {
		difficulty = 10
	}
	if difficulty > 90 {
		difficulty = 90
	}

	catchChance := 100 - difficulty

	// fmt.Printf("\n- roll: %v\n- chance: %v\n", roll, catchChance)

	return catchChance > roll
}

func GetPokemon(name string) (Pokemon, error) {

	var pokemon Pokemon

	res, err := http.Get(basePokemonURL + name)
	if err != nil {
		return Pokemon{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}

	// Unmarshal json
	if err = json.Unmarshal(body, &pokemon); err != nil {
		return Pokemon{}, err
	}

	return pokemon, nil

}
