FROM golang:alpine
LABEL maintainer="Rafael Teixeira <rafaelteixeiradev@gmail.com>"
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./
EXPOSE 3000
ENV CGO_ENABLED=0 GOOS=linux GARCH=amd64 GOGC=10000
RUN go build main.go
RUN apk update && apk add tesseract-ocr && apk add tesseract-ocr-data-por
CMD ["./main"]