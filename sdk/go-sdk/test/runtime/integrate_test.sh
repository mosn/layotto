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

# fail fast
set -e

# start storage systems, e.g. redis, zk, etcd
nohup redis-server --port 6380 &
bash /usr/share/zookeeper/bin/zkServer.sh start
nohup bash /usr/share/zookeeper/bin/zkCli.sh < ./sdk/go-sdk/test/runtime/zkCreateZnode.sh
cd ..
nohup etcd &

# build and run Layotto
cd layotto
go build ./cmd/layotto
nohup ./layotto start -c ./configs/config_integrate_test.json &

# run integrate_test
cd sdk/go-sdk/test/runtime
go test -p 1 -v ./...

# run demos
cd ../../../../demo/configuration/common
go build -o client
names="etcd_config_demo"
for key in ${names}; do
    ./client -s $key
done

cd ../../state/common
go build -o client
names="redis_state_demo zookeeper_state_demo"
for key in ${names}; do
    ./client -s $key
done

cd ../../lock/common
go build -o client
names="redis_lock_demo etcd_lock_demo zookeeper_lock_demo"
for key in ${names}; do
    ./client -s $key
done

cd ../../sequencer/common
go build -o client
names="redis_sequencer_demo etcd_sequencer_demo zookeeper_sequencer_demo"
for key in ${names}; do
    ./client -s $key
done

cd ../../pubsub/server
names="redis_pub_subs_demo"
for key in ${names}; do
    cd ../server
    go build -o subscriber
    ./subscriber -s $key &
    cd ../client
    go build -o publisher
    ./publisher -s $key  
done

cd ../../secret/common
go build -o client
names="local_file_secret_demo"
for key in ${names}; do
    ./client -s $key
done