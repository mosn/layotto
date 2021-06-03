package apollo

const (
	storename                  = "apollo"
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
