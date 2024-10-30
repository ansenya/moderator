FROM golang:1.23.0-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /build

FROM alpine:latest

COPY --from=builder /build /app/build

WORKDIR /app

CMD ["/app/build"]
