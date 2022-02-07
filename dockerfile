FROM golang:1.17-alpine
LABEL maintainer="Rafael Teixeira <rafaelteixeiradev@gmail.com>"
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./
EXPOSE 3000
RUN GOOS=linux GARCH=amd64 go build -o app main.go
RUN apk update && apk add tesseract-ocr && apk add tesseract-ocr-data-por
CMD ["./app"]
