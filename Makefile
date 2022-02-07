
docker-build:
	@echo "===> Creating docker container"
	sudo docker build -t docker-tesseract .

docker-run:
	@echo "===> Running docker container"
	docker run -d -p 3000:3000 docker-tesseract

docker-all: docker-build docker-run

build:
	@echo "===> Building"
	CGO_ENABLED=1 GOOS=linux GARCH=amd64 go build -race -o app main.go

