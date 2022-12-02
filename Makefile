all: requests

install-swagger:
	which swagger || go install github.com/go-swagger/go-swagger/cmd/swagger@v0.28.0

swagger.json: install-swagger
	go mod vendor && swagger generate spec -o ./swagger.json --scan-models

requests: swagger.json
	go build --buildvcs=false

clean:
	rm -rf swagger.json requests vendor

.PHONY: install-swagger clean all
