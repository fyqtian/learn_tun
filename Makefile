build-client:
	go build -o bin/client client/client.go

build-server:
	go build -o bin/server server/server.go

run-client: build-client
	sudo ./bin/client --serverAddress ubuntu:8000

run-server: build-server
	sudo ./bin/server

