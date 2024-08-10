package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

/*
TestHealthcheckEqualNumber: Given I have set a healthcheck status
code a number of times when I call the GET method then the value
for that status code will be that number.
*/
func TestHealthcheckEqualNumber(t *testing.T) {
	us, err := NewUserService()
	if err != nil {
		t.Fatal(err.Error())
	}

	statuses := map[int]int{
		http.StatusOK:                  420,
		http.StatusNotFound:            69,
		http.StatusInternalServerError: 9001,
	}

	for status, count := range statuses {
		for i := 0; i < count; i++ {
			us.healthcheck.increment(status)
		}
	}

	get, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	w := httptest.NewRecorder()
	us.ServeHTTP(w, get)

	body := map[string]int{}

	err = json.NewDecoder(get.Body).Decode(&body)
	if err != nil {
		t.Fatal(err.Error())
	}

	for status, count := range statuses {
		key := fmt.Sprint(status)
		got := body[key]

		if got != count {
			t.Fatalf("expecting status %q to be %d but got %d instead", key, got, count)
		}

	}
}
