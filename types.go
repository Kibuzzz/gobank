package main

import (
	"math/rand"
	"time"
)

type CreateAccountRequest struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Account struct {
	ID         int       `json:"id"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	Number     int64     `json:"number"`
	Balance    int64     `json:"balance"`
	Created_at time.Time `json:"created_at"`
}

func NewAccount(firstName, lastName string) *Account {
	return &Account{
		FirstName:  firstName,
		LastName:   lastName,
		Number:     int64(rand.Intn(100000)),
		Balance:    0,
		Created_at: time.Now().UTC(),
	}
}