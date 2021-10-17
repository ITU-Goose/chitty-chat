package main

import (
	"log"
	"net"

	pb "github.com/goose-alt/chitty-chat/api/v1/pb/chat"
	pkg "github.com/goose-alt/chitty-chat/pkg/server"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to open listener: %v", err)
	}

	server := pkg.NewChatServer()

	s := grpc.NewServer()
	pb.RegisterChatServer(s, &server)

	log.Printf("Listening on: %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
