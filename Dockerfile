FROM golang:latest

WORKDIR ~/server

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o=./bin/floodcontrol main.go

CMD ["./bin/floodcontrol 1>&1"]