build:
	go build -v ./cmd/apigateway/.

run:
	go build -v ./cmd/apigateway/. && ./apigateway -config_path=./configs/config.yaml

run_dev:
	go build -v ./cmd/apigateway/. && ./apigateway -config_path=./configs/config.yaml.devel

test:
	go test -v -race ./...

.DEFAULT_GOAL := run