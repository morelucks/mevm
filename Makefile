all: build

build:
	@mkdir -p ../../bin

	@cd cmd/ && go build -o ../bin/mevm
run:
	@cd cmd/ && go run main.go
