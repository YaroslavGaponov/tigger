.PHONY: demo build

demo:
		go run cmd/tigger/main.go assets

build:
		go build cmd/tigger/main.go -o tigger