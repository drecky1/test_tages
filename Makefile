generate_proto:
	protoc --go_out=. --go-grpc_out=. proto/file.proto

run_server:
	go run server/main.go

run_client:
	go run example_client/main.go

run_tests:
	go test ./tests
	go fmt ./...
	go vet ./...

compile_server:
	CGO_ENABLED=0 go build -v -ldflags=" \
        -X 'build_data.User=$(id -u -n)'  \
        -X 'build_data.Time=$(date)'  \
        -X 'build_data.Version=$(cat version.txt)'" -o build/test_tages_server server/main.go