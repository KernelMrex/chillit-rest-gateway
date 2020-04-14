# Chillit REST Gateway

### About
Service for Web app, working with JSON format.

### Using

Compile `make build` and run `./chillit-rest-gateway[-config_path=<path>]` or just run with `make run`

### Configuration

Add file `config.yaml` to working directory.
 
``` yaml
api_server:
  hostname: ":8080"
  
store_service:
  url: "localhost:10050"
``` 