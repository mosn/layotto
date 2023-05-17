#!/bin/bash -e

set -e

# won't test:
# docs/en/start/configuration/start-apollo.md
# docs/zh/start/configuration/start-apollo.md
# because the github workflow can not connect to the apollo server due to the great firewall

GO_VERSION=${1:-"default"}

# By default,we test these docs with golang 1.18
quickstarts_in_default="docs/en/start/configuration/start.md
  docs/zh/start/lifecycle/start.md
  docs/en/start/lifecycle/start.md
  docs/zh/start/configuration/start.md
  docs/en/start/state/start.md
  docs/zh/start/state/start.md
  docs/en/start/pubsub/start.md
  docs/zh/start/pubsub/start.md
  docs/en/start/lock/start.md
  docs/zh/start/lock/start.md
  docs/en/start/sequencer/start.md
  docs/zh/start/sequencer/start.md
  docs/en/start/rpc/helloworld.md
  docs/zh/start/rpc/helloworld.md
  docs/zh/start/file/minio.md
  docs/en/start/api_plugin/helloworld.md
  docs/zh/start/api_plugin/helloworld.md
  docs/zh/start/actuator/start.md
  docs/en/start/actuator/start.md
  docs/zh/start/trace/trace.md
  docs/en/start/trace/trace.md
  docs/en/start/trace/skywalking.md
  docs/zh/start/trace/skywalking.md
  docs/zh/start/trace/prometheus.md
  docs/en/start/trace/prometheus.md
  docs/zh/start/trace/zipkin.md
  docs/zh/start/trace/jaeger.md
  docs/en/start/wasm/start.md
  docs/zh/start/wasm/start.md
  docs/en/start/secret/start.md
  docs/zh/start/secret/start.md
  docs/en/start/secret/secret_ref.md
  docs/zh/start/secret/secret_ref.md
  docs/zh/start/uds/start.md
"

# In advance mod, we test these docs with golang 1.17
quickstarts_in_advance="docs/en/start/rpc/dubbo_json_rpc.md
docs/zh/start/rpc/dubbo_json_rpc.md"

# download mdx
if ! test -e $(pwd)/etc/script/mdx; then
  curl -o $(pwd)/etc/script/mdx https://raw.githubusercontent.com/seeflood/mdx/main/mdx
fi
chmod +x $(pwd)/etc/script/mdx

# download etcd
sh etc/script/download_etcd.sh

# release all resources
release_resource() {
  # kill processes
  processes="layotto layotto_wasmtime etcd server client go"
  for key in ${processes}; do
    if killall $key; then
      echo "$key released"
    fi
  done

  # remove containers
  keywords="redis skywalking hangzhouzk minio"
  for key in ${keywords}; do
    if [ $(docker container ls -a | grep $key | wc -l) -gt 0 ]; then
      echo "Deleting containers of $key : " $to_delete
      docker rm -f $(docker container ls -a | grep $key | awk '{ print $1 }')
    fi
  done
}

release_resource

quickstarts=${quickstarts_in_advance}
if test "${GO_VERSION}" = "default"; then
  quickstarts=${quickstarts_in_default}
fi

# test quickstarts
echo "quickstarts contain ${quickstarts}"
for doc in ${quickstarts}; do
  echo "Start testing $doc......"

  #./mdx docs/en/start/state/start.md
  $(pwd)/etc/script/mdx $doc

  echo "End testing $doc......"
  release_resource
done
