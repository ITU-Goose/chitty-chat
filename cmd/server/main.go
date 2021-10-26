package main

import (
	"net"

	pb "github.com/goose-alt/chitty-chat/api/v1/pb/chat"
	pkg "github.com/goose-alt/chitty-chat/pkg/server"
	"google.golang.org/grpc"
)

const (
	port = ":50051"
)

func main() {
	server := pkg.NewChatServer()

	listener, err := net.Listen("tcp", port)
	if err != nil {
		server.Logger.EPrintf("Failed to open listener: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterChatServer(s, &server)

	server.Logger.IPrintf("Listening on: %v", listener.Addr())
	if err := s.Serve(listener); err != nil {
		server.Logger.EPrintf("Failed to serve: %v", err)
	}
}
