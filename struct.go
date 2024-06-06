package main

import (
	"sync"
)

// Employee represents an employee's data
type Employee struct {
	ID       int
	Name     string
	Position string
	Salary   float64
}

// Store represents the in-memory employee store
type Store struct {
	sync.Mutex
	employees map[int]Employee
	nextID    int
}

// NewStore initializes a new employee store
func NewStore() *Store {
	return &Store{
		employees: make(map[int]Employee),
		nextID:    1,
	}
}

type JsonWrappedContent struct {
	StatusCode  int         `json:"statusCode"`
	LastUpdated string      `json:"last_updated,omitempty"`
	Content     interface{} `json:"content"`
}
