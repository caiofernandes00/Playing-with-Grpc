protogen:
	protoc --proto_path=pkg/proto --go_out=pkg/proto/pb --go_opt=paths=source_relative \
	--go-grpc_out=pkg/proto/pb --go-grpc_opt=paths=source_relative,require_unimplemented_servers=false \
	--grpc-gateway_out pkg/proto/pb --grpc-gateway_opt paths=source_relative \
	--openapiv2_out pkg/proto/openapi --openapiv2_opt logtostderr=true \
	pkg/proto/*.proto

test:
	go test -cover -race ./...

server-grpc:
	go run cmd/server/main.go -port 8080 -tls true -type grpc

server-rest:
	go run cmd/server/main.go -port 8080 -tls false -type rest

client-create:
	go run cmd/client/main.go -address 0.0.0.0:8080 -operation create -tls true

client-search:
	go run cmd/client/main.go -address 0.0.0.0:8080 -operation search -tls true

client-upload:
	go run cmd/client/main.go -address 0.0.0.0:8080 -operation upload -tls true

client-rate:
	go run cmd/client/main.go -address 0.0.0.0:8080 -operation rate -tls true

cert:
	./cert/gen.sh

.PHONY: protogen test server-grpc server-rest client-create client-search client-upload client-rate cert