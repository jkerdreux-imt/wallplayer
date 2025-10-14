build:
	go mod tidy
	go build -o wallplayer ./cmd

dev:
	VIDEOS_DIR=~/Videos DEV=1 go run ./cmd

test:clean build
	mkdir ./target
	mv ./wallplayer ./target
	cd ./target && VIDEOS_DIR=~/Videos ./wallplayer

clean:
	rm -rf ./target wallplayer
