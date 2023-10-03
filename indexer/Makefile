.PHONY: run-dev run-prd

prof.mode?=

run-dev:	# run in dev mode (using go run command)
	go run main.go -prof.mode=$(prof.mode)

run-prd:	# run in prd mode (building and running binary)
	go build -o ./bin/indexer
	./bin/indexer -prof.mode=$(prof.mode)
