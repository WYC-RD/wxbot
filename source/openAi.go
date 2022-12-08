package source

import (
	"context"
	"fmt"
	"github.com/PullRequestInc/go-gpt3"
	"github.com/bitly/go-simplejson"
	"github.com/eatmoreapple/openwechat"
	"github.com/joho/godotenv"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// 最近n条信息用作训练
var Prompt = map[string]string{}

func OpHandle(message *openwechat.Message) error {
	if !message.IsText() {
		message.ReplyText("?")
		return nil
	}
	s, err := message.Sender()
	if err != nil {
		return err
	}
	//fmt.Println(s.NickName)
	u, _ := message.Bot.GetCurrentUser()
	msg := strings.Replace(message.Content, "@"+u.NickName+" ", "", 1)

	//画画接口调用
	if strings.Contains(msg, "我看到") {
		if err := picHandle(msg, message); err != nil {
			message.ReplyText(err.Error())
			return err
		}
		return nil
	}

	//调用问答接口
	reply, err := AiReply(msg, s.NickName)
	if err != nil {
		if err := secReq(message, msg, s.NickName); err != nil {
			message.ReplyText(err.Error())
			return err
		}
	}
	message.ReplyText(reply)
	//fmt.Println(Prompt[s.NickName])
	return nil
}

func secReq(message *openwechat.Message, msg string, nick string) error {
	slice := strings.Split(Prompt[nick], " ")
	Prompt[nick] = strings.Join(slice[2:], " ")
	reply, err := AiReply(msg, nick)
	message.ReplyText(reply)
	if err != nil {
		return err
	}
	return nil
}

func AiReply(msg string, nick string) (string, error) {
	if err := godotenv.Load(); err != nil {
		return "", err
	}
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Println("Missing API KEY")
	}
	ctx := context.Background()
	client := gpt3.NewClient(apiKey)
	pp := Prompt[nick]
	ask := fmt.Sprintf(" Human:%s AI:", msg)
	pp += ask
	for len(pp) > 4000 {
		slice := strings.Split(pp, " ")
		slice = slice[2:]
		pp = strings.Join(slice, " ")
	}
	x := gpt3.CompletionRequest{
		//Model:  []string{"text-davinci-003"},
		Prompt:           []string{"The following is a conversation with an AI assistant." + pp},
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
	reply = strings.TrimPrefix(reply, "\n\n")
	//fmt.Println(reply)
	pp += reply
	Prompt[nick] = pp
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
	//fmt.Println(string(body))
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
	p, err := png.Decode(h.Body)
	if err != nil {
		return nil, err
	}
	return p, nil
}
func picHandle(msg string, message *openwechat.Message) error {
	picname := strings.Replace(msg, "我看到", "", -1)
	pic, err := AiPic(picname)
	if err != nil {
		return err
	}
	fname := fmt.Sprintf("./wxbot-pic-log/openai/%s.png", message.MsgId)
	f, err := os.Create(fname)
	defer f.Close()
	if err != nil {
		return err
	}
	if err := png.Encode(f, pic); err != nil {
		return err
	}
	o, err := os.Open(fname)
	message.ReplyImage(o)
	return nil
}
