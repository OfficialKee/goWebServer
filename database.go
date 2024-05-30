package main

import (
    "encoding/json"
    "fmt"
    "log"
)
type chirp struct {
	Id int `json:"id"`
	Body string `json:"body"`
}

type dataBase struct {
	Chirps  map[int]chirp `json:"chirps"`
}

Db := new dataBase{
	Chirps: make(map[int]chirp),
}