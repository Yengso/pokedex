package pokeapi

import (
	"encoding/json"
	"net/http"
	"io"
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

// This helper function takes the api url and returns a page form or err.

func PokedexAPI(url string) (Page, error) {
	var emptyPage Page
	defaultApi := "https://pokeapi.co/api/v2/location-area/"

	finalURL := url

	if len(url) == 0 {
		finalURL = defaultApi
	}

	resp, err := http.Get(finalURL)
	if err != nil {
		return emptyPage, err
	}
	defer resp.Body.Close()

	apiBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return emptyPage, err
	}

	apiForm := Page{}
	err = json.Unmarshal(apiBytes, &apiForm)
	if err != nil {
		return emptyPage, err
	}

	return apiForm, nil
}