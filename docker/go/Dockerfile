FROM golang:latest

WORKDIR ../..

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

COPY config/* ./config/ 

RUN go build -o main .

EXPOSE 8080

CMD ["./main", "config/gomf.config"]
