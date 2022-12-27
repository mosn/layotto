FROM golang:1.18

RUN apt-get update && apt-get install -y redis-server \
        && apt-get install -y etcd-server   \
        && apt-get install -y net-tools   \
        && apt-get install -y zookeeper   \