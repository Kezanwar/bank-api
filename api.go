package main

import "net/http"

type APIServer struct {
	listernAddr string
}

func NewApiServer(listenAddr string) *APIServer {
	return &APIServer{
		listernAddr: listenAddr,
	}
}

func (server *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (server *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (server *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (server *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (server *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	return nil
}
