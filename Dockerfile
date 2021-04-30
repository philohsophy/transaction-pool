FROM golang:1.16.3-alpine3.13 AS builder
RUN mkdir /app
COPY . /app
WORKDIR /app
RUN go build -o main .

FROM alpine:3.13.5
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8010
CMD ["./main"]