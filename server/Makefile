run-dev:
	go run main.go

run-prd:
	go build -o ./bin/server
	./bin/server

generate-api-docs:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g ./main.go -o ./docs
