package server

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/goose-alt/chitty-chat/internal"
	pb "github.com/goose-alt/chitty-chat/api/v1/pb/chat"
)

type chatServer struct {
	pb.UnimplementedChatServer

	// List of clients, mapped by their generated id
	clients map[string]internal.Client
}

func NewChatServer() chatServer {
	return chatServer {
		clients: make(map[string]internal.Client),
	}
}

// TODO: Delete this
func fakeLamport() *pb.Lamport {
	lamport := pb.Lamport {
		Clients: make(map[string]int32),
	}

	return &lamport
}

func (s *chatServer) Register(ctx context.Context, in *pb.RegisterMessage) (*pb.RegisterResponse, error) {
	// TODO: Do something with lamport

	// Generate a new uuid for the client
	id := uuid.New().String()

	// Register client
	s.clients[id] = internal.Client {
		Uuid: id,
		Name: in.Name,
	}

	// Return response
	return &pb.RegisterResponse {
		Timestamp: fakeLamport(), // TODO: Add proper lamport response
		Info: &pb.ClientInfo {
			Uuid: id,
			Name: in.Name,
		},
	}, nil
}

// TODO: Implement this
func (s *chatServer) Chat(stream pb.Chat_ChatServer) error {
	return errors.New("Unimplemented")
}
