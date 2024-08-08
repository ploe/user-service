package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

/*
TestSingleUsersGet: Given I have created a User when I call the GET
method then the User I have created will return with the following
fields populated: created_at, country, email, first_name, last_name,
nickname, password and updated_at with the values I created them with
and the types/formats specified in ./docs/endpoints/users/SCHEMA.md
*/
func TestSingleUsersGet(t *testing.T) {
	us, err := NewUserService()
	if err != nil {
		t.Fatal(err.Error())
	}

	data := map[string]string{
		"country":    "UK",
		"email":      "alice@bob.com",
		"first_name": "Alice",
		"last_name":  "Bob",
		"nickname":   "AB123",
		"password":   "f6b7e19e0d867de6c0391879050e8297165728d89d7c4e9e8839972b356c4d9d",
	}

	req_body, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err.Error())
	}

	post, err := http.NewRequest("POST", "/users", bytes.NewReader(req_body))
	if err != nil {
		t.Fatal(err.Error())
	}

	us.ServeHTTP(httptest.NewRecorder(), post)

	get, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	w := httptest.NewRecorder()
	us.ServeHTTP(w, get)

	resp_body := []map[string]string{}

	err = json.NewDecoder(w.Body).Decode(&resp_body)
	if err != nil {
		t.Fatal(err.Error())
	}

	user := resp_body[0]

	for key, expected := range data {
		got := user[key]

		if expected != got {
			t.Fatalf("Expected attribute %q to be %q but got %q", key, expected, got)
		}
	}

	for _, key := range []string{"created_at", "updated_at"} {
		_, err = time.Parse(DtLayout, user[key])
		if err != nil {
			t.Fatalf("Expected attribute %q to be in correct layout but got %q instead", key, user[key])
		}
	}

	err = uuid.Validate(user["id"])
	if err != nil {
		t.Fatalf("Expected 'id' to be a uuid but %q", err.Error())
	}
}

/*
TestCorrectUsersPost: Given I have rendered a request body with the
attributes country, email, first_name, last_name, nickname and
password when I send the request body via the HTTP POST method then
in the response the HTTP status code will be set to 201 Created.
*/
func TestCorrectUsersPost(t *testing.T) {
	data := map[string]string{
		"country":    "UK",
		"email":      "alice@bob.com",
		"first_name": "Alice",
		"last_name":  "Bob",
		"nickname":   "AB123",
		"password":   "f6b7e19e0d867de6c0391879050e8297165728d89d7c4e9e8839972b356c4d9d",
	}

	body, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err.Error())
	}

	req, err := http.NewRequest("POST", "/users", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err.Error())
	}

	w := httptest.NewRecorder()

	us, err := NewUserService()
	if err != nil {
		t.Fatal(err.Error())
	}

	us.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("Unexpected error code. Got %d, %d expected.", resp.StatusCode, http.StatusCreated)
	}
}

/*
TestWrongUsersPost: Given I have rendered a request body missing any
of the following attributes: country, email, first_name, last_name,
nickname and password when I send the request body via the HTTP POST
method then in the response the HTTP status code will be set to 400 Bad Request.
*/
func TestWrongUsersPost(t *testing.T) {
	expected := []string{"country", "email", "first_name", "last_name", "nickname", "password"}

	data := map[string]string{
		"country":    "UK",
		"email":      "alice@bob.com",
		"first_name": "Alice",
		"last_name":  "Bob",
		"nickname":   "AB123",
		"password":   "f6b7e19e0d867de6c0391879050e8297165728d89d7c4e9e8839972b356c4d9d",
	}

	for _, attribute := range expected {
		wrong := map[string]string{}
		for key, value := range data {
			if key == attribute {
				continue
			}

			wrong[key] = value
		}

		body, err := json.Marshal(wrong)
		if err != nil {
			t.Fatal(err.Error())
		}

		req, err := http.NewRequest("POST", "/users", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err.Error())
		}

		w := httptest.NewRecorder()

		us, err := NewUserService()
		if err != nil {
			t.Fatal(err.Error())
		}

		us.ServeHTTP(w, req)

		resp := w.Result()

		if resp.StatusCode != http.StatusBadRequest {
			t.Fatalf("Unexpected error code. Got %d, %d expected.", resp.StatusCode, http.StatusBadRequest)
		}
	}
}
