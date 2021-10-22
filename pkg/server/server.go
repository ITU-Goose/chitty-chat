package server

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"

	pb "github.com/goose-alt/chitty-chat/api/v1/pb/chat"
	"github.com/goose-alt/chitty-chat/internal"
)

type chatServer struct {
	pb.UnimplementedChatServer

	// List of clients, mapped by their generated id
	clients map[string]internal.Client
	lock sync.Mutex
}

func NewChatServer() chatServer {
	return chatServer {
		clients: make(map[string]internal.Client),
	}
}

func (s *chatServer) addClient(id string, stream pb.Chat_ChatServer) {

	s.lock.Lock()
	defer s.lock.Unlock()

	s.clients[id] = internal.Client{Uuid: id, Name: "Client-"+id, Chat: stream}
}

func (s *chatServer) removeClient(id string) {

	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.clients, id)
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

/*
Is a stream to send chat messages. This is bidirectional.

The implementation is inspired by: https://github.com/castaneai/grpc-broadcast-example/blob/master/server/server.go 
*/
func (s *chatServer) Chat(stream pb.Chat_ChatServer) error {

	id := uuid.New().String() // NOTE: this is temporarily and will not be generated from here. This should come as metadata

	s.addClient(id, stream) // Register client
	defer s.removeClient(id)

	for {
		req, err := stream.Recv()
		if err != nil {
			log.Printf("Recieve error: %v", err)
			return err
		}

		for key, ss := range s.clients {
			
			if key == id {
				continue // Do not send message back to client that submitted the message
			}

			if err := ss.Chat.Send(&pb.Message{Content: req.Content, Timestamp: req.Timestamp, Info: &pb.ClientInfo{Uuid: id, Name: s.clients[id].Name}}); err != nil {
				log.Printf("Could not send message for client id %s: %v", key, err)
			}
		}
	}

	return nil
}
