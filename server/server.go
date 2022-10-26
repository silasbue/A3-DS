package main

import (
	"io"
	"log"
	"net"

	"github.com/silasbue/A3-DS.git/chitty_chat"
	"google.golang.org/grpc"
)

type Server struct {
	chitty_chat.UnimplementedChittyChatServer
	streams []chitty_chat.ChittyChat_ChatServer
}

func (s *Server) connect(newStream chitty_chat.ChittyChat_ChatServer) {
	s.streams = append(s.streams, newStream)
	for _, client := range s.streams {
		client.SendMsg(&chitty_chat.Message{Id: 1, Msg: "A new client has joined", T: 1})
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
				client.SendMsg(&chitty_chat.Message{Msg: msg.GetMsg()})
			}
		}
	}()

	waitc := make(chan struct{})
	s.connect(stream)

	<-waitc
	return nil
}

func main() {
	lis, _ := net.Listen("tcp", "localhost:5400")

	grpcServer := grpc.NewServer()
	chitty_chat.RegisterChittyChatServer(grpcServer, &Server{})

	log.Printf("server listening at %v", lis.Addr())

	grpcServer.Serve(lis)
}
