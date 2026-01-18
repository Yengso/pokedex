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

var cache = pokecache.NewCache(5 * time.Second)

// This helper function takes the api url and returns a page form or err.

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

func PokemonAPI(locURL string) (LocationArea, error) {
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