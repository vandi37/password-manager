FROM golang:1.24.1-alpine

WORKDIR /app

COPY go.mod go.sum ./
COPY . .
RUN go mod download
RUN go mod tidy
RUN go build -o main cmd/main.go

CMD [ "./main"]