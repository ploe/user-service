package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
)

/*
TestPatchModifiesUpdatedAt: Given I have created a User and modified
them with PATCH when I call the GET method then the User's created_at
and updated_at attributes will be no longer be the same.
*/
func TestPatchModifiesUpdatedAt(t *testing.T) {
	/* create user */

	us, err := NewUserService()
	if err != nil {
		t.Fatal(err.Error())
	}

	debut := map[string]string{
		"country":    "UK",
		"email":      "alice@bob.com",
		"first_name": "Alice",
		"last_name":  "Bob",
		"nickname":   "AB123",
		"password":   "f6b7e19e0d867de6c0391879050e8297165728d89d7c4e9e8839972b356c4d9d",
	}

	post_body, err := json.Marshal(debut)
	if err != nil {
		t.Fatal(err.Error())
	}

	post_req, err := http.NewRequest("POST", "/users", bytes.NewReader(post_body))
	if err != nil {
		t.Fatal(err.Error())
	}

	us.ServeHTTP(httptest.NewRecorder(), post_req)

	created_get_req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	/* after created get the user */

	created_get_resp := httptest.NewRecorder()
	us.ServeHTTP(created_get_resp, created_get_req)

	created_get_body := []map[string]string{}

	err = json.NewDecoder(created_get_resp.Body).Decode(&created_get_body)
	if err != nil {
		t.Fatal(err.Error())
	}

	debut_updated_at := created_get_body[0]["updated_at"]

	id := created_get_body[0]["id"]
	url := fmt.Sprintf("/users/%s", id)

	finale := map[string]string{
		"country":    "USA",
		"email":      "ken@bob.com",
		"first_name": "Ken",
		"last_name":  "Thompson",
		"nickname":   "ken",
		"password":   "b3bb4cd67f11e1f6350a5792c8a0f91c2e7920ab93ccd7e964d97d79ad9f8270",
	}

	patch_body, err := json.Marshal(&finale)
	if err != nil {
		t.Fatal(err.Error())
	}

	patch_req, err := http.NewRequest("PATCH", url, bytes.NewReader(patch_body))
	if err != nil {
		t.Fatal(err.Error())
	}

	patch_resp := httptest.NewRecorder()
	us.ServeHTTP(patch_resp, patch_req)

	/* after patch get users */

	patched_get_req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	patched_get_resp := httptest.NewRecorder()
	us.ServeHTTP(patched_get_resp, patched_get_req)

	patched_get_body := []map[string]string{}

	err = json.NewDecoder(patched_get_resp.Body).Decode(&patched_get_body)
	if err != nil {
		t.Fatal(err.Error())
	}

	finale_updated_at := patched_get_body[0]["updated_at"]

	if debut_updated_at == finale_updated_at {
		t.Fatalf("expected updated_at attribute to be modified but both were %q", debut_updated_at)
	}
}

/*
TestPatchModifiesAttributes: Given I have created a User and
modified any of the following attributes: country, email, first_name,
last_name, nickname and password when I call the GET method then the
modified attributes will be what I updated them to.
*/
func TestPatchModifiesAttributes(t *testing.T) {
	/* create user */

	us, err := NewUserService()
	if err != nil {
		t.Fatal(err.Error())
	}

	debut := map[string]string{
		"country":    "UK",
		"email":      "alice@bob.com",
		"first_name": "Alice",
		"last_name":  "Bob",
		"nickname":   "AB123",
		"password":   "f6b7e19e0d867de6c0391879050e8297165728d89d7c4e9e8839972b356c4d9d",
	}

	post_body, err := json.Marshal(debut)
	if err != nil {
		t.Fatal(err.Error())
	}

	post_req, err := http.NewRequest("POST", "/users", bytes.NewReader(post_body))
	if err != nil {
		t.Fatal(err.Error())
	}

	us.ServeHTTP(httptest.NewRecorder(), post_req)

	created_get_req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	/* after created get the user */

	created_get_resp := httptest.NewRecorder()
	us.ServeHTTP(created_get_resp, created_get_req)

	created_get_body := []map[string]string{}

	err = json.NewDecoder(created_get_resp.Body).Decode(&created_get_body)
	if err != nil {
		t.Fatal(err.Error())
	}

	id := created_get_body[0]["id"]
	url := fmt.Sprintf("/users/%s", id)

	/* patch each attribute on the user and check if it has updated */

	finale := map[string]string{
		"country":    "USA",
		"email":      "ken@bob.com",
		"first_name": "Ken",
		"last_name":  "Thompson",
		"nickname":   "ken",
		"password":   "b3bb4cd67f11e1f6350a5792c8a0f91c2e7920ab93ccd7e964d97d79ad9f8270",
	}

	for attribute, stale := range debut {
		/* patch the user */

		fresh := finale[attribute]

		patch_data := map[string]string{
			attribute: fresh,
		}

		patch_body, err := json.Marshal(&patch_data)
		if err != nil {
			t.Fatal(err.Error())
		}

		patch_req, err := http.NewRequest("PATCH", url, bytes.NewReader(patch_body))
		if err != nil {
			t.Fatal(err.Error())
		}

		patch_resp := httptest.NewRecorder()
		us.ServeHTTP(patch_resp, patch_req)

		/* after patch get users */

		patched_get_req, err := http.NewRequest("GET", "/users", nil)
		if err != nil {
			t.Fatal(err.Error())
		}

		patched_get_resp := httptest.NewRecorder()
		us.ServeHTTP(patched_get_resp, patched_get_req)

		patched_get_body := []map[string]string{}

		err = json.NewDecoder(patched_get_resp.Body).Decode(&patched_get_body)
		if err != nil {
			t.Fatal(err.Error())
		}

		/* has the user been modified like we expect? */

		current := patched_get_body[0][attribute]

		if (current == stale) || (current != fresh) {
			t.Fatalf("expected attribute %q to be modified from %q to %q but got %q", attribute, stale, fresh, current)
		}
	}

}

/*
TestPostAndPatchStatusIsNoContent: Given I have created a User
when I call the PATCH method then the HTTP status code will be 204
No Content.
*/
func TestPostAndPatchStatusIsNoContent(t *testing.T) {
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

	post_body, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err.Error())
	}

	post_req, err := http.NewRequest("POST", "/users", bytes.NewReader(post_body))
	if err != nil {
		t.Fatal(err.Error())
	}

	us.ServeHTTP(httptest.NewRecorder(), post_req)

	get_req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	get_resp := httptest.NewRecorder()
	us.ServeHTTP(get_resp, get_req)

	get_body := []map[string]string{}

	err = json.NewDecoder(get_resp.Body).Decode(&get_body)
	if err != nil {
		t.Fatal(err.Error())
	}

	user := get_body[0]
	url := fmt.Sprintf("/users/%s", user["id"])

	body := []byte("{}")

	patch_req, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err.Error())
	}

	patch_resp := httptest.NewRecorder()
	us.ServeHTTP(patch_resp, patch_req)

	status := patch_resp.Result().StatusCode

	if status != http.StatusNoContent {
		t.Fatalf("Unexpected error code. Got %d, %d expected.", status, http.StatusNoContent)
	}
}

/*
TestPatchStatusIsNotFound: Given I have not created a User when I
call the DELETE method then the HTTP status code will be 404 Not
Found.
*/
func TestPatchStatusIsNotFound(t *testing.T) {
	us, err := NewUserService()
	if err != nil {
		t.Fatal(err.Error())
	}

	url := fmt.Sprintf("/users/%s", uuid.NewString())

	body := []byte("{}")

	req, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err.Error())
	}

	w := httptest.NewRecorder()
	us.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Unexpected error code. Got %d, %d expected.", resp.StatusCode, http.StatusNotFound)
	}
}

/*
TestPostAndDeleteStatusIsNoContent: Given I have created a User
when I call the DELETE method then the HTTP status code will be 204
No Content.
*/
func TestPostAndDeleteStatusIsNoContent(t *testing.T) {
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

	post_body, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err.Error())
	}

	post_req, err := http.NewRequest("POST", "/users", bytes.NewReader(post_body))
	if err != nil {
		t.Fatal(err.Error())
	}

	us.ServeHTTP(httptest.NewRecorder(), post_req)

	get_req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	get_resp := httptest.NewRecorder()
	us.ServeHTTP(get_resp, get_req)

	get_body := []map[string]string{}

	err = json.NewDecoder(get_resp.Body).Decode(&get_body)
	if err != nil {
		t.Fatal(err.Error())
	}

	user := get_body[0]
	url := fmt.Sprintf("/users/%s", user["id"])

	delete_req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	delete_resp := httptest.NewRecorder()
	us.ServeHTTP(delete_resp, delete_req)

	status := delete_resp.Result().StatusCode

	if status != http.StatusNoContent {
		t.Fatalf("Unexpected error code. Got %d, %d expected.", status, http.StatusNoContent)
	}
}

