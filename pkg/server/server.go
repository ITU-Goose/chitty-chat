package server

import (
	"sync"

	"github.com/google/uuid"

	pb "github.com/goose-alt/chitty-chat/api/v1/pb/chat"
	"github.com/goose-alt/chitty-chat/internal"
	"github.com/goose-alt/chitty-chat/internal/logging"
	"github.com/goose-alt/chitty-chat/internal/time"
)

type chatServer struct {
	pb.UnimplementedChatServer

	// List of clients, mapped by their generated id
	clients   map[string]*internal.Client
	Logger    logging.Log
	lock      sync.Mutex
	timestamp time.VectorTimestamp
}

const (
	serverId   = "11111111-1111-1111-1111-111111111111"
	serverName = "Server"
)

func NewChatServer() chatServer {
	return chatServer{
		clients:   make(map[string]*internal.Client),
		Logger:    logging.New(),
		timestamp: time.CreateVectorTimestamp(serverId),
	}
}

func (s *chatServer) addClient(stream pb.Chat_ChatServer) *internal.Client {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Generate a new uuid for the client
	id := uuid.New().String()

	// TODO: Replace name with username
	client := internal.Client{
		Uuid: id,
		Name: "",
		Chat: stream,
	}

	s.clients[id] = &client

	s.Logger.IPrintf("Client connected. ID: %s\n", client.Uuid)

	return &client
}

func (s *chatServer) removeClient(client *internal.Client) {
	s.lock.Lock()
	defer s.lock.Unlock()

	delete(s.clients, client.Uuid)

	s.Logger.IPrintf("Client disconnected. ID: %s\n", client.Uuid)

	s.timestamp.Increment()
	s.broadcast(&pb.Message{
		Content:   "User disconnected: " + client.Name,
		Timestamp: &pb.Lamport{Clients: s.timestamp.GetVectorTime()}, // TODO: Hmmmm, what to put here?
		Info: &pb.ClientInfo{
			Uuid: serverId,
			Name: serverName,
		},
	})
}

/*
Is a stream to send chat messages. This is bidirectional.

The implementation is inspired by: https://github.com/castaneai/grpc-broadcast-example/blob/master/server/server.go
*/
func (s *chatServer) Chat(stream pb.Chat_ChatServer) error {
	client := s.addClient(stream) // Register client
	defer s.removeClient(client)

	for {
		req, err := stream.Recv()
		if err != nil {
			s.Logger.EPrintf("Recieve error: %v\n", err)
			return err
		}

		s.Logger.IPrintf("Recieved message: %v\n", req)

		if client.Name == "" {
			if req.Info.Name != "" {
				s.setClientName(client.Uuid, req.Info.Name)
			} else {
				client.Chat.Send(&pb.Message{
					Content:   "Error: Your name is not yet set",
					Timestamp: req.Timestamp,
					Info:      &pb.ClientInfo{Uuid: client.Uuid, Name: client.Name},
				})

				continue
			}
		}

		s.broadcast(&pb.Message{
			Content:   req.Content,
			Timestamp: req.Timestamp,
			Info:      &pb.ClientInfo{Uuid: client.Uuid, Name: client.Name},
		})
	}

	return nil
}

func (s *chatServer) setClientName(id string, name string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	client := s.clients[id]
	client.Name = name
	
	s.timestamp.Increment()
	client.Chat.Send(&pb.Message{
		Timestamp: &pb.Lamport{Clients: s.timestamp.GetVectorTime()},
		Info: &pb.ClientInfo{
			Uuid: id,
			Name: name,
		},
	})

	s.timestamp.Increment()
	s.broadcast(&pb.Message{
		Content:   "User joined: " + name,
		Timestamp: &pb.Lamport{Clients: s.timestamp.GetVectorTime()}, // TODO: Hmmmm, what to put here?
		Info: &pb.ClientInfo{
			Uuid: serverId,
			Name: serverName,
		},
	})
}

func (s *chatServer) broadcast(message *pb.Message) {
	s.Logger.IPrintf("Broadcasting message: %v\n", message)

	for key, ss := range s.clients {
		if err := ss.Chat.Send(message); err != nil {
			s.Logger.EPrintf("Could not send message for client id %s: %v\n", key, err)
		}
	}
}
