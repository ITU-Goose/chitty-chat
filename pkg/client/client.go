package client

import (
	"github.com/google/uuid"
	pb "github.com/goose-alt/chitty-chat/api/v1/pb/chat"
	"github.com/goose-alt/chitty-chat/internal/time"
)

type Client struct {
	Uuid string
	Name string

	Timestamp *time.VectorTimestamp
	Stream    *pb.Chat_ChatClient
}

func NewClient(name string, stream *pb.Chat_ChatClient) *Client {
	id := uuid.New().String()

	client := Client{
		Uuid: id,
		Name: name,

		Timestamp: time.CreateVectorTimestamp(id),
		Stream:    stream,
	}

	return &client
}

func (c *Client) SendMessage(message string) error {
	c.Timestamp.Increment()

	msg := &pb.Message{
		Timestamp: &pb.Lamport{Clients: c.Timestamp.GetVectorTime()},
		Info: &pb.ClientInfo{
			Uuid: c.Uuid,
			Name: c.Name,
		},
		Content: message,
	}

	return (*c.Stream).Send(msg)
}
