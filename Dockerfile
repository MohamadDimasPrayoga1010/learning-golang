FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o learning-golang main.go

FROM alpine:latest

WORKDIR /backend

COPY --from=builder /app/learning-golang /backend/learning-golang

COPY --from=builder /app/templates /backend/templates

EXPOSE 8080

ENTRYPOINT ["/backend/learning-golang"]
