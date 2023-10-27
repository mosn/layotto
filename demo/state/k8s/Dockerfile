FROM golang:1.18-alpine AS build-env

WORKDIR /app

COPY . .

RUN go build -o client client.go

FROM alpine:latest
COPY --from=build-env /app/client .
CMD [ "./client", "-s", "state_demo" ]
