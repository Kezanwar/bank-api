package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/joho/godotenv/autoload"
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
	Message string `json:"message"`
}

type EmptySuccessResponse struct {
	Message string `json:"message"`
}

func makeHTTPHandler(f ApiHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle the error
			WriteJSON(w, http.StatusBadRequest, ApiError{Message: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func getUserFromHeader(r *http.Request) (*Account, error) {
	user := r.Header.Get("x-server-user")

	if len(user) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	userInterface := Account{}

	err := json.Unmarshal([]byte(user), &userInterface)

	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return &userInterface, nil

}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandler(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPHandler(s.handleAccountByID))
	router.HandleFunc("/transfer", makeHTTPHandler(s.handleTransfer))

	router.Use(loggingMiddleware)
	router.Use(makeAuthMiddleware(s))

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
		// user, err := getUserFromHeader(r)

		// print_map(user)

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

	newAcc := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)

	createdAcc, err := s.store.CreateAccount(newAcc)

	if err != nil {
		return fmt.Errorf("unable to create an account")
	}

	tokenString, err := createJWT(createdAcc)

	log.Println(tokenString)

	if err != nil {
		return WriteJSON(w, http.StatusInternalServerError, &ApiError{Message: "a server error occured"})
	}

	return WriteJSON(w, http.StatusOK, createdAcc)
}

// * handles DELETE /account
func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {

	id, err := getID(r)

	if err != nil {
		return err
	}

	deleteErr := s.store.DeleteAccountByID(id)

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
