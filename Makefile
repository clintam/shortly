.PHONY: all test test-suite

PERF_TEST=docker-compose run tester test-suite --server http://server:8080 --concurrency 100 --iterations 1000 --initial-writes 100000

IMAGE = clintam/shortly

export GOPATH

all: test

build:
	docker build -t $(IMAGE) .

redis:
	docker-compose up -d redis

redis-cli:
	docker-compose exec redis redis-cli

mongo:
	docker-compose up -d mongo

mongo-cli:
	docker-compose exec mongo mongo

test: build redis mongo
	docker-compose run --rm server bash -c 'go test -race -coverprofile=c.out -v ./src/server && go tool cover -func=c.out'

memory-perftest: build
	STORAGE=memory docker-compose up -d server
	$(PERF_TEST)
	@echo === Memory

redis-perftest: build redis
	STORAGE=redis docker-compose up -d server
	$(PERF_TEST)
	@echo === Redis

mongo-perftest: build mongo
	STORAGE=mongo docker-compose up -d server
	$(PERF_TEST)
	@echo === Mongo

perftest-suite: memory-perftest redis-perftest mongo-perftest

clean:
	docker-compose down -v