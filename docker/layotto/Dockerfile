# A simplified docker file for building layotto base image

FROM centos:centos7

COPY ./layotto /runtime/

WORKDIR /runtime

RUN chmod +x  /runtime/layotto

ENTRYPOINT ["/runtime/layotto"]
