FROM golang:latest

WORKDIR ../..

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

COPY config/* ./config/ 

RUN go build -o gomf .

EXPOSE 8080

RUN mkdir -p /gomf/files

CMD ["./gomf", "config/gomf.config"]
