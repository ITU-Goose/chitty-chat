package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
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
	// Generate a new uuid for the client
	id := uuid.New().String()
	timestamp := ts.CreateVectorTimestamp(id)

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

	waitc := make(chan struct{})
	stream, err := c.Chat(context.Background())
	if err!=nil{
		return
	}

	uuid,name := register(c,ctx,timestamp,stream,waitc)
	chat(c,ctx,timestamp,stream,waitc,uuid,name)

}

func readInput() string {
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return input
}

func register(c pb.ChatClient, ctx context.Context, timestamp ts.VectorTimestamp, stream pb.Chat_ChatClient, waitc chan struct{}) (userId string, name string) {
	name = askForUsername()

	timestamp.Increment()

	stream.Send(&pb.Message{Content: "",Timestamp: &pb.Lamport{Clients: timestamp.GetVectorTime()}, Info: &pb.ClientInfo{Uuid: "", Name: name}})
	in, err := stream.Recv()
	if err == io.EOF {
		// read done.
		close(waitc)
		return
	}
	if err != nil {
		log.Fatalf("Failed to receive a note : %v", err)
	}
	userId = in.Info.Uuid
	timestamp.Sync(in.Timestamp.Clients)
	log.Printf("Got message %s. From: %s. Timestamp: %s", in.Content, in.Info.Name, timestamp.GetDisplayableContent())
	  
	return userId,name		
}

func askForUsername() string {
	msg := "Enter your username"
	fmt.Println(msg)
	fmt.Print("-> ")

	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')

	return username
}

func chat(c pb.ChatClient, ctx context.Context, timestamp ts.VectorTimestamp, stream pb.Chat_ChatClient, waitc chan struct{}, uuid string, name string) {
	message := readInput()
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a message: %v", err)
			}
			timestamp.Sync(in.Timestamp.Clients)
			timestamp.Increment()
		  log.Printf("Got message %s. From: %s. Timestamp: %s", in.Content, in.Info.Name, timestamp.GetDisplayableContent())
		}
	  }()
	  mes := pb.Message{Content: message,Timestamp: &pb.Lamport{Clients: timestamp.GetVectorTime()}, Info: &pb.ClientInfo{Uuid: uuid, Name: name}}
		if err := stream.Send(&mes); err != nil {
		  log.Fatalf("Failed to send a note: %v", err)
		}
	stream.CloseSend()
	<-waitc

}
