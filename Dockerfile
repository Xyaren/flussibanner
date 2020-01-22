FROM golang:alpine

RUN mkdir /app
ADD ./ /app/
WORKDIR /app
RUN go build -o main ./cmd/flussibanner-server/

EXPOSE 8080
CMD ["/app/main"]