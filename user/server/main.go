package main

import (
	"context"
	"fmt"
	pb "github.com/mem-memov/buildsystemswithgo/user/server/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
}

func (u *UserServer) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.User, error) {
	fmt.Println("Server received: ", req.String())
	return &pb.User{UserId: "John", Email: "john@gmail.com"}, nil
}

func AuthServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
	) (interface{}, error) {

	md, found := metadata.FromIncomingContext(ctx)
	if ! found {
		return nil, status.Errorf(codes.InvalidArgument, "metadata not found")
	}

	password, found := md["password"]
	if !found {
		return nil, status.Errorf(codes.Unauthenticated, "password not found")
	}

	if password[0] != "go" {
		return nil, status.Errorf(codes.Unauthenticated, "password not valid")
	}

	return handler(ctx, req)
}

func main() {
	lis, err := net.Listen("tcp", "localhost:50051")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(AuthServerInterceptor))
	pb.RegisterUserServiceServer(s, &UserServer{})

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}
