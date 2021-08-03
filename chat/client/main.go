package main

import (
	"context"
	"fmt"
	pb "github.com/mem-memov/buildsystemswithgo/chat/client/chat"
	"google.golang.org/grpc"
	"time"
)

func Chat(stream pb.ChatService_SendTxtClient, done chan bool) {
	t := time.NewTicker(500 * time.Millisecond)
	for {
		select {
		case <- done:
			return
		case <- t.C:
			err := stream.Send(&pb.ChatRequest{Txt: "Hello", Id: 1, To: 2})
			if err != nil {
				panic(err)
			}
		}
	}
}

func Stats(stream pb.ChatService_SendTxtClient, done chan bool) {
	for {
		stats, err := stream.Recv()
		if err != nil {
			panic(err)
		}
		fmt.Println(stats.String())
		if stats.TotalChar > 35 {
			fmt.Println("Beyond the limit!!!")
			done <- true
			stream.CloseSend()
			return
		}
	}
}

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	c := pb.NewChatServiceClient(conn)

	stream, err := c.SendTxt(context.Background())
	if err != nil {
		panic(err)
	}
	done := make(chan bool)

	go Chat(stream, done)
	go Stats(stream, done)

	<- done
}
