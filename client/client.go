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

var T int32

func main() {
	if len(os.Args) != 3 {
		log.Printf("Please run the client with an URL and a username")
		return
	}

	T = 0

	waitc := make(chan struct{})
	conn, _ := grpc.Dial(os.Args[1], grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()

	client := chitty_chat.NewChittyChatClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	//Joining the chat
	stream, _ := client.Chat(ctx)
	T++
	stream.Send(&chitty_chat.Message{Username: os.Args[2], T: T})

	//Recieve messages
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
			// Update timestamp
			T = Max(T, in.GetT()) + 1
			// Log message
			log.Println(in.Username+": "+in.Msg, "Lamport:", in.GetT())
		}
	}()

	// Send messages
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		T++
		err := stream.Send(&chitty_chat.Message{Username: os.Args[2], Msg: msg, T: T})

		if err != nil {
			panic(err)
		}

	}

	<-waitc
}

func Max(i int32, j int32) int32 {
	if i > j {
		return i
	}

	return j
}
