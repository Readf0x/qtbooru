//go:build !debug

package main

import "qtbooru/config"

func loadEnv() *config.ApiConfig {
	return config.LoadConfig()
}
