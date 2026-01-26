FROM golang:1.25.5-alpine

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN mkdir -p build && go build -o build/todo cmd/todo/main.go

CMD ["./build/todo"]