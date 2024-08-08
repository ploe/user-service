package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type datetime struct {
	tm time.Time
}

type UserService struct {
	callback chan func()
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
		mux:      http.NewServeMux(),
		users:    make(map[string]*user),
	}

	us.mux.Handle("/users", us)

	go func() {
		for {
			(<-us.callback)()
		}
	}()

	return us, nil
}

/* Add a new user to the in-memory storage mechanism. */
func (us *UserService) AddUser(user *user) {
	us.callback <- func() {
		us.users[user.ID] = user
	}
}

func (us *UserService) GetUsers() []*user {
	ch := make(chan []*user)

	us.callback <- func() {
		users := []*user{}
		for _, value := range us.users {
			users = append(users, value)
		}

		ch <- users
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
	case http.MethodPost:
		us.post(w, r)
	case http.MethodGet:
		us.get(w, r)
	}
}

func (us *UserService) get(w http.ResponseWriter, r *http.Request) {
	sender := r.RemoteAddr

	log.Printf("[%s] GET /users: attempting to get users", sender)

	users := us.GetUsers()

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

	us.AddUser(&user)
	log.Printf("[%s] POST /users: added %q", sender, id)

	w.WriteHeader(http.StatusCreated)
}
