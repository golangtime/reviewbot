all: build test deploy

build:
	go build main.go -o reviewbot

test:
	go test ./...

deploy:
	echo "deploy"