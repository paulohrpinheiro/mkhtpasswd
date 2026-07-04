.PHONY: test coverage build vet clean

test:
	go test -race -cover ./...

coverage:
	go test -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

build:
	go build ./cmd/mkhtpasswd

vet:
	go vet ./...

clean:
	rm -f mkhtpasswd mkhtpasswd.exe coverage.out coverage.html