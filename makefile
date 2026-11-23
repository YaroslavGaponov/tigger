.PHONY: demo build gen genapi genclient up


gen: genapi genclient

genapi:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest --config=api/oapi-codegen.yaml  api/openapi/tigger.spec.yaml

genclient:
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest --config=api/oapi-codegen-client.yaml  api/openapi/tigger.spec.yaml

demo:
		go run cmd/tigger/main.go assets

build:
		go build cmd/tigger/main.go -o tigger

up:
		docker-compose  -f ./docker-compose/docker-compose.yaml up --build