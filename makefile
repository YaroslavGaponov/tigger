.PHONY: demo build

demo:
		go run cmd/tigger/main.go assets/test.txt

build:
		go build cmd/tigger/main.go -o tigger