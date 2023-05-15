package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func readConfigFromENV() map[string]string {
	var envs = make(map[string]string)

	envs["MASTODON_SERVER"] = os.Getenv("MASTODON_SERVER")
	envs["MASTODON_CLIENT_ID"] = os.Getenv("MASTODON_CLIENT_ID")
	envs["MASTODON_CLIENT_SECRET"] = os.Getenv("MASTODON_CLIENT_SECRET")
	envs["MASTODON_APP_TOKEN"] = os.Getenv("MASTODON_APP_TOKEN")
	envs["RAINDROP_SERVER"] = os.Getenv("RAINDROP_SERVER")
	envs["RAINDROP_CLIENT_ID"] = os.Getenv("RAINDROP_CLIENT_ID")
	envs["RAINDROP_CLIENT_SECRET"] = os.Getenv("RAINDROP_CLIENT_SECRET")
	envs["RAINDROP_APP_TOKEN"] = os.Getenv("RAINDROP_APP_TOKEN")
	return envs
}

func GetConfig() (map[string]string, error) {

	if os.Getenv("APP_ENV") != "production" {
		envs, error := godotenv.Read(".env")
		if error != nil {
			fmt.Println("Error loading .env file")
		}
		return envs, error
	} else {
		envs := readConfigFromENV()
		return envs, nil
	}
}
