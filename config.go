package main

import (
	_ "embed"
	"encoding/json"
	"log"
)

//go:embed config.json
var configJSON string

type Config struct {
	GeminiApiKey string `json:"GEMINI-API-KEY"`
	JiraDomain   string `json:"JIRA-DOMAIN"`
	JiraEmail    string `json:"JIRA-EMAIL"`
	JiraApiKey   string `json:"JIRA-API-KEY"`
}

func GetConfig() Config {
	var config Config
	err := json.Unmarshal([]byte(configJSON), &config)
	if err != nil {
		log.Fatal(err)
	}
	return config
}
