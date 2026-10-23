.PHONY: build test run docker-build docker-run example clean

build:
	CGO_ENABLED=1 go build -o sampler .

test:
	CGO_ENABLED=1 go test -v ./...

docker-build:
	docker build -t sampler .

docker-run:
	docker run -it sampler

example: docker-build docker-run

clean:
	rm -f sampler
