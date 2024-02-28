package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIserver struct {
	listenAddress string
	store         Storage
}

type ApiError struct {
	Error string
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func NewAPIServer(address string, store Storage) *APIserver {
	return &APIserver{
		listenAddress: address,
		store:         store,
	}
}

func (s *APIserver) Run() error {
	router := mux.NewRouter()
	router.Handle("/account", makeHTTPHandleFunc(s.handleAccount))
	log.Println("Json api server running on port: ", s.listenAddress)
	return http.ListenAndServe(s.listenAddress, router)
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func (s *APIserver) handleAccount(w http.ResponseWriter, r *http.Request) error {
	rt := r.Method
	switch rt {
	case "GET":
		return s.handleGetAccount(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	}
	return fmt.Errorf("method not allowed %s", rt)
}

func (s *APIserver) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	account := NewAccount("Pablo", "Technik")
	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIserver) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, accounts)
}

func (s *APIserver) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := &CreateAccountRequest{}
	if err := json.NewDecoder(r.Body).Decode(&createAccountReq); err != nil {
		return err
	}
	new_account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if err := s.store.CreateAccount(new_account); err != nil {
		return err
	}
	fmt.Printf("%+v", new_account)
	return WriteJSON(w, http.StatusOK, new_account)
}
func (s *APIserver) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIserver) handleTransaction(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}
