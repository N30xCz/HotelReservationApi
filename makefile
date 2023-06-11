build:
	@go build -o bin/api

run: build
	./bin/api

rm :
	@rm -rf ./bin

seed:
	@go run scripts/seed.go
test:
	@go test -v ./...