/*
TestDeleteStatusIsNotFound: Given I have not created a
User when I call the DELETE method then the HTTP status code will be
404 Not Found.
*/
func TestDeleteStatusIsNotFound(t *testing.T) {
	us, err := NewUserService()
	if err != nil {
		t.Fatal(err.Error())
	}

	url := fmt.Sprintf("/users/%s", uuid.NewString())

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	w := httptest.NewRecorder()
	us.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("Unexpected error code. Got %d, %d expected.", resp.StatusCode, http.StatusNotFound)
	}
}

/*
TestPostAndDeleteRemovesUserFromGet: Given I have created a User and
deleted it when I call the GET method then the user will not be in
the returned data.
*/
func TestPostAndDeleteRemovesUserFromGet(t *testing.T) {
	us, err := NewUserService()
	if err != nil {
		t.Fatal(err.Error())
	}

	/* create user */

	data := map[string]string{
		"country":    "UK",
		"email":      "alice@bob.com",
		"first_name": "Alice",
		"last_name":  "Bob",
		"nickname":   "AB123",
		"password":   "f6b7e19e0d867de6c0391879050e8297165728d89d7c4e9e8839972b356c4d9d",
	}

	post_body, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err.Error())
	}

	post_req, err := http.NewRequest("POST", "/users", bytes.NewReader(post_body))
	if err != nil {
		t.Fatal(err.Error())
	}

	us.ServeHTTP(httptest.NewRecorder(), post_req)

	/* after create get users */

	created_get_req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	created_get_resp := httptest.NewRecorder()
	us.ServeHTTP(created_get_resp, created_get_req)

	created_get_body := []map[string]string{}

	err = json.NewDecoder(created_get_resp.Body).Decode(&created_get_body)
	if err != nil {
		t.Fatal(err.Error())
	}

	/* is there only one user? */

	if len(created_get_body) != 1 {
		t.Fatalf("expected created_get_body == 1 but got %q", len(created_get_body))
	}

	/* delete the user */

	user := created_get_body[0]
	url := fmt.Sprintf("/users/%s", user["id"])

	delete_req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	delete_resp := httptest.NewRecorder()
	us.ServeHTTP(delete_resp, delete_req)

	/* after delete get users */

	deleted_get_req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	deleted_get_resp := httptest.NewRecorder()
	us.ServeHTTP(deleted_get_resp, deleted_get_req)

	deleted_get_body := []map[string]string{}

	err = json.NewDecoder(deleted_get_resp.Body).Decode(&deleted_get_body)
	if err != nil {
		t.Fatal(err.Error())
	}

	/* are there no users? */

	if len(deleted_get_body) != 0 {
		t.Fatalf("expected created_get_body == 1 but got %q", len(deleted_get_body))
	}
}

