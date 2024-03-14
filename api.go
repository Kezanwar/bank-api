package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      *PostGresDB
}

func NewApiServer(listenAddr string, store *PostGresDB) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

type ApiHandler func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandler(f ApiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle the error
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandler(s.handleAccount))

	router.HandleFunc("/account/{id}", makeHTTPHandler(s.handleGetAccount))

	log.Println("Bank API running on port", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)

	return nil
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}

	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allow %s", r.Method)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	idStr := mux.Vars(r)["id"]

	if len(idStr) > 0 {

		id, err := strconv.Atoi(idStr)

		// if ID cant cast to an int
		if err != nil {
			return fmt.Errorf("invalid ID given: %s", idStr)
		}

		account, err := s.store.GetAccountByID(id)

		// if Account doesn't exist
		if err != nil {
			return fmt.Errorf("no account found for ID given: %s", idStr)
		}

		WriteJSON(w, http.StatusOK, account)
	} else {
		accounts, err := s.store.GetAccounts()

		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, accounts)
	}

	return nil
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
