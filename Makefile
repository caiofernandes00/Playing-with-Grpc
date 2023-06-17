protogen:
	protoc --proto_path=pkg/proto --go_out=pkg/proto/pb --go_opt=paths=source_relative \
	--go-grpc_out=pkg/proto/pb --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false \
	pkg/proto/*.proto

test:
	go test -cover -race ./...

server:
	go run cmd/server/main.go -port 8080

client-create:
	go run cmd/client/main.go -address 0.0.0.0:8080 -operation create

client-search:
	go run cmd/client/main.go -address 0.0.0.0:8080 -operation search

client-upload:
	go run cmd/client/main.go -address 0.0.0.0:8080 -operation upload

client-rate:
	go run cmd/client/main.go -address 0.0.0.0:8080 -operation rate