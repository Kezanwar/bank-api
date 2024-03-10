package main

import (
	"log"
)

func main() {

	store, dbErr := NewPostgresDB()

	if dbErr != nil {
		log.Fatal(dbErr)
	}

	initErr := store.Init()

	if initErr != nil {
		log.Fatal(initErr)
	}

	server := NewApiServer(":8000", store)
	server.Run()
}
