FROM golang:1-alpine

RUN mkdir /pancors
WORKDIR /pancors
COPY . .
RUN go build -o pancors ./cmd/pancors/main.go

EXPOSE 8080

CMD [ "./pancors" ]