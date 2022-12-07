package source

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

// 最近n条信息用作训练
var Req string

func AiReply(msg string) (string, error) {
	godotenv.Load()
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API KEY")
	}
	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	ask := fmt.Sprintf(" Human:%s AI:", msg)
	Req += ask
	for len(Req) > 4000 {
		slice := strings.Split(Req, " ")
		slice = slice[2:]
		Req = strings.Join(slice, " ")
	}
	x := gpt3.CompletionRequest{
		//Model:  []string{"text-davinci-003"},
		Prompt:           []string{"The following is a conversation with an AI assistant." + Req},
		MaxTokens:        gpt3.IntPtr(800),
		Stop:             []string{" Human:", " AI:"},
		Echo:             false,
		Temperature:      gpt3.Float32Ptr(0.9),
		TopP:             gpt3.Float32Ptr(1),
		FrequencyPenalty: 0,
		PresencePenalty:  0.6,
	}

	var reply string
	err := client.CompletionStreamWithEngine(ctx, "text-davinci-003", x, func(response *gpt3.CompletionResponse) {
		reply += response.Choices[0].Text
	})
	if err != nil {
		log.Println(err)
		return "", err
	}
	//fmt.Println(reply)
	Req += reply
	return reply, nil
}
