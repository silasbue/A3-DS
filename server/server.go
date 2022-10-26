package main

import (
	"log"
	"net"

	"github.com/silasbue/A3-DS.git/chitty_chat"
	"google.golang.org/grpc"
)

type Server struct {
	chitty_chat.UnimplementedChittyChatServer
	streams []chitty_chat.ChittyChat_ChatServer
}

func (s *Server) Chat(stream chitty_chat.ChittyChat_ChatServer) error {
	s.streams = append(s.streams, stream)
	for _, client := range s.streams {
		client.Send(&chitty_chat.Message{Id: 1, Msg: "A new client has joined", T: 1})
	}

	return nil
}

func main() {
	lis, _ := net.Listen("tcp", "localhost:5400")

	grpcServer := grpc.NewServer()
	chitty_chat.RegisterChittyChatServer(grpcServer, &Server{})

	log.Printf("server listening at %v", lis.Addr())

	grpcServer.Serve(lis)
}
