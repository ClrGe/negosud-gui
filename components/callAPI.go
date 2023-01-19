package components

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Call bottle API and return the list of all bottles
func fetchBottles() {
	env, err := LoadConfig(".")
	res, err := http.Get(env.SERVER + "/api/bottle")

	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&bottles); err != nil {
		fmt.Println(err)
	}
}

// Call producer API and return individual producer
func fetchProducer() {
	env, err := LoadConfig(".")
	res, err := http.Get(env.SERVER + "/api/producer/1")

	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&producer); err != nil {
		fmt.Println(err)
	}
}

// Call producer API and return the list of all producers
func fetchProducers() {
	env, err := LoadConfig(".")
	res, err := http.Get(env.SERVER + "/api/producer")

	if err != nil {
		fmt.Println(err)
		fmt.Println(err)
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(&producers); err != nil {
		fmt.Println(err)
	}
}
