FROM golang:1.19.3

RUN mkdir -p /api
WORKDIR /api
COPY . .
RUN go mod download
RUN go build -o api
ENTRYPOINT ["./api"]