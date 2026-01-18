package pokeapi

import (
	"fmt"
	"encoding/json"
	"net/http"
	"time"
	"io"
	"github.com/yengso/pokedexcli/internal/pokecache"
)

type Result struct {
	Name string	`json:"name"`
	Url	 string `json:"url"`
}

type Page struct {
	Results 	[]Result `json:"results"`
	Next 		string	 `json:"next"`
	Previous  	string	 `json:"previous"`
}

type LocationArea struct {
	PokemonEncounters []struct {
		Pokemon struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		} `json:"pokemon"`
		VersionDetails []struct {
			EncounterDetails []struct {
				Chance          int   `json:"chance"`
				ConditionValues []any `json:"condition_values"`
				MaxLevel        int   `json:"max_level"`
				Method          struct {
					Name string `json:"name"`
					URL  string `json:"url"`
				} `json:"method"`
				MinLevel int `json:"min_level"`
			} `json:"encounter_details"`
		} `json:"version_details"`
	} `json:"pokemon_encounters"`
}

type Pokemon struct {
	BaseExperience int `json:"base_experience"`
	Height    int `json:"height"`
	Name          string `json:"name"`
	Stats []struct {
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
	Weight int `json:"weight"`
}

var cache = pokecache.NewCache(5 * time.Second)

func LocationsAPI(url string) (Page, error) {
	var emptyPage Page
	defaultApi := "https://pokeapi.co/api/v2/location-area/"

	finalURL := url
	var apiBytes []byte

	if len(url) == 0 {
		finalURL = defaultApi
	}

	cacheData, ok := cache.Get(finalURL)
	if ok == true {
		apiBytes = cacheData
		fmt.Println("You just used cached data!")
	}
	if ok == false {
		resp, err := http.Get(finalURL)
		if err != nil {
			return emptyPage, err
		}
		defer resp.Body.Close()

		apiBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return emptyPage, err
		}

		cache.Add(finalURL, apiBytes)
	}

	apiForm := Page{}
	err := json.Unmarshal(apiBytes, &apiForm)
	if err != nil {
		return emptyPage, err
	}

	return apiForm, nil
}

func ExploreAPI(locURL string) (LocationArea, error) {
	var emptyLocArea LocationArea
	var apiBytes [] byte
	var startURL = "https://pokeapi.co/api/v2/location-area/"

	fullURL := startURL + locURL

	if len(locURL) == 0 {
		fmt.Println("No location written")
		return emptyLocArea, nil
	}

	cacheData, ok := cache.Get(fullURL)
	if ok == true {
		apiBytes = cacheData
		fmt.Println("You just used cached data!")
	}
	if ok == false {
		resp, err := http.Get(fullURL)
		if err != nil {
			return emptyLocArea, err
		}
		defer resp.Body.Close()

		apiBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return emptyLocArea, err
		}
		cache.Add(fullURL, apiBytes)
	}

	apiForm := LocationArea{}
	err := json.Unmarshal(apiBytes, &apiForm)
	if err != nil {
		return emptyLocArea, err
	}

	return apiForm, nil 
}

func PokemonAPI(pokemon string) (Pokemon, error) {
	var emptyPokemon Pokemon
	var apiBytes [] byte
	var startURL = "https://pokeapi.co/api/v2/pokemon/"

	fullURL := startURL + pokemon

	if len(pokemon) == 0 {
		fmt.Println("No pokemon mentioned")
		return emptyPokemon, nil
	}

	cacheData, ok := cache.Get(fullURL)
	if ok == true {
		apiBytes = cacheData
	}
	if ok == false {
		resp, err := http.Get(fullURL)
		if err != nil {
			return emptyPokemon, err
		}
		defer resp.Body.Close()

		apiBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return emptyPokemon, err
		}
		cache.Add(fullURL, apiBytes)
	}

	apiForm := Pokemon{}
	err := json.Unmarshal(apiBytes, &apiForm)
	if err != nil {
		return emptyPokemon, err
	}

	return apiForm, nil 
}




