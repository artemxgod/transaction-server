package model

import (
	"context"
	"fmt"


	"github.com/segmentio/kafka-go"
)

type User struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Balance   float64 `json:"balance"`
	Reader    *kafka.Reader `json:"-"`
	Writer    *kafka.Writer `json:"-"`
	Writechan chan string `json:"-"`
	Readchan  chan string `json:"-"`
	Started   bool `json:"-"`
}

func (u *User) Read() {
	ctx := context.Background()
	for {
		m, err := u.Reader.ReadMessage(ctx)
		if err != nil {
			break
		}

		fmt.Println("reading a message", string(m.Value))
		u.Readchan <- string(m.Value)
	}
	fmt.Println("broken")
}

func (u *User) Write() {
	ctx := context.Background()
	for {
		msg := <-u.Writechan
		fmt.Println("writing a message", msg)
		u.Writer.WriteMessages(ctx, kafka.Message{
			Value: []byte(msg),
		})
	}
}
