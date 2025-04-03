# Build
FROM golang:latest as build

WORKDIR /go/src/app

COPY . .
RUN CGO_ENABLED=0 make build

# Run
FROM gcr.io/distroless/static

COPY --from=build /go/src/app/bin/aiseg /aiseg

ENTRYPOINT ["/aiseg", "serve", "-v", "-d", "/db"]
