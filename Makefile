all:
	go build -C ./cmd/wdscmd -o ../../bin/wdscmd

swag:
	swag fmt && swag init -d api/ -g server.go -ot json -o ./api/docs

test:
	go test -v ./...

clean:
	rm -rf bin/