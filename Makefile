
docker-build:
	@echo "===> Creating docker container"
	sudo docker build -t docker-tesseract .

docker-run:
	@echo "===> Running docker container"
	docker run -d -p 3000:3000 docker-tesseract

docker-all: docker-build docker-run

build:
	@echo "===> Building"
	CGO_ENABLED=0 GOOS=linux GARCH=amd64 go build -o app main.go

autocannon:
	@echo "===> Running autocannon with 5 connections for 10 seconds"
	autocannon http://localhost:3000/?image=https://image.slidesharecdn.com/portugus2b-170225215804/95/texto-verbal-e-noverbal-13-638.jgp -d 10 -c 5

autocannon-worker:
	@echo "===> Running autocannon with 7 connections and 5 workers for 10 seconds"
	autocannon http://localhost:3000/?image=https://image.slidesharecdn.com/portugus2b-170225215804/95/texto-verbal-e-noverbal-13-638.jgp -d 10 -c 7 -w 5