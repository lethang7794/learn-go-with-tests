package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func TestUserServer_RegisterUser(t *testing.T) {
	t.Run("", func(t *testing.T) {
		user := User{
			Username: "harry",
			Password: "12345",
		}
		expectedUserId := "whatever"
		service := &MockUserService{
			RegisterFunc: func(user User) (insertedID string, err error) {
				return expectedUserId, nil
			},
		}
		srv := NewUserServer(service)
		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/", userToJson(user))

		srv.RegisterUser(response, request)

		assertStatus(t, response.Code, http.StatusCreated)
		if response.Body.String() != expectedUserId {
			t.Errorf("expected got %#v, want %#v", response.Body.String(), expectedUserId)
		}
		if len(service.registeredUsers) != 1 {
			t.Fatalf("expected 1 user created: got %v", service.registeredUsers)
		}
		if !reflect.DeepEqual(service.registeredUsers[0], user) {
			t.Errorf("got %v, want %v", service.registeredUsers[0], user)
		}
	})

	t.Run("returns 400 Bad Request if request body is not a valid JSON", func(t *testing.T) {
		body := strings.NewReader("Not a valid JSON")
		service := &MockUserService{
			RegisterFunc: func(user User) (insertedID string, err error) {
				return "", nil
			},
		}
		srv := NewUserServer(service)
		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/", body)

		srv.RegisterUser(response, request)

		assertStatus(t, response.Code, http.StatusBadRequest)
	})
	t.Run("returns 500 Internal if user service is down", func(t *testing.T) {
		user := User{}
		service := &MockUserService{
			RegisterFunc: func(user User) (insertedID string, err error) {
				return "", errors.New("user service is down")
			},
		}
		srv := NewUserServer(service)
		response := httptest.NewRecorder()
		request := httptest.NewRequest("GET", "/", userToJson(user))

		srv.RegisterUser(response, request)

		assertStatus(t, response.Code, http.StatusInternalServerError)
	})
}

func assertStatus(t *testing.T, got int, want int) {
	if got != want {
		t.Errorf("got %#v, want %#v", got, want)
	}
}

func userToJson(user User) io.Reader {
	bs, _ := json.Marshal(user)
	return bytes.NewReader(bs)
}
