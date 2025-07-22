FROM golang:1.24-alpine3.22 AS  builder

WORKDIR /app

COPY go.mod go.sum ./
COPY . .
RUN go mod download
RUN go mod tidy
RUN go build -o main cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/configs/ ./

CMD [ "/app/main"]