build-server:
	go build -o serv/server serv/server.go

build-client:
	go build -o cli/client cli/client.go

run-server:
	serv/server start

client-recieve-channel-2:
	cli/client recieve -channel 2

client-recieve-channel-3:
	cli/client recieve -channel 3

client-send-mytxt-channel-2:
	cli/client send mytxt.txt -channel 2

client-send-cities-channel-2:
	cli/client send cities.csv -channel 2

client-send-quiz-channel-3:
	cli/client send cities.csv -channel 3

clean:
	rm -rf cli/client-*/
	rm cli/client
	rm -rf serv/files/
	rm serv/server
