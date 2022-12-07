package source

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/bitly/go-simplejson"
	"github.com/joho/godotenv"
	"image"
	png2 "image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// 最近n条信息用作训练
var Req string

func AiReply(msg string) (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", err
	}
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

func AiPic(msg string) (image.Image, error) {
	url := "https://api.openai.com/v1/images/generations"
	method := "POST"
	r := fmt.Sprintf(`{"prompt":"%s","n":1,"size": "512x512"}`, msg)
	payload := strings.NewReader(r)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer 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")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Connection", "keep-alive")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(string(body))
	a, err := simplejson.NewJson(body)
	if err != nil {
		return nil, err
	}
	l, err := a.Get("data").GetIndex(0).Get("url").String()
	if err != nil {
		return nil, err
	}
	h, err := http.Get(l)
	if err != nil {
		return nil, err
	}
	p, err := png2.Decode(h.Body)
	if err != nil {
		return nil, err
	}
	return p, nil
}
