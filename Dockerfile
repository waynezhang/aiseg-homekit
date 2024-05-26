FROM golang:latest

WORKDIR /app

COPY . .

RUN make build

ENTRYPOINT ["/app/bin/aiseg", "serve", "-v", "-d", "/db"]
