package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type datetime struct {
	tm time.Time
}

type UserService struct {
	callback chan func()
	hc       *healthchecker
	mux      *http.ServeMux
	users    map[string]*user
}

type user struct {
	CreatedAt datetime `json:"created_at"`
	Country   string   `json:"country"`
	Email     string   `json:"email"`
	FirstName string   `json:"first_name"`
	ID        string   `json:"id"`
	LastName  string   `json:"last_name"`
	Nickname  string   `json:"nickname"`
	Password  string   `json:"password"`
	UpdatedAt datetime `json:"updated_at"`
}

const DtLayout = "2006-01-02T15:04.05Z"

func NewDatetime() datetime {
	return datetime{
		tm: time.Now(),
	}
}

/* Marshal the datetime in custom datetime format */
func (dt *datetime) MarshalJSON() ([]byte, error) {
	b := fmt.Sprintf(`"%s"`, dt.tm.Format(DtLayout))

	return []byte(b), nil
}

/* Marshal the datetime in custom datetime format */
func (dt *datetime) UnmarshalJSON(b []byte) error {
	tm, err := time.Parse(DtLayout, string(b))
	if err != nil {
		return err
	}

	dt.tm = tm

	return nil
}

/* Create a new UserService. */
func NewUserService() (*UserService, error) {
	us := &UserService{
		callback: make(chan func()),
		hc:       newHealthchecker(),
		mux:      http.NewServeMux(),
		users:    make(map[string]*user),
	}

	//	us.mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNotFound) }))
	us.mux.Handle("/healthcheck", us.hc)
	us.mux.Handle("/users", us)

	go func() {
		for {
			(<-us.callback)()
		}
	}()

	return us, nil
}

/* Add a new user to the in-memory storage mechanism. */
func (us *UserService) addUser(user *user) {
	us.callback <- func() {
		us.users[user.ID] = user
	}
}

/* Delete a user from the in-memory storage mechanism. */
func (us *UserService) deleteUser(id string) bool {
	ch := make(chan bool)

	us.callback <- func() {
		_, ok := us.users[id]

		if ok {
			delete(us.users, id)
		}

		ch <- ok
	}

	return <-ch
}

/*
Get a filtered list of the Users from the in-memory storage
mechanism.
*/
func (us *UserService) getUsers(filters map[string]string) []*user {
	ch := make(chan []*user)

	us.callback <- func() {
		users := []*user{}
		for _, user := range us.users {
			current := map[string]string{
				"country":    user.Country,
				"email":      user.Email,
				"first_name": user.FirstName,
				"last_name":  user.LastName,
				"nickname":   user.Nickname,
			}

			add := true
			for key, filter := range filters {
				if current[key] != filter {
					add = false
					break
				}
			}

			if !add {
				continue
			}

			users = append(users, user)
		}

		ch <- users
	}

	return <-ch
}

/* Modify a user from the in-memory storage mechanism. */
func (us *UserService) modifyUser(id string, data map[string]string) bool {
	ch := make(chan bool)

	us.callback <- func() {
		user, ok := us.users[id]
		if !ok {
			ch <- false
			return
		}

		attributes := []struct {
			key    string
			target *string
		}{
			{"country", &(user.Country)},
			{"email", &(user.Email)},
			{"first_name", &(user.FirstName)},
			{"last_name", &(user.LastName)},
			{"nickname", &(user.Nickname)},
			{"password", &(user.Password)},
		}

		modified := false
		for _, attribute := range attributes {
			value, ok := data[attribute.key]

			if !ok {
				continue
			}

			*(attribute.target) = value
			modified = true
		}

		if modified {
			user.UpdatedAt.tm = time.Now()
		}

		ch <- true
	}

	return <-ch
}

/* Serves HTTP on the requested addr */
func (us *UserService) ListenAndServe(addr string) error {
	return http.ListenAndServe(addr, us.mux)
}

/* Handler method ServeHTTP for UserService. */
func (us *UserService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodDelete:
		us.delete(w, r)
	case http.MethodGet:
		us.get(w, r)
	case http.MethodPatch:
		us.patch(w, r)
	case http.MethodPost:
		us.post(w, r)
	}
}

