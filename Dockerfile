FROM golang:latest

WORKDIR /app

COPY . .

RUN go build -o ./go-url-shortner ./cmd/go-url-shortner/main.go

EXPOSE 8000

RUN chmod a+x ./go-url-shortner

CMD ["./go-url-shortner"]