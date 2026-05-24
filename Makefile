APP_NAME = nsdata

build:
	go build -o ${APP_NAME} .
	mkdir -p bin
	mv ${APP_NAME} bin

test:
	go test ./...
