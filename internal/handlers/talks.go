package handlers

import (
	"context"
	"io"
	"log"
	"os"

	"github.com/sashabaranov/go-openai"
	"golang.org/x/net/websocket"
)

type Talks struct {
	conn    *websocket.Conn
	speaker *openai.Client
}

func NewTalks() *Talks {
	return &Talks{
		speaker: openai.NewClient(os.Getenv("OPENAI_API_KEY")),
	}
}

func (t *Talks) Handle(conn *websocket.Conn) {
	t.conn = conn
	t.Read(conn)
}

func (t *Talks) Read(conn *websocket.Conn) {
	defer conn.Close()

	msg := make([]byte, 512)
	for {
		n, err := conn.Read(msg)
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		t.Write(conn, msg[:n])
	}
}

func (t *Talks) Write(conn *websocket.Conn, content []byte) {
	resp, err := t.speaker.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: string(content),
				},
			},
		},
	)

	if err != nil {
		log.Println(err)
		return
	}

	if _, err := conn.Write([]byte(resp.Choices[0].Message.Content)); err != nil {
		log.Println(err)
	}
}
