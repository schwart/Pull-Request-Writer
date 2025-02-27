package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"log"
)

type PullRequestResponse struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func CallGemini(gitPatch string, modelName string, geminiApiKey string) *PullRequestResponse {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(geminiApiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)
	model.SystemInstruction = genai.NewUserContent(genai.Text("You are a developer who is writing a pull request on Github."))
	model.ResponseMIMEType = "application/json"
	model.GenerationConfig.ResponseSchema = &genai.Schema{
		Type: genai.TypeObject,
		Properties: map[string]*genai.Schema{
			"title": {
				Type: genai.TypeString,
			},
			"body": {
				Type: genai.TypeString,
			},
		},
		Required: []string{
			"title",
			"body",
		},
	}
	resp, err := model.GenerateContent(ctx, genai.Text(gitPatch))
	if err != nil {
		log.Fatal(err)
	}

	return formatResponse(resp)
}

func formatResponse(resp *genai.GenerateContentResponse) *PullRequestResponse {
	responseObject := PullRequestResponse{
		Title: "",
		Body:  "",
	}
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					currentPart := PullRequestResponse{}
					if err := json.Unmarshal([]byte(txt), &currentPart); err != nil {
						log.Fatal(err)
					}

					if currentPart.Title != "" {
						responseObject.Title = currentPart.Title
					}

					responseObject.Body += currentPart.Body
				}
			}
		}
	}
	if responseObject.Title == "" || responseObject.Body == "" {
		return nil
	}
	return &responseObject
}