func (us *UserService) delete(w http.ResponseWriter, r *http.Request) {
	sender := r.RemoteAddr

	id := filepath.Base(r.URL.Path)

	log.Printf("[%s] DELETE /users: attempting to delete user %q", sender, id)

	err := uuid.Validate(id)
	if err != nil {
		log.Printf("[%s] DELETE /users: %q is not a valid user id", sender, id)

		w.WriteHeader(http.StatusNotFound)
		return
	}

	ok := us.deleteUser(id)
	if !ok {
		log.Printf("[%s] DELETE /users: %q is not a user", sender, id)

		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Printf("[%s] DELETE /users: deleted %q", sender, id)

	w.WriteHeader(http.StatusNoContent)
}

func (us *UserService) get(w http.ResponseWriter, r *http.Request) {
	sender := r.RemoteAddr

	url := r.URL.Query()

	limit := 0
	for _, value := range url["limit"] {
		i, err := strconv.Atoi(value)
		if err != nil {
			continue
		}

		limit = i
		break
	}

	page := 0
	for _, value := range url["page"] {
		i, err := strconv.Atoi(value)
		if err != nil {
			continue
		}

		page = i
		break
	}

	filters := map[string]string{}
	for _, key := range []string{"country", "email", "first_name", "last_name", "nickname"} {
		query, ok := url[key]

		if !ok {
			continue
		}

		for _, value := range query {
			filters[key] = value
		}
	}

	log.Printf("[%s] GET /users: attempting to get users", sender)

	users := us.getUsers(filters)

	if limit != 0 {
		start := (page * limit)
		end := start + limit
		users = users[start:end]
	}

	body, err := json.Marshal(users)
	if err != nil {
		log.Printf("[%s] GET /users: unable to marshal users %q", sender, err.Error())

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	if len(users) == 0 {
		w.WriteHeader(http.StatusNoContent)
	}

	log.Printf("[%s] GET /users: got users", sender)
	w.Write(body)
}

func (us *UserService) patch(w http.ResponseWriter, r *http.Request) {
	sender := r.RemoteAddr

	id := filepath.Base(r.URL.Path)

	log.Printf("[%s] PATCH /users: attempting to patch user %q", sender, id)

	err := uuid.Validate(id)
	if err != nil {
		log.Printf("[%s] PATCH /users: %q is not a valid user id", sender, id)

		w.WriteHeader(http.StatusNotFound)
		return
	}

	data := map[string]string{}

	if r.Body == nil {
		log.Printf("[%s] PATCH /users: no body sent with %q", sender, id)

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf("[%s] PATCH /users: unable to decode JSON on %q: %s", sender, id, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ok := us.modifyUser(id, data)
	if !ok {
		log.Printf("[%s] PATCH /users: %q is not a user", sender, id)

		w.WriteHeader(http.StatusNotFound)
		return
	}

	log.Printf("[%s] PATCH /users: patched %q", sender, id)

	w.WriteHeader(http.StatusNoContent)
}

func (us *UserService) post(w http.ResponseWriter, r *http.Request) {
	id := uuid.NewString()
	sender := r.RemoteAddr

	log.Printf("[%s] POST /users: attempting to add %q", sender, id)

	data := map[string]string{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Printf("[%s] POST /users: unable to decode JSON on %q: %s", sender, id, err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expected := []string{"country", "email", "first_name", "last_name", "nickname", "password"}
	for _, key := range expected {
		_, ok := data[key]
		if !ok {
			log.Printf("[%s] POST /users: unable to add as attribute %q was missing on %q", sender, id, key)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	now := time.Now()
	user := user{
		CreatedAt: datetime{tm: now},
		Country:   data["country"],
		Email:     data["email"],
		FirstName: data["first_name"],
		ID:        id,
		LastName:  data["last_name"],
		Nickname:  data["nickname"],
		Password:  data["password"],
		UpdatedAt: datetime{tm: now},
	}

	us.addUser(&user)
	log.Printf("[%s] POST /users: added %q", sender, id)

	w.WriteHeader(http.StatusCreated)
}
