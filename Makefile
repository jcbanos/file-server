build-server:
	go build serv/server.go

build-client:
	go build cli/client.go

run-server:
	serv/server start

clean:
	rm -rf 
