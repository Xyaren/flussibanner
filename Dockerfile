FROM golang:1.17-alpine

RUN mkdir /app
ADD ./ /app/
WORKDIR /app
RUN go build -o main ./cmd/flussibanner-server/

EXPOSE 8080
CMD ["/app/main"]
