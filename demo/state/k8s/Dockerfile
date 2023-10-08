FROM golang:1.18-alpine AS build-env

WORKDIR /app

COPY ../../../deploy/k8s/state .

RUN go build -o client client.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=build-env /app/client .
COPY config.json /usr/share/layotto/configure/
EXPOSE 8443
CMD [ "./client", "-s", "state_demo" ]
