build:
	go build -o bin/server ./cmd/go-pagination/main.go

run:
	./bin/server

test:
	go test ./...

clean:
	rm -rf bin/server && clear