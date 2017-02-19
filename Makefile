.PHONY: all test test-suite

IMAGE = clintam/shortly

export GOPATH

all: test test-suite

build:
	docker build -t $(IMAGE) .

test: build
	docker run --rm $(IMAGE) bash -c 'go test -race -coverprofile=c.out -v ./src/server && go tool cover -func=c.out'

test-suite: build
	docker-compose up -d
	docker-compose exec server test-suite --concurrency 100

clean:
	docker-compose down -v