# A simplified docker file for building layotto base image
# set default path to load config to `config.json`

FROM centos:centos7

COPY ./layotto /runtime/

RUN chmod +x  /runtime/layotto

ENTRYPOINT ["/runtime/layotto", "-c", "config.json"]
