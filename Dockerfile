FROM golang:1.24-alpine as build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
RUN apk add --no-cache gcc musl-dev

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bot cmd/bot/main.go

FROM alpine:latest

RUN apk add --no-cache sqlite-libs

WORKDIR /root/
COPY --from=build /app/bot .
COPY .env .

CMD ["./bot"]

