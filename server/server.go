package main

import (
	"context"
	"log"
	"net"

	"github.com/silasbue/A3-DS.git/chitty_chat"
	"google.golang.org/grpc"
)

type Server struct {
	chitty_chat.UnimplementedChittyChatServer
}

func (s *Server) GetMessage(ctx context.Context, in *chitty_chat.MessageRequest) (*chitty_chat.MessageReply, error) {
	log.Printf("Received message: %v", in.GetMsg())
	return &chitty_chat.MessageReply{Msg: "recieved: " + in.GetMsg()}, nil
}

func main() {
	lis, _ := net.Listen("tcp", "localhost:5400")

	grpcServer := grpc.NewServer()
	chitty_chat.RegisterChittyChatServer(grpcServer, &Server{})

	log.Printf("server listening at %v", lis.Addr())

	grpcServer.Serve(lis)
}
