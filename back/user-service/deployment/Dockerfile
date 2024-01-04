FROM golang:1.21.5

WORKDIR /app/project

COPY ./ ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./user-service/api/main.go

WORKDIR /app

RUN  rm -rf ./project

CMD ["/server"]