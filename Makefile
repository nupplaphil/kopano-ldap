test:
	sh test.sh

build:
	go build kopano-ld.go

.DEFAULT_GOAL := build
