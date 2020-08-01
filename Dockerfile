FROM golang:1.14

WORKDIR /go/src/github.com/abeatrice/acl
COPY . .

RUN go get -u github.com/go-sql-driver/mysql
RUN go get github.com/joho/godotenv
RUN go get -u github.com/gorilla/mux

CMD ["go","run","main.go"]
