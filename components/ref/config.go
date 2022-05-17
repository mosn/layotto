package ref

//RefConfig is ref json config
type RefConfig struct {
	ComponentRef     []*ComponentRef
	SecretRef        []*RefItem
	ConfigurationRef []*RefItem
}

type ComponentRef struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

type RefItem struct {
	Name string `json:"name"`
	Key  string `json:"key"`
}
