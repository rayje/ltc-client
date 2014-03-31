package main

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "net/http"
    "encoding/json"
    "os"
)

type ApigeeToken struct {
	ApigeeToken    string	`json:"apigeeToken"`
}

func buildApigeeRequest(apigee *ApigeeConfig) (http.Request, error) {
	apigeeUrl := "https://ecollege-prod.apigee.net"
	apigeeUrl += "/latency-test/idm/service/identity/login/basic"
    apigeeUrl += "?apikey=" + apigee.Apikey

    dataFmt := "{\"email\":\"%s\",\"password\":\"%s\"}"
    dataString := fmt.Sprintf(dataFmt, apigee.Email, apigee.Password)

    req, err := http.NewRequest("POST", apigeeUrl, bytes.NewBufferString(dataString))
    if err != nil {
    	return *req, err
    }

    req.Header.Add("Content-Type", "application/json")
    return *req, err
}

func getApigeeToken(config *Config) (string, error) {
	apigeeConfig := config.Apigee
	if apigeeConfig.Email == "" {
		fmt.Println("Apigee Error: Email not set")
		os.Exit(1)
	}
	if apigeeConfig.Password == "" {
		fmt.Println("Apigee Error: Password not set")
		os.Exit(1)
	}

	req, err := buildApigeeRequest(&apigeeConfig)
    if err != nil {
  		return "", err
  	}

	client := &http.Client{}
    resp, err := client.Do(&req)
    if err != nil {
  		return "", err
  	}

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
  		return "", err
  	}

	var apigee ApigeeToken
  	err = json.Unmarshal(body, &apigee)
  	if err != nil {
  		return "", err
  	}

  	return apigee.ApigeeToken, nil
}