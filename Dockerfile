FROM docker.io/golang:1.23-alpine as build

WORKDIR /app

COPY . .

RUN go build -o bin cmd/service/main.go

FROM alpine:3.21

COPY --from=build /app/bin /app/bin

CMD ["/app/bin"]