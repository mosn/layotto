FROM centos:centos7 as builder

ENV CHROOT_RUNTIME_PREFIX /ROOT/home/admin/layotto
RUN useradd -ms /bin/bash admin


COPY ./cmd/layotto/layotto $CHROOT_RUNTIME_PREFIX/bin/layotto
RUN mkdir -p $CHROOT_RUNTIME_PREFIX/conf \
 && mkdir -p $CHROOT_RUNTIME_PREFIX/logs

COPY ./configs $CHROOT_RUNTIME_PREFIX/configs
RUN chmod +x $CHROOT_RUNTIME_PREFIX/bin/layotto
RUN chown -R admin:admin /ROOT/home/admin

FROM centos:centos7

RUN yum install -y \
       ssh wget curl perl logrotate make build-essential procps etcd \
       tsar tcpdump mpstat iostat vmstat sysstat \
       python-setuptools; yum clean all

RUN useradd -ms /bin/bash admin
# copy all in one layer
COPY --from=builder /ROOT /
WORKDIR /home/admin/layotto
ENV RUNTIME_CONFIG /home/admin/layotto/configs/runtime_config.json
ENV RUNTIME_PREFIX /home/admin/layotto
ARG RUNTIME_CONFIG_DIR /home/admin/layotto/configs

ENTRYPOINT ["/home/admin/layotto/bin/layotto", "start"]
