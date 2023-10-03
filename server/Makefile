.PHONY: run-dev run-prd generate-api-docs

run-dev: # run in dev mode (using go run command)
	go run main.go

run-prd: # run in prd mode (building and running binary)
	go build -o ./bin/server
	./bin/server

generate-api-docs:	# generate openapi docs via swagger
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init -g ./main.go -o ./docs