/*
TestMultiplePostAndDeleteRemovesUserFromGet: Given I have created
multiple users and deleted one of them when I call the GET method
then the deleted user will no longer be in the results.
*/
func TestMultiplePostAndDeleteRemovesUserFromGet(t *testing.T) {
	us, err := NewUserService()
	if err != nil {
		t.Fatal(err.Error())
	}

	/* create users */

	data := map[string](map[string]string){
		"rob": {
			"country":    "Canada/Australia",
			"email":      "rob@bob.com",
			"first_name": "Rob",
			"last_name":  "Pike",
			"nickname":   "rob",
			"password":   "f9c33006f81d188494d2b108a7977ec2710d9fe6c7d33b1b01792eac812d5069",
		},
		"ken": {
			"country":    "USA",
			"email":      "ken@bob.com",
			"first_name": "Ken",
			"last_name":  "Thompson",
			"nickname":   "ken",
			"password":   "b3bb4cd67f11e1f6350a5792c8a0f91c2e7920ab93ccd7e964d97d79ad9f8270",
		},
		"griesemer": {
			"country":    "Switzerland",
			"email":      "robert@bob.com",
			"first_name": "Robert",
			"last_name":  "Griesemer",
			"nickname":   "griesemer",
			"password":   "cb3a8f635e2afa535e2597817cffc0a6aae7698bf63f5d2b3e396de2a6cfb743",
		},
	}

	for _, datum := range data {
		post_body, err := json.Marshal(datum)
		if err != nil {
			t.Fatal(err.Error())
		}

		post_req, err := http.NewRequest("POST", "/users", bytes.NewReader(post_body))
		if err != nil {
			t.Fatal(err.Error())
		}

		us.ServeHTTP(httptest.NewRecorder(), post_req)
	}

	/* after create get users */

	created_get_req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	created_get_resp := httptest.NewRecorder()
	us.ServeHTTP(created_get_resp, created_get_req)

	created_get_body := []map[string]string{}

	err = json.NewDecoder(created_get_resp.Body).Decode(&created_get_body)
	if err != nil {
		t.Fatal(err.Error())
	}

	/* delete the user at index 1 */

	victim := created_get_body[1]["id"]

	url := fmt.Sprintf("/users/%s", victim)

	delete_req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	delete_resp := httptest.NewRecorder()
	us.ServeHTTP(delete_resp, delete_req)

	/* after delete get users */

	deleted_get_req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	deleted_get_resp := httptest.NewRecorder()
	us.ServeHTTP(deleted_get_resp, deleted_get_req)

	deleted_get_body := []map[string]string{}

	err = json.NewDecoder(deleted_get_resp.Body).Decode(&deleted_get_body)
	if err != nil {
		t.Fatal(err.Error())
	}

	/* has the user been removed from the users? */

	for _, user := range deleted_get_body {
		if user["id"] == victim {
			t.Fatalf("expecting user with id %q to be removed", victim)
		}
	}

}

/*
TestAddedStatusUsersGet: Given I have created Users when I call the
GET method then the HTTP status code will be 200 OK.
*/
func TestAddedStatusUsersGet(t *testing.T) {
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

	resp := w.Result()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected error code. Got %d, %d expected.", resp.StatusCode, http.StatusBadRequest)
	}
}

/*
TestEmptyStatusUsersGet: Given I have created no Users when I call
the GET method then the HTTP status code will be 204 No Content.
*/
func TestEmptyStatusUsersGet(t *testing.T) {
	us, err := NewUserService()
	if err != nil {
		t.Fatal(err.Error())
	}

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err.Error())
	}

	w := httptest.NewRecorder()
	us.ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("Unexpected error code. Got %d, %d expected.", resp.StatusCode, http.StatusNoContent)
	}
}

/*
TestMultiUsersGet: Given I have created multiple Users when I call
the GET method then the Users I have created will return with the
following fields populated: created_at, country, email, first_name,
last_name, nickname, password and updated_at with the values I
created them with and the types/formats specified in
./docs/endpoints/users/SCHEMA.md
*/
func TestMultiUsersGet(t *testing.T) {
	us, err := NewUserService()
	if err != nil {
		t.Fatal(err.Error())
	}

	data := map[string](map[string]string){
		"rob": {
			"country":    "Canada/Australia",
			"email":      "rob@bob.com",
			"first_name": "Rob",
			"last_name":  "Pike",
			"nickname":   "rob",
			"password":   "f9c33006f81d188494d2b108a7977ec2710d9fe6c7d33b1b01792eac812d5069",
		},
		"ken": {
			"country":    "USA",
			"email":      "ken@bob.com",
			"first_name": "Ken",
			"last_name":  "Thompson",
			"nickname":   "ken",
			"password":   "b3bb4cd67f11e1f6350a5792c8a0f91c2e7920ab93ccd7e964d97d79ad9f8270",
		},
		"griesemer": {
			"country":    "Switzerland",
			"email":      "robert@bob.com",
			"first_name": "Robert",
			"last_name":  "Griesemer",
			"nickname":   "griesemer",
			"password":   "cb3a8f635e2afa535e2597817cffc0a6aae7698bf63f5d2b3e396de2a6cfb743",
		},
	}

	for _, datum := range data {
		req_body, err := json.Marshal(datum)
		if err != nil {
			t.Fatal(err.Error())
		}

		post, err := http.NewRequest("POST", "/users", bytes.NewReader(req_body))
		if err != nil {
			t.Fatal(err.Error())
		}

		us.ServeHTTP(httptest.NewRecorder(), post)
	}

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

	for _, got_user := range resp_body {
		posted_user := data[got_user["nickname"]]

		for key, expected := range posted_user {
			got := got_user[key]

			if expected != got {
				t.Fatalf("Expected attribute %q to be %q but got %q", key, expected, got)
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
	}
}

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
