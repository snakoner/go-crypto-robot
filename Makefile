default_goal := all

all: build clean

build:
	go build -o app -v ./cmd/app

clean:
	rm output.log