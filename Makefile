test:
	go test -v ./... -coverprofile cover.out

run:
	go mod download
	go mod verify
	go run main.go

docker-build: test
	docker build -t url-shortener:1.0 .

docker-run:
	docker run -d -p 3000:3000  --name url-shortener  url-shortener:1.0

docker-stop:
	docker stop url-shortener
	docker rm url-shortener

