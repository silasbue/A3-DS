package main

import (
	"io"
	"log"
	"net"
	"os"

	"github.com/silasbue/A3-DS.git/chitty_chat"
	"google.golang.org/grpc"
)

type Server struct {
	chitty_chat.UnimplementedChittyChatServer
	streams []chitty_chat.ChittyChat_ChatServer
}

func (s *Server) connect(newStream chitty_chat.ChittyChat_ChatServer) {
	s.streams = append(s.streams, newStream)
	// Notify all clients that a new client has joined
	nameMsg, _ := newStream.Recv()
	for _, client := range s.streams {
		client.SendMsg(&chitty_chat.Message{Username: "Server", Msg: nameMsg.Username + " has joined the chat"})
	}
}

func (s *Server) Chat(stream chitty_chat.ChittyChat_ChatServer) error {
	go func() {
		for {
			msg, err := stream.Recv()

			if err == io.EOF {
				return
			} else if err != nil {
				return
			}
			for _, client := range s.streams {
				client.SendMsg(&chitty_chat.Message{Username: msg.Username, Msg: msg.GetMsg()})
			}
		}
	}()

	waitc := make(chan struct{})
	s.connect(stream)

	<-waitc
	return nil
}

func main() {
	if len(os.Args) != 2 {
		log.Printf("Please input the port to run the server on")
	}

	lis, _ := net.Listen("tcp", "localhost:"+os.Args[1])

	grpcServer := grpc.NewServer()
	chitty_chat.RegisterChittyChatServer(grpcServer, &Server{})

	log.Printf("server listening at %v", lis.Addr())

	grpcServer.Serve(lis)
}
