#
# Copyright 2021 Layotto Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

go build ./cmd/layotto
nohup redis-server --port 6380 &
nohup ./layotto start -c ./configs/config_redis.json &
cd sdk/go-sdk/test/runtime
go test -p 1 -v -run TestRedis ./...

cd ../../../../
kill -9 $(netstat -tlnp|grep 34904|awk '{print $7}'|awk -F '/' '{print $1}')
nohup etcd  &
nohup ./layotto start -c ./configs/runtime_config.json &
cd sdk/go-sdk/test/runtime
go test -p 1 -v -run TestEtcd ./...

# cd ../../../../
# go build ./cmd/layotto
# nohup bash /usr/share/zookeeper/bin/zkServer.sh start&
# bash /usr/share/zookeeper/bin/zkCli.sh
# create -s /sequencer|||app||MyKey "MyValue"
# quit
# nohup ./layotto start -c ./configs/config_zookeeper.json &
# cd sdk/go-sdk/test/runtime
# go test -p 1 -v -run TestZK ./...