package main

import "math/rand"

type Account struct {
	ID        int
	FirstName string
	LastName  string
	Number    int64
	Balance   uint64
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		ID:        rand.Intn(1000),
		FirstName: firstName,
		LastName:  lastName,
		Number:    int64(rand.Intn(10000000)),
	}
}
