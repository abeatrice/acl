FROM golang:1.14

WORKDIR /go/src/github.com/abeatrice/acl
COPY . .

RUN go get -v -t -d ./...

CMD ["go", "build", "-v", "."]
