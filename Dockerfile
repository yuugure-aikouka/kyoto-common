FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN GOARCH=amd64 GOOS=linux go build -o bin/main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN adduser -D -g '' myuser

WORKDIR /root/

COPY --from=build /app/bin/main .

RUN chown myuser:myuser ./main

USER myuser

EXPOSE 8080

CMD ["./main"]
