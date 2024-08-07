package http

import (
	"encoding/json"
	"log"
	"net/http"
)

type UserService struct {
	mux *http.ServeMux
}

func NewUserService() (*UserService, error) {
	mux := http.NewServeMux()

	us := &UserService{
		mux: mux,
	}

	mux.Handle("/users", us)

	return us, nil
}

func (us *UserService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		us.post(w, r)
	}
}

func (us *UserService) post(w http.ResponseWriter, r *http.Request) {
	id := 0

	data := map[string]string{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf("/users POST: unable to decode JSON: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expected := []string{"country", "email", "first_name", "last_name", "nickname", "password"}
	for _, key := range expected {
		_, ok := data[key]
		if !ok {
			log.Printf("/users POST: unable to create as attribute %q was missing", key)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	log.Printf("/users POST: created %q", id)
	w.WriteHeader(http.StatusCreated)
}
