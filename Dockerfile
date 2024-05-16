FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy && go mod download

COPY . .

WORKDIR /app/cmd/url-shortener

RUN go build -o main .

EXPOSE 3002

CMD ["./main"]
