package config

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func ProducerAPIConfig() string {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load configuration")
	}

	producerUrl := env.SERVER + "/api/producer"
	return producerUrl
}

// Call producer API and return the list of all producers
func FetchProducers() {
	apiUrl := ProducerAPIConfig()

	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&Producers); err != nil {
		fmt.Println(err)
	}
}

// Call producer API and return producer matching ID
func FetchIndividual(id string) io.ReadCloser {
	apiUrl := ProducerAPIConfig() + "/" + id
	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}

	return res.Body
}

func BottleAPIConfig() string {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load configuration")
	}

	bottleUrl := env.SERVER + "/api/bottle"
	return bottleUrl
}

// Retrieve a list of bottles from an API endpoint
// The list of bottles is unmarshalled from the response body and stored in the bottles variable.
func FetchBottles() {
	res, err := http.Get(BottleAPIConfig())

	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&Bottles); err != nil {
		fmt.Println(err)
	}
}
