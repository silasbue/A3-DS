package main

import (
	"context"
	"log"

	"github.com/silasbue/A3-DS.git/chitty_chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	waitc := make(chan struct{})
	conn, _ := grpc.Dial("localhost:5400", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := chitty_chat.NewChittyChatClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, _ := client.Chat(ctx)

	go func() {
		for {
			in, err := stream.Recv()
			/*if err == io.EOF {
				// read done.
				close(waitc)
				return
			}*/
			// if err != nil {
			// 	log.Fatalf("Failed to receive a note : %v", err)
			// }
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
				break
			}
			log.Println("Got message: ", in.Msg)
		}
	}()
	<-waitc
}
