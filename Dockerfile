FROM golang:1.25-alpine AS builder

WORKDIR /go/app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/department ./cmd/main.go

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=builder /go/app/bin/department ./department

EXPOSE 8080

USER app

ENTRYPOINT ["./department"]
