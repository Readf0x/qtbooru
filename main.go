package main

import (
	"fmt"
	"log"
	"os"
	"qtbooru/pkg/api"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	req, err := api.NewRequest(
		api.E926,
		&[]string{"limit=5"},
		&[]string{"rating:safe"},
		os.Getenv("API_USER"),
		os.Getenv("API_KEY"),
	)

	posts := api.Process(req)

	for _, p := range *posts {
		fmt.Println(p.CreatedAt)
		fmt.Println(p.Preview.URL)
	}
}

