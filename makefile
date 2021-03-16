build:
	GOARCH=wasm GOOS=js go build -o ./web ./...
	mv ./web/app ./web/app.wasm
	go build -o Project1.exe ./server

run: build
	PORT=8000 ./Project1.exe

test:
	go test -p 1 ./...