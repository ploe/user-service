package http

import "net/http"

type UserService struct{}

func NewUserService() (*UserService, error) {
	return &UserService{}, nil
}

func (us *UserService) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
