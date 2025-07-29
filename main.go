package main

import (
	"booru/pkg/api"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

const Agent = "QtBooru/indev_v0 (created by readf0x)"

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}
	rq, err := api.NewRequest(
		api.E926,
		&[]string{
			"limit=1",
		},
		os.Getenv("API_USER"),
		os.Getenv("API_KEY"),
		Agent,
	)

	resp, err := client.Do(rq)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Body: %s\n", body)
}

