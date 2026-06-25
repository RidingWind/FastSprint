FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
COPY ../common ../common
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o project-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /root/

COPY --from=builder /app/project-service .
COPY config ./config

EXPOSE 8002

CMD ["./project-service"]
