build:
	go build -v -o bin/simdrone cmd/main.go

run: build
	./bin/simdrone
