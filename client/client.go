package main

import (
	"context"
	"fmt"
	"log"

	"github.com/silasbue/A3-DS.git/chitty_chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, _ := grpc.Dial("localhost:5400", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := chitty_chat.NewChittyChatClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for i := 0; i < 2; i++ {
		var input string

		fmt.Scanln(&input)

		r, _ := client.GetMessage(ctx, &chitty_chat.MessageRequest{Msg: input})

		log.Printf("Reply from server: %s", r.GetMsg())
	}

	conn.Close()

}
