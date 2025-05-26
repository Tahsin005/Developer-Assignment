FROM golang:1.24.2-alpine as builder

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

RUN go mod tidy

COPY . .

COPY .env .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/.env .

ENTRYPOINT ["./main"]

EXPOSE 8080