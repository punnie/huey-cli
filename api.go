package main

import (
	"bytes"
	"fmt"
	"log"

	"encoding/json"
	"net/http"
)

func RequestApi[T any](verb string, path string, payloadData interface{}) (T, error) {
	jsonPayload, err := json.Marshal(payloadData)

	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	url := fmt.Sprintf("%s%s", k.String("api.url"), path)

	req, err := http.NewRequest(verb, url, bytes.NewBuffer(jsonPayload))

	if err != nil {
		log.Fatalf("Error creating Request: %v", err)
	}

	authorizationString := fmt.Sprintf("Bearer %s", k.String("api.token"))

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authorizationString)

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Error requesting stuff")
	}

	defer resp.Body.Close()

	if resp.StatusCode/100 != 2 {
		log.Fatalf("Request failed with status: %s", resp.Status)
	}

	var responseData T

	if resp.StatusCode != http.StatusNoContent {
		if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
			log.Fatalf("Error decoding response JSON: %v", err)
		}
	}

	return responseData, nil
}
