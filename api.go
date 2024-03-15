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
	Error string `json:"error"`
}

type EmptySuccessResponse struct {
	Message string `json:"message"`
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
	router.HandleFunc("/account/{id}", makeHTTPHandler(s.handleAccountByID))
	router.HandleFunc("/transfer", makeHTTPHandler(s.handleTransfer))

	log.Println("Bank API running on port", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)

	return nil
}

// * routes to all REST handlers for /account
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}

	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	return fmt.Errorf("method not allow %s", r.Method)
}

// * routes to all REST handlers for /account/{id}
func (s *APIServer) handleAccountByID(w http.ResponseWriter, r *http.Request) error {

	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allow %s", r.Method)
}

// * routes to all REST handlers for /transfer
func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "POST" {
		return s.handleTransferFunds(w, r)
	}

	return fmt.Errorf("method not allow %s", r.Method)
}

// * handles GET /account and /account/:id
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {

	idStr := mux.Vars(r)["id"]

	if len(idStr) > 0 {

		id, err := getID(r)

		// if ID cant cast to an int
		if err != nil {
			return err
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

// * handles POST /account
func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	defer r.Body.Close()

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)

	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

// * handles DELETE /account
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	id, err := getID(r)

	if err != nil {
		return err
	}

	deleteErr := s.store.DeleteAccount(id)

	if deleteErr != nil {
		return fmt.Errorf("unable to delete account with ID %d", id)
	}

	return WriteJSON(w, http.StatusOK, &EmptySuccessResponse{Message: "Account deleted succesfully"})
}

func (s *APIServer) handleTransferFunds(w http.ResponseWriter, r *http.Request) error {
	transfer := &TransferRequest{}

	err := json.NewDecoder(r.Body).Decode(transfer)

	if err != nil {
		return err
	}

	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transfer)
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]

	if len(idStr) == 0 {
		return 0, fmt.Errorf("no ID given")
	}

	id, err := strconv.Atoi(idStr)

	// if ID cant cast to an int
	if err != nil {
		return 0, fmt.Errorf("invalid ID given: %s", idStr)
	}

	return id, nil
}
