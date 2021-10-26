package server

import (
	"errors"
	"fmt"
	"io"
	"sync"

	pb "github.com/goose-alt/chitty-chat/api/v1/pb/chat"
	"github.com/goose-alt/chitty-chat/internal"
	"github.com/goose-alt/chitty-chat/internal/logging"
	lm "github.com/goose-alt/chitty-chat/internal/message"
	"github.com/goose-alt/chitty-chat/internal/time"
)

type chatServer struct {
	pb.UnimplementedChatServer

	// List of clients, mapped by their generated id
	clients   map[string]*internal.Client
	Logger    logging.Log
	lock      sync.Mutex
	timestamp *time.VectorTimestamp
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

func (s *chatServer) sendMessage(target *pb.Chat_ChatServer, message *pb.Message) error {
	s.timestamp.Increment()
	message.Timestamp = &pb.Lamport{Clients: s.timestamp.GetVectorTime()}

	lm.PrintMessage("Sending message", s.Logger, message.Content, s.timestamp, message.Info.Uuid)

	return (*target).Send(message)
}

func (s *chatServer) sendSystemMessage(target *pb.Chat_ChatServer, msg string) error {
	s.timestamp.Increment()
	message := &pb.Message{
		Content: msg,
		Info: &pb.ClientInfo{
			Uuid: serverId,
			Name: serverName,
		},
	}
	lm.PrintMessage("Sending system message", s.Logger, msg, s.timestamp, message.Info.Uuid)

	return (*target).Send(message)
}

func (s *chatServer) broadcastSystemMessage(message string) {
	s.Logger.IPrintf("Broadcasting system message: \"%s\"\n", message)
	s.broadcast(&pb.Message{
		Content: message,
		Info: &pb.ClientInfo{
			Uuid: serverId,
			Name: serverName,
		},
	})
}

func (s *chatServer) broadcast(message *pb.Message) {
	lm.PrintMessage("Broadcasting message", s.Logger, message.Content, s.timestamp, message.Info.Uuid)

	for key, ss := range s.clients {
		if err := s.sendMessage(&ss.Chat, message); err != nil {
			s.Logger.EPrintf("Could not send message for client id %s: %v\n", key, err)
		}
	}
}

func (s *chatServer) recieveMessage(stream *pb.Chat_ChatServer) (*pb.Message, error) {
	req, err := (*stream).Recv()
	if err != nil {
		return nil, err
	}

	s.timestamp.Sync(req.Timestamp.Clients)
	s.timestamp.Increment()
	lm.PrintMessage("Recieved message", s.Logger, req.Content, s.timestamp, req.Info.Uuid)

	return req, nil
}

func (s *chatServer) addClient(req *pb.Message, stream pb.Chat_ChatServer) (*internal.Client, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	id, name := req.Info.Uuid, req.Info.Name

	// TODO: Validate Id to be unique
	if id == "" || name == "" {
		return nil, errors.New(fmt.Sprintf("Missing information. id=%s,name=%s", id, name))
	} else if _, exists := s.clients[id]; exists {
		return nil, errors.New(fmt.Sprintf("Client with id: %s, already exists", id))
	}

	client := &internal.Client{
		Uuid: id,
		Name: name,
		Chat: stream,
	}
	s.clients[id] = client

	s.Logger.IPrintf("Client connected. (%s)%s: %s\n", client.Name, client.Uuid)
	s.broadcastSystemMessage("User joined: " + name)

	return client, nil
}

func (s *chatServer) removeClient(client *internal.Client) {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.Logger.IPrintf("Client disconnected. (%s)%s: %s\n", client.Name, client.Uuid)

	delete(s.clients, client.Uuid)

	s.broadcastSystemMessage("User left: " + client.Name)
}

/*
Is a stream to send chat messages. This is bidirectional.

The implementation is inspired by: https://github.com/castaneai/grpc-broadcast-example/blob/master/server/server.go
*/
func (s *chatServer) Chat(stream pb.Chat_ChatServer) error {
	clientRegistered := false

	client := &internal.Client{}
	s.timestamp.Increment()

	for {
		req, err := s.recieveMessage(&stream)
		if err != nil {
			s.removeClient(client)

			fmt.Printf("%T\n", err)
			if errors.Is(err, io.EOF) {
				return nil
			}

			s.Logger.EPrintf("Recieve error: %v\n", err)
			return err
		}

		if !clientRegistered {
			client, err = s.addClient(req, stream) // Register Client
			if err != nil {
				s.Logger.EPrintf("Failed to add client: %v\n", err)
				s.sendSystemMessage(&stream, err.Error())
				return err
			}

			clientRegistered = true
		}

		s.broadcast(&pb.Message{
			Content: req.Content,
			Info:    &pb.ClientInfo{Uuid: client.Uuid, Name: client.Name},
		})
	}

	return nil
}
