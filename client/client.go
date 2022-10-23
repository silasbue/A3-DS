package main

import (
	"context"
	"log"
	"time"

	"github.com/silasbue/A3-DS.git/chitty_chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, _ := grpc.Dial("localhost:5400", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := chitty_chat.NewChittyChatClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, _ := client.GetMessage(ctx, &chitty_chat.MessageRequest{Msg: "hi"})

	log.Printf("Reply from server: %s", r.GetMsg())

}
