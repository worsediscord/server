all:
	go build -C ./cmd/ -o ../bin/wdscmd

swag:
	swag fmt && swag init -d api/ -g server.go -ot yaml -o ./

test:
	go test -v ./...

clean:
	rm -rf bin/