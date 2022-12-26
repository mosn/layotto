FROM golang:1.18

RUN apt-get update && apt-get install -y redis-server
RUN curl -L -O https://github.com/wasmerio/wasmer/releases/download/2.2.1/wasmer-linux-amd64.tar.gz
RUN tar zvxf wasmer-linux-amd64.tar.gz
RUN cp lib/libwasmer.so /usr/lib/libwasmer.so
