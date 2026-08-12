FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY ../common ../common
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /root/

COPY --from=builder /app/user-service .
COPY config ./config

EXPOSE 8001

CMD ["./user-service"]
