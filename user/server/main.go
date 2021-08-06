package main

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/mem-memov/buildsystemswithgo/user/server/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"net"
	"net/http"
)

type UserServer struct {
	httpAddr string
	grpcAddr string
	pb.UnimplementedUserServiceServer
}

func (u *UserServer) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.User, error) {
	fmt.Println("Server received (GetUser): ", req.String())
	return &pb.User{UserId: "John", Email: "john@gmail.com"}, nil
}

func (u *UserServer) Create(ctx context.Context, req *pb.User) (*pb.User, error) {
	fmt.Println("Server received (Create): ", req.String())
	return &pb.User{UserId: "John", Email: "john@gmail.com"}, nil
}

func (u *UserServer) ServeGrpc() {
	lis, err := net.Listen("tcp", u.grpcAddr)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(grpc.UnaryInterceptor(AuthServerInterceptor))

	pb.RegisterUserServiceServer(s, u)
	fmt.Println("Server listening GRPC:")

	if err := s.Serve(lis); err != nil {
		panic(err)
	}
}

func (u *UserServer) ServeHttp() {
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	endpoint := u.grpcAddr

	err := pb.RegisterUserServiceHandlerFromEndpoint(context.Background(), mux, endpoint, opts)
	if err != nil {
		panic(err)
	}

	httpServer := &http.Server{
		Addr: u.httpAddr,
		Handler: mux,
	}

	fmt.Println("Server listening HTTP:")
	if err := httpServer.ListenAndServe(); err != nil {
		panic(err)
	}
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
	us := UserServer{httpAddr: ":8080", grpcAddr: ":50051"}
	go us.ServeGrpc()
	us.ServeHttp()
}

