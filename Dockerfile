FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN apk add build-base && go build -o forum cmd/main.go

FROM alpine:3.6
LABEL Authors="Markhabat&Asemokamichi" Project="Forum"
WORKDIR /app
COPY --from=builder /app .
CMD ["/app/forum"]