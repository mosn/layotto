FROM golang:1.18-alpine AS build-env

WORKDIR /app

COPY . .

RUN go build -o injector cmd/layotto_injector/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build-env /app/injector .
EXPOSE 8443
