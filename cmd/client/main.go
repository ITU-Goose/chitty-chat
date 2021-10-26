package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"

	pb "github.com/goose-alt/chitty-chat/api/v1/pb/chat"
	"github.com/goose-alt/chitty-chat/internal/logging"
	pkgClient "github.com/goose-alt/chitty-chat/pkg/client"
	"google.golang.org/grpc"
)

const (
	defaultName = "world"
)

func main() {
	logger := logging.New()

	host := flag.String("host", "localhost", "The host to connect to, usually an IP address")
	port := flag.String("port", "50051", "The port of the server to connect to")
	random := flag.Bool("random", false, "Chat randomly with the server, requiring no user input")
	flag.Parse()

	address := fmt.Sprintf("%s:%s", *host, *port)

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.EPrintf("Could not connect: %v\n", err)
	}
	defer conn.Close()

	c := pb.NewChatClient(conn)

	// Contact the server and print out its response.
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, err := c.Chat(context.Background())
	if err != nil {
		logger.EPrintf("Could not open stream: %v\n", err)
	}

	name := defaultName
	if !*random {
		name = askForUsername()
	} else {
		t, err := os.Hostname()
		if err != nil {
			logger.EPrintf("Could not open stream: %v\n", err)
			return
		}

		name = t
	}

	client := pkgClient.NewClient(name, &stream)

	if !*random {
		chat(client, logger)
	} else {
		randomChat(client, logger)
	}
}

func readInput() string {
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	return input
}

func askForUsername() string {
	msg := "Enter your username"
	fmt.Println(msg)
	fmt.Print("-> ")

	reader := bufio.NewReader(os.Stdin)
	username, _ := reader.ReadString('\n')
	username = strings.Trim(username, "\n \r")
	return username
}

func chat(client *pkgClient.Client, logger logging.Log) {
	callback := func(message string) {
		client.SendMessage(message)
	}
	ui := pkgClient.NewUi(callback)

	go func() {
		for {
			in, err := (*client.Stream).Recv()
			if err == io.EOF {
				// read done.
				(*client.Stream).CloseSend()
				return
			} else if err != nil {
				logger.EPrintf("Error while receiving: %v", err)
			}

			client.Timestamp.Sync(in.Timestamp.Clients)
			client.Timestamp.Increment()
			ui.AddMessage(fmt.Sprintf("%s (TS: %s): %s", in.Info.Name, client.Timestamp.GetDisplayableContent(), in.Content))
		}
	}()

	ui.Run()
}

func randomChat(client *pkgClient.Client, logger logging.Log) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	go func() {
		for {
			in, err := (*client.Stream).Recv()
			if err != nil {
				logger.EPrintf("Error while receiving: %v", err)
			}

			client.Timestamp.Sync(in.Timestamp.Clients)
			client.Timestamp.Increment()
			logger.IPrintf("Recieved message: %s, from: %s, at: %s\n", in.Content, in.Info.Name, client.Timestamp.GetDisplayableContent())
		}
	}()

	for {
		client.SendMessage(fmt.Sprintf("Random number: %d", r.Intn(100)))

		time.Sleep(time.Duration(r.Intn(10)) * time.Second)
	}
}
