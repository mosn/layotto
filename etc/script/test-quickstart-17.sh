#!/bin/bash -e

set -e

# won't test:
# docs/en/start/configuration/start-apollo.md
# docs/zh/start/configuration/start-apollo.md
# because the github workflow can not connect to the apollo server due to the great firewall

quickstarts="docs/en/start/rpc/dubbo_json_rpc.md
docs/zh/start/rpc/dubbo_json_rpc.md"


# download mdx
if ! test -e $(pwd)/etc/script/mdx
then
  curl -o $(pwd)/etc/script/mdx https://raw.githubusercontent.com/seeflood/mdx/main/mdx
fi
chmod +x $(pwd)/etc/script/mdx


# release all resources
release_resource() {
  if killall layotto; then
    echo "layotto released"
  fi
  if killall etcd; then
    echo "etcd released"
  fi
  if killall go; then
    echo "golang processes released"
  fi

  # remove all the docker containers
  if [ $(docker ps -a -q | wc -l) -gt 0 ]; then
    docker rm -f $(docker ps -a -q)
  fi
}

release_resource

# download etcd
sh etc/script/download_etcd.sh

# test quickstarts
for doc in ${quickstarts}; do
  echo "Start testing $doc......"

  #./mdx docs/en/start/state/start.md
  $(pwd)/etc/script/mdx $doc

  echo "End testing $doc......"
  release_resource
done
