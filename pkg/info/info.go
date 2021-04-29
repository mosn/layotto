package info

// Runtimeinfo
type RuntimeInfo struct {
	Services ServiceInfo `json:"services"`
}

// 注册的服务信息
type ServiceInfo map[string]*ComponentInfo

// 组件信息
type ComponentInfo struct {
	// 启动时注册了哪些组件
	Registered []string `json:"registered"`
	// 启动时加载（初始化）了哪些组件
	Loaded []string `json:"loaded"`
}

func NewRuntimeInfo() *RuntimeInfo {
	return &RuntimeInfo{
		Services: ServiceInfo{},
	}
}

func (info *RuntimeInfo) AddService(service string) {
	info.Services[service] = &ComponentInfo{}
}

func (info *RuntimeInfo) RegisterComponent(service string, name string) {
	if c, ok := info.Services[service]; ok {
		c.Registered = append(c.Registered, name)
	}
}

func (info *RuntimeInfo) LoadComponent(service string, name string) {
	if c, ok := info.Services[service]; ok {
		c.Loaded = append(c.Loaded, name)
	}
}
