package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	pb "github.com/goose-alt/chitty-chat/api/v1/pb/chat"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:50051"
	defaultName = "world"
)

func main() {
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

}

func readInput() string{
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return input
}
