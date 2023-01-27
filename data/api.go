package data

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

var Token string

// AuthGetRequest sends a req to the endpoint (route parameter) and returns the response body
func AuthGetRequest(route string) io.ReadCloser {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("Could not load config")
	}
	Token = env.APIKEY

	req, err := http.NewRequest("GET", env.SERVER+"/"+route, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+Token)
	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		fmt.Println("Success")
	} else {
		fmt.Println("Failed")
		return nil
	}
	return resp.Body
}

func AuthPostRequest(route string, body []byte) int {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("Could not load config")
	}

	Token = env.APIKEY

	req, err := http.NewRequest("POST", env.SERVER+"/"+route, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+Token)
	req.Header.Set("Content-Type", "application/json")

	// send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode == 200 || resp.StatusCode == 201 {
		fmt.Println("Success")
	} else {
		fmt.Println("Failed")
		return resp.StatusCode
	}

	return resp.StatusCode
}
