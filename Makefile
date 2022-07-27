build:
	go build -o serv/server serv/server.go
	go build -o cli/client cli/client.go


clean:
	rm -rf cli/client-*/
	rm cli/client
	rm -rf serv/files/
	mkdir serv/files
	rm serv/server
