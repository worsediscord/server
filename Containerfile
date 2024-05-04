FROM golang:1.22.1

WORKDIR /app

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o server

EXPOSE 8069

CMD ["./server"]