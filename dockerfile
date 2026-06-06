FROM golang:1.26.2

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o main ./cmd/app

RUN ls -la /app/main

CMD ["./main"]