package main

import (
	"context"
	"encoding/json"
	"log"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type PullRequestResponse struct {
	Title string `json:"title"`
	Body string `json:"body"`
}

func CallGemini(gitPatch string, modelName string) *PullRequestResponse {
	ctx := context.Background()
	apiKey := "AIzaSyCxX5Auy-LfDTeYuKWza6xmKvzfMtd8PxU"
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
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
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					responseObject := PullRequestResponse{}
					if err := json.Unmarshal([]byte(txt), &responseObject); err != nil {
						log.Fatal(err)
					}
					return &responseObject
				}
			}
		}
	}
	return nil
}
