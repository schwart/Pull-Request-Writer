package main

import (
	"context"
	"log"
	"strings"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func CallGemini(gitPatch string) string {
	ctx := context.Background()
	apiKey := "AIzaSyCxX5Auy-LfDTeYuKWza6xmKvzfMtd8PxU"
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-2.0-flash")
	model.SystemInstruction = genai.NewUserContent(genai.Text("You are a developer who is writing a pull request on Github."))
	resp, err := model.GenerateContent(ctx, genai.Text(gitPatch))
	if err != nil {
		log.Fatal(err)
	}

	return formatResponse(resp)
}

func formatResponse(resp *genai.GenerateContentResponse) string {
	var sb strings.Builder

	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				if txt, ok := part.(genai.Text); ok {
					sb.WriteString(string(txt))
				}
			}
		}
	}
	return sb.String()
}
