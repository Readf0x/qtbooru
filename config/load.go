package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type ApiConfig struct {
	Username string `json:"username"`
	Key      string `json:"key"`
}

func LoadConfig() *ApiConfig {
	c := os.Getenv("XDG_CONFIG_HOME")
	if c == "" {
		c = os.Getenv("HOME") + "/.config"
	}
	paths := []string{ c + "/qtbooru.json" }
	conf := &ApiConfig{}
	for _, p := range paths {
		file, err := os.Open(p)
		if err == nil {
			b, err := io.ReadAll(file)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			json.Unmarshal(b, conf)
			return conf
		}
	}
	return nil
}

