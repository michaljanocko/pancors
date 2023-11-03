FROM golang:1-alpine AS builder

RUN mkdir /pancors
WORKDIR /pancors

COPY go.mod go.mod
COPY pancors.go pancors.go
COPY cmd/ cmd/

RUN go build -o pancors ./cmd/pancors/main.go

FROM alpine:latest

COPY --from=builder /pancors/pancors /usr/bin/
EXPOSE 8080
USER nobody
CMD ["/usr/bin/pancors"]
