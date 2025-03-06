package main

import (
	"encoding/json"
	"errors"
	"net/http"
)

type api struct {
	addr string
}

var users = []User{}

func (a *api) getUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	//below step put the users into ResponseWriter in json format
	err := json.NewEncoder(w).Encode(users)
	if err != nil {
		//new introduction to http.Error
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
}

func (a *api) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload User
	//decoding the payload from r
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	u := User{
		FirstName: payload.FirstName,
		LastName:  payload.LastName,
	}
	//appending newly created user to user slice after checks
	err = insertUser(u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
}

func insertUser(u User) error {
	if u.FirstName == "" || u.LastName == "" {
		return errors.New("invalid user")
	}
	for _, user := range users {
		if user.FirstName == u.FirstName && user.LastName == u.LastName {
			return errors.New("duplicate user")
		}
	}
	users = append(users, u)
	return nil
}
