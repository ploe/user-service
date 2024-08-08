package main

import (
	"log"
	"user-service/http"
)

func main() {
	us, err := http.NewUserService()
	if err != nil {
		log.Fatalf("Unable to create UserService %q", err.Error())
	}

	const addr = "0.0.0.0:8080"

	log.Printf("Coming up on %q", addr)
	log.Fatal(us.ListenAndServe(addr))
}
