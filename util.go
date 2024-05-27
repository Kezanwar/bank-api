package main

import (
	"encoding/json"
	"log"
)

func print_map(m any) {
	b, err := json_stringify(m)
	if err != nil {
		log.Fatal(err)
	}
	println(b)
}

func json_stringify(m any) (string, error) {
	b, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
