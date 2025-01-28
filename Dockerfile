FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o bin/main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN adduser -D -g '' myuser

WORKDIR /root/

COPY --from=build /app/bin/main .

RUN chown myuser:myuser /root/main
RUN chmod 755 /root/main  # Make it executable for the user

USER myuser

EXPOSE 8080

CMD ["/root/main"]  # Ensure absolute path, as the user may not have the root working directory by default
