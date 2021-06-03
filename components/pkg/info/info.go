package info

// Runtimeinfo
type RuntimeInfo struct {
	Services ServiceInfo `json:"services"`
}

// ServiceInfo
type ServiceInfo map[string]*ComponentInfo

// ComponentInfo
type ComponentInfo struct {
	// Registered Component
	Registered []string `json:"registered"`
	// Loaded Component
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
