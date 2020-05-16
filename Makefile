build:
	go build -v ./cmd/apigateway/.

run:
	go run -v ./cmd/apigateway/. -config_path=./configs/config.yaml

run_dev:
	go run -v ./cmd/apigateway/. -config_path=./configs/config.yaml.devel

test:
	go test -v -race ./...

.DEFAULT_GOAL := run