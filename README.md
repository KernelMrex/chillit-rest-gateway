# Chillit REST Gateway

### About
Service for Web app, working with JSON format.

### Using

Compile `go build` and run `./chillit-rest-gateway [-config_path=<path>]`

### Configuration

Add file `config.yaml` to working directory.
 
``` yaml
store_service:
  url: "localhost:10050"
``` 