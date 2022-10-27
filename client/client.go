package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"

	"github.com/silasbue/A3-DS.git/chitty_chat"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	if len(os.Args) != 3 {
		log.Printf("Please run the client with an URL and a username")
		return
	}

	waitc := make(chan struct{})
	conn, _ := grpc.Dial(os.Args[1], grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := chitty_chat.NewChittyChatClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stream, _ := client.Chat(ctx)
	stream.Send(&chitty_chat.Message{Username: os.Args[2]})

	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			log.Println(in.Username + ": " + in.Msg)
		}
	}()

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		err := stream.Send(&chitty_chat.Message{Username: os.Args[2], Msg: msg})

		if err != nil {
			panic(err)
		}

	}

	<-waitc
}
