# buildsystemswithgo

    cd server/chat
    protoc chat.proto --go_out=./ --go-grpc_out=./

    cd client/chat
    protoc chat.proto --go_out=./ --go-grpc_out=./

    cd client
    go mod tidy
    go build -o ../chat_client  github.com/mem-memov/buildsystemswithgo/chat/client

    cd server
    go mod tidy
    go build -o ../chat_server  github.com/mem-memov/buildsystemswithgo/chat/server

    ./chat_server
    ./chat_client