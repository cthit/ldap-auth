FROM golang:latest
WORKDIR /go/src/app
COPY . .

RUN go mod download
RUN go get -d
RUN go get github.com/codegangsta/gin

CMD ["gin", "-i", "run", "main.go"]