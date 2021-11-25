# build
FROM golang:1.17-alpine AS builder

RUN mkdir /app
ADD ./ /app/
WORKDIR /app

RUN go build -o main ./cmd/flussibanner-server/
RUN ls -ahln .
# image
FROM scratch
COPY --chmod=777 --from=builder /app/main /main
EXPOSE 8080
ENTRYPOINT ["/main"]
