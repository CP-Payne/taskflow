compile:
	protoc api/v1/*.proto \
		--go_out=./internal/proto \
		--go_opt=paths=source_relative \
		--proto_path=. 
test:
	go test -race ./...
