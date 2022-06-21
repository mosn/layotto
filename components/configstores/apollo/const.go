/*
 * Copyright 2021 Layotto Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package apollo

const (
	defaultTagsNamespace       = "sidecar_config_tags"
	defaultDelimiter           = "@$"
	defaultNamespace           = "application"
	defaultEnv                 = "DEV"
	defaultTimeoutWhenResponse = 2000
	defaultIsBackupConfig      = true
	configKeyAppId             = "app_id"
	setUrlTpl                  = "%v/openapi/v1/envs/%v/apps/%v/clusters/%v/namespaces/%v/items/%v"
	commitUrlTpl               = "%v/openapi/v1/envs/%v/apps/%v/clusters/%v/namespaces/%v/releases"
	deleteUrlTpl               = "%v/openapi/v1/envs/%v/apps/%v/clusters/%v/namespaces/%v/items/%v"
	createNamespaceUrlTpl      = "%v/openapi/v1/apps/%v/appnamespaces"
)
