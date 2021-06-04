package health

type Indicator interface {
	Report() (string, map[string]interface{})
}

type IndicatorAdapter func() (string, map[string]interface{})

func (ca IndicatorAdapter) Report() (string, map[string]interface{}) {
	return ca()
}

// Status is the enumeration value of component health status.
type Status = string

var (
	// INIT means it is starting
	INIT = Status("INIT")
	// UP means it is healthy
	UP = Status("UP")
	// DOWN means it is unhealthy
	DOWN = Status("DOWN")
)

// Details hold additional contextual details about the health of a component.
type Details = map[string]interface{}

func NewDetails() Details {
	m := make(map[string]interface{})
	return Details(m)
}

// Health carries information about the health of a component.
// Details are optional.
type Health struct {
	Status  Status  `json:"status"`
	Details Details `json:"details,omitempty"`
}

func NewHealth(status Status) Health {
	return Health{
		Status:  status,
		Details: NewDetails(),
	}
}

// SetDetail sets a message v into the health details, indexed by k.
// Note that the previous message of k, if exists, will be overriden.
// v MUST be a valid json marshable type, otherwise runtime panic or
// error occurs which fails the actuator health API.
func (h *Health) SetDetail(k string, v interface{}) {
	if h == nil {
		return
	}
	h.Details[k] = v
}

// GetDetail returns the detailed message indexed by k.
func (h *Health) GetDetail(k string) interface{} {
	if h == nil {
		return nil
	}
	return h.Details[k]
}
