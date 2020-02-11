FROM golang:latest
WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go get -d

CMD ["go", "run", "main.go"]