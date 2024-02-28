package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	userService := &MockUserService{}
	userServer := NewUserServer(userService)
	server := http.Server{Addr: ":8080", Handler: http.HandlerFunc(userServer.RegisterUser)}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

type User struct {
	Username string
	Password string
}

type UserServer struct {
	service UserService
}

func NewUserServer(service UserService) *UserServer {
	return &UserServer{service: service}
}

func (u *UserServer) RegisterUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// request parsing and validation
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("could not decode user payload: %v", err), http.StatusBadRequest)
		return
	}

	// call a service thing to take care of the hard work
	insertedID, err := u.service.Register(newUser)

	// depending on what we get back, respond accordingly
	if err != nil {
		//todo: handle different kinds of errors differently
		http.Error(w, fmt.Sprintf("problem registering new user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, insertedID)
}
