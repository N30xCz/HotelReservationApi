build:
	@go build -o bin/api

run: build
	./bin/api

rm :
	@rm -rf ./bin