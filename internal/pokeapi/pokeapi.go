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

var cache = pokecache.NewCache(5 * time.Second)

// This helper function takes the api url and returns a page form or err.

func PokedexAPI(url string) (Page, error) {
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