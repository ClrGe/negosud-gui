package data

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func UserAPIConfig() string {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load configuration")
	}

	userUrl := env.SERVER + "/api/users"
	return userUrl
}

func ProducerAPIConfig() string {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load configuration")
	}

	producerUrl := env.SERVER + "/api/producer"
	return producerUrl
}

func OrderAPIConfig() string {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load configuration")
	}

	orderUrl := env.SERVER + "/api/orders"
	return orderUrl
}

func CustomerOrderAPIConfig() string {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load configuration")
	}

	orderUrl := env.SERVER + "/api/orders/customers"
	return orderUrl
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
func FetchIndividualProducer(id string) io.ReadCloser {
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

func LoginAPIConfig(email string, password string) string {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load configuration")
	}

	bottleUrl := env.SERVER + "/api/authentication/login?email=" + email + "&password=" + password
	return bottleUrl
}

func UpdateBottleAPI() string {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load configuration")
	}

	bottleUrl := env.SERVER + "/api/updatebottle/"
	return bottleUrl
}

func UpdateProducerAPI() string {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("cannot load configuration")
	}

	bottleUrl := env.SERVER + "/api/updateproducer/"
	return bottleUrl
}

// Call producer API and return producer matching ID
func FetchIndividualBottle(id string) io.ReadCloser {
	apiUrl := BottleAPIConfig() + "/" + id

	res, err := http.Get(apiUrl)
	if err != nil {
		fmt.Println(err)
	}

	return res.Body
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
