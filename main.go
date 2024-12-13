package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"

	openai "github.com/sashabaranov/go-openai"
)

const prompt = `
Please create me a well formed git commit message and description
for the following code changes.

Be sure to prefix with a type of change that this diff implies, such as a feat, chore, fix etc.

Please be both concise but make sure to cover all the changes in this code.

Git diff begins in the next message
  `

func main() {
	log.SetFlags(0)

	cmd := exec.Command("git", "diff", "--staged")
	output, err := cmd.Output()

	if err != nil {
		log.Fatalln("failed to run command:", err)
	}

	client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: prompt,
			},
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: string(output),
			},
		},
	}

	resp, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		fmt.Printf("Completion error: %v\n", err)
		return
	}

	msg := resp.Choices[0].Message.Content

	command := exec.Command("git", "commit", "-em", msg)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Run()
}
