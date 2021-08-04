package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	pb "github.com/mem-memov/buildsystemswithgo/sum/client/sum"
	"time"
)

func main() {
	conn, err := grpc.Dial(":50051", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		panic(err)
	}

	c := pb.NewNumServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	stream, err := c.Sum(ctx)
	if err != nil {
		panic(err)
	}

	from, to := 1, 1000000

	for i := from; i<=to; i ++ {
		err := stream.Send(&pb.NumRequest{X: int64(i)})
		if err != nil {
			panic(err)
		}
	}

	fmt.Println("Waiting for response...")

	result, err := stream.CloseAndRecv()
	if err != nil {
		panic(err)
	}

	fmt.Printf("The sum from %d to %d is %d.\n", from, to, result.Total)
}