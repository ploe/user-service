package http

import (
	"encoding/json"
	"fmt"
	"log"
)

func Blep() {
	req := map[string]string{
		"country":    "UK",
		"email":      "alice@bob.com",
		"first_name": "Alice",
		"last_name":  "Bob",
		"nickname":   "AB123",
		"password":   "f6b7e19e0d867de6c0391879050e8297165728d89d7c4e9e8839972b356c4d9d",
	}

	body, err := json.Marshal(req)
	if err != nil {
		log.Fatalf(err.Error())
	}

	fmt.Println(string(body))
}
