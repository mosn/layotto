#!/usr/bin/env bash

docs=(
#  docs/en/start/configuration/start-apollo.md
#  docs/zh/start/configuration/start-apollo.md
 docs/en/start/state/start.md
 docs/zh/start/state/start.md
 docs/en/start/lock/start.md
 docs/zh/start/lock/start.md
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
)

export projectpath=$(pwd)
export project_path=$(pwd)

for doc in "${docs[@]}" ; do
  echo "Start testing $doc......"

  #./mdsh.sh docs/en/start/state/start.md
  #./mdsh.sh docs/zh/start/state/start.md
  ./mdsh.sh $doc

  # release all resources
  killall layotto
  if [ $(docker ps -a -q|wc -l) -gt 0 ]
  then
       docker rm -f $(docker ps -a -q)
  fi

done
