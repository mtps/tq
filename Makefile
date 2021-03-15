all: build


build:
	go build -mod=vendor ./cmd/tq

install:
	go install -mod=vendor ./cmd/tq
