# build
FROM golang:1.17-alpine AS builder

RUN mkdir /app
ADD ./ /app/
WORKDIR /app

RUN CGO_ENABLED=0 go build -o main ./cmd/flussibanner-server/
# image
FROM scratch
COPY --from=alpine:latest /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=alpine:latest /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/main /main
EXPOSE 8080
ENTRYPOINT ["/main"]
