quickstarts="docs/en/start/state/start.md
  docs/zh/start/state/start.md
  docs/en/start/pubsub/start.md
  docs/zh/start/pubsub/start.md
  docs/en/start/lock/start.md
  docs/zh/start/lock/start.md
  docs/en/start/rpc/helloworld.md
  docs/zh/start/rpc/helloworld.md
  docs/en/start/rpc/dubbo_json_rpc.md
  docs/zh/start/rpc/dubbo_json_rpc.md
  docs/zh/start/file/minio.md
  docs/en/start/api_plugin/helloworld.md
  docs/zh/start/api_plugin/helloworld.md
  docs/zh/start/actuator/start.md
  docs/en/start/actuator/start.md
  docs/zh/start/trace/trace.md
  docs/en/start/trace/trace.md
  docs/zh/start/trace/skywalking.md
  docs/en/start/wasm/start.md
  docs/zh/start/wasm/start.md
"

quickstarts=("docs/zh/start/sequencer/start.md")

export projectpath=$(pwd)
export project_path=$(pwd)

# release all resources
release_resource() {
  killall layotto
  killall etcd
  # remove all the docker containers
  if [ $(docker ps -a -q | wc -l) -gt 0 ]; then
    docker rm -f $(docker ps -a -q)
  fi
}

release_resource

# download etcd
if [ "$(uname)" == "Darwin" ]; then
  # Mac OS X
  sh etc/script/download_etcd_mac.sh
elif [ "$(expr substr $(uname -s) 1 5)" == "Linux" ]; then
  # GNU/Linux
  sh etc/script/download_etcd_linux.sh
else
  # Windows or other OS
  echo "Your OS is not supported!"
  exit 1
fi

# test quickstarts
for doc in ${quickstarts}; do
  echo "Start testing $doc......"

  #./mdsh.sh docs/en/start/state/start.md
  ./mdsh.sh $doc

  release_resource
done
