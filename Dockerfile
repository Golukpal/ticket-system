FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o ticket-server ./cmd/server


FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/ticket-server .

EXPOSE 8080

CMD ["./ticket-server"]