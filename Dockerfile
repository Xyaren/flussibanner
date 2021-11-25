# build
FROM golang:1.17-alpine AS builder

RUN mkdir /app
ADD ./ /app/
WORKDIR /app

RUN CGO_ENABLED=0 go build -o main ./cmd/flussibanner-server/
# image

FROM alpine:latest AS builder2
RUN apk add --no-cache tzdata

FROM scratch
COPY --from=builder2 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder2 /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /app/main /flussibanner
EXPOSE 8080
ENTRYPOINT ["/flussibanner"]
