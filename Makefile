include .env

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-w -s" -o ./order-service ./cmd/order_service/main.go

.PHONY: prepare
prepare:
	go mod download

.PHONY: order-service
order-service:
	./order-service

.PHONY: clean
clean:
	rm ./order-service

.PHONY: lint
lint:
	golangci-lint run ./... --fix

.PHONY: test
test:
	go test -v ./... --cover

.PHONE: docker-build
docker-build:
	docker build -t flykarlikimages/order:latest .