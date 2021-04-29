/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package roundrobin

import (
	"github.com/zouyx/agollo/v4/env/config"
)

//RoundRobin 轮询调度
type RoundRobin struct {
}

//Load 负载均衡
func (r *RoundRobin) Load(servers map[string]*config.ServerInfo) *config.ServerInfo {
	var returnServer *config.ServerInfo
	for _, server := range servers {
		// if some node has down then select next node
		if server.IsDown {
			continue
		}
		returnServer = server
		break
	}
	return returnServer
}
