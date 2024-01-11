build:
	@go build -o ./bin/go-url-shortner ./cmd/go-url-shortner/main.go

run: build
	@./bin/go-url-shortner

tidy:
	@go mod tidy
