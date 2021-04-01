FROM golang:1-alpine

RUN mkdir /pancors
WORKDIR /pancors
COPY main.go main.go
COPY go.mod go.mod
RUN go build -o pancors .

EXPOSE 8080

CMD [ "/pancors/pancors" ]