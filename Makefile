all: clean windows linux docker

clean:
	rm -rf build

windows:
	GOOS=windows GOARCH=amd64 go build -o build/windows/preservationnc-server.exe main.go

linux:
	GOOS=linux GOARCH=amd64 go build -o build/linux/preservationnc-server main.go

docker:
	docker-compose build
