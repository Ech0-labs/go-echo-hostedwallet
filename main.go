package main

import (
	"log"

	"github.com/Ech0-labs/go-echo-prototype/api"
	"github.com/joho/godotenv"
)

const message = "Test keynote 14/12"

func Handle(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	Handle(godotenv.Load())

	_, err := api.InitFromEnv()
	Handle(err)

	select {}
}
