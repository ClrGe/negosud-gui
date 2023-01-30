package data

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var Token string

// AuthGetRequest sends a req to the endpoint (route parameter) and returns the response body
func AuthGetRequest(route string) io.ReadCloser {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("Could not load config")
	}
	Token = env.KEY
	var server string

	if env.ENV == "dev" {
		server = env.SERVER_DEV
	} else {
		server = env.SERVER_PROD
	}

	url := server + "/" + route

	req, err := http.NewRequest("GET", url, nil)
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

func AuthPostRequest(route string, body *bytes.Buffer) int {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Println("Could not load config")
	}

	Token = env.KEY
	var server string

	if env.ENV == "dev" {
		server = env.SERVER_DEV
	} else {
		server = env.SERVER_PROD
	}

	url := server + "/" + route

	req, err := http.NewRequest("POST", url, body)
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

func LoginAndSaveToken(email string, password string) int {
	env, err := LoadConfig(".")
	if err != nil {
		fmt.Print(err, "Error loading config")
	}

	var response int
	var token string
	var server string

	if env.ENV == "dev" {
		server = env.SERVER_DEV
	} else {
		server = env.SERVER_PROD
	}

	hostname, _ := os.Hostname()
	url := server + "/authentication/login"

	userInfo := &User{
		Email:    email,
		Password: password,
	}

	jsonValue, err := json.Marshal(userInfo)
	if err != nil {
		Logger(true, "LOGIN Method", err.Error())
		fmt.Println(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Print(err, "Error creating request")
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		Logger(true, "LOGIN", err.Error())
		fmt.Print(err, "Error sending request")
	}

	if resp.StatusCode != 200 {
		message := "Login failed from origin " + hostname + "\n"
		fmt.Print(bytes.NewBuffer(jsonValue))
		fmt.Print(err)
		Logger(true, "LOGIN ", message)
		response = resp.StatusCode
		return response
	}
	Logger(false, "LOGIN", "Login successful from origin "+hostname+"\n")
	body, err := ioutil.ReadAll(resp.Body)
	token = string(body)

	if err != nil {
		fmt.Print(err, "Error retrieving token from response")
	}
	defer resp.Body.Close()
	SaveConfig("KEY", token)

	response = resp.StatusCode
	return response
}

//func SendInReqBody(route string, body *bytes.Buffer) int {
//	env, err := LoadConfig(".")
//	if err != nil {
//		fmt.Println("Could not load config")
//	}
//
//	Token = env.KEY
//	var server string
//
//	if env.ENV == "dev" {
//		server = env.SERVER_DEV
//	} else {
//		server = env.SERVER_PROD
//	}
//
//	url := server + "/" + route
//
//	req, err := http.NewRequest("POST", url, body)
//	if err != nil {
//		panic(err)
//	}
//	req.Header.Add("Authorization", "Bearer "+Token)
//	req.Header.Set("Content-Type", "application/json")
//
//	// send the request
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		panic(err)
//	}
//	if resp.StatusCode == 200 || resp.StatusCode == 201 {
//		fmt.Println("Success")
//	} else {
//		fmt.Println("Failed")
//		return resp.StatusCode
//	}
//
//	return resp.StatusCode
//}
