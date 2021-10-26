package main

import (
	"bufio"
	"context"
	"log"
	"os"

	"github.com/google/uuid"
	pb "github.com/goose-alt/chitty-chat/api/v1/pb/chat"
	ts "github.com/goose-alt/chitty-chat/internal/time"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
	
)

func main() {

	timestamp := ts.CreateVectorTimestamp("abe")

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewChatClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	message := readInput()

	chat(c,ctx,message,timestamp)

}

func readInput() string {
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return input
}

func chat(c pb.ChatClient, ctx context.Context, message string, timestamp ts.VectorTimestamp) {
	timestamp.Increment()
	c.Chat(ctx, pb.Message{Content: message,Timestamp: &pb.Lamport{Clients: timestamp.GetVectorTime()}, Info: &pb.ClientInfo{Uuid: "søde smukke", Name: "Amalie"}})
}
