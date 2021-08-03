package main

import (
	"fmt"
	pb "github.com/mem-memov/buildsystemswithgo/chat/server/chat"
	"google.golang.org/grpc"
	"io"
	"net"
	"time"
)

type ChatServer struct{
	pb.UnimplementedChatServiceServer
}

func (c *ChatServer) SendTxt(stream pb.ChatService_SendTxtServer) error {
	var total int64 = 0
	go func() {
		t := time.NewTicker(2 * time.Second)
		for {
			select {
			case <- t.C:
				stream.Send(&pb.StatsResponse{TotalChar: total})
			}
		}
	}()

	for {
		next, err := stream.Recv()
		if err == io.EOF {
			fmt.Println("Client closed")
			return nil
		}
		if err != nil {
			return err
		}
		fmt.Println("->", next.Txt)
		total += int64(len(next.Txt))
	}
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()

	pb.RegisterChatServiceServer(s, &ChatServer{})

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}