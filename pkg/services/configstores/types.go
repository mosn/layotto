package configstores

// StoreConfig wraps configuration for a store implementation
type StoreConfig struct {
	StoreName string            `json:"store_name"`
	Address   []string          `json:"address"`
	TimeOut   string            `json:"timeout"`
	Metadata  map[string]string `json:"metadata"`
}

// GetRequest is the object describing a get configuration request
type GetRequest struct {
	AppId           string
	Group           string
	Label           string
	Keys            []string
	Metadata        map[string]string
}

// SetRequest is the object describing a save configuration request
type SetRequest struct {
	StoreName string
	AppId string
	Items []*ConfigurationItem
}

// ConfigurationItem represents a configuration item with key, content and other information.
type ConfigurationItem struct {
	Key      string
	Content  string
	Group    string
	Label    string
	Tags     map[string]string
	Metadata map[string]string
}

// DeleteRequest is the object describing a delete configuration request
type DeleteRequest struct {
	AppId    string
	Group    string
	Label    string
	Keys     []string
	Metadata map[string]string
}

// SubscribeReq is the object describing a subscription request
type SubscribeReq struct {
	AppId string
	Group string
	Label string
	Keys []string
	Metadata map[string]string
}

// SubscribeResp is the object describing a response for subscription
type SubscribeResp struct {
	StoreName string
	AppId string
	Items []*ConfigurationItem
}

