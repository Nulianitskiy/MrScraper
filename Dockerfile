FROM golang:1.21-alpine

RUN apk update && apk add --no-cache poppler-utils

WORKDIR /go/src/app

COPY . .

RUN go mod tidy

RUN go build -o Scraper ./cmd/app

EXPOSE 8080

CMD ["./Scraper"]