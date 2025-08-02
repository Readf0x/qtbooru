package main

import (
	// "fmt"
	"log"
	"qtbooru/ui"
	// "os"
	// "os/exec"
	// "qtbooru/pkg/api"
	// "strings"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// tags := os.Args[1:]
	// req, err := api.NewRequest(
	// 	api.E926,
	// 	&[]string{"limit=1"},
	// 	&tags,
	// 	os.Getenv("API_USER"),
	// 	os.Getenv("API_KEY"),
	// )
	//
	// posts := *api.Process(req)
	//
	// url := strings.Replace(posts[0].Preview.URL, "localhost", "loki2", 1)
	// fmt.Println(url)
	// cmd := exec.Command("kitten", "icat", url)
	// cmd.Env = os.Environ()
	// err = cmd.Run()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	ui.Spawn()
}

