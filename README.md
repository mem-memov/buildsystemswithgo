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




    brew tap bufbuild/buf
    brew install buf
    
    go install \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
        github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
        google.golang.org/protobuf/cmd/protoc-gen-go \
        google.golang.org/grpc/cmd/protoc-gen-go-grpc

    cd user/client
    go mod tidy
    buf generate

    cd user/server
    go mod tidy
    buf generate

    cd user/client
    go mod tidy
    go build -o ../user_client

    cd user/server
    go mod tidy
    go build -o ../user_server

    ./user_server

    curl http://localhost:8080/v1/user/john -H "Grpc-Metadata-password: go"
    curl http://localhost:8080/v1/user -H "Grpc-Metadata-password: go" -d "{\"userId\":\"John\",\"email\":\"john@gmail.com\"}"

    ./user_client