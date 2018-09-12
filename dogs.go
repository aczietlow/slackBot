package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// Gets dog images
func getDog() string {
	resp, err := http.Get("https://dog.ceo/api/breeds/image/random")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	input, err := ioutil.ReadAll(resp.Body)
	var dogs DogResp
	json.Unmarshal(input, &dogs)
	return dogs.Message
}
