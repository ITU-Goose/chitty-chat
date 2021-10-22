package internal

import pb "github.com/goose-alt/chitty-chat/api/v1/pb/chat"

// Internal representation of a client
type Client struct {
	Uuid       string
	Name       string
	Chat pb.Chat_ChatServer
}
