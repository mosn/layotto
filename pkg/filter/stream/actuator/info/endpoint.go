package info

import (
	"context"
	"github.com/layotto/layotto/pkg/filter/stream/actuator"
	"mosn.io/pkg/log"
)

func init() {
	actuator.GetDefault().AddEndpoint("info", NewEndpoint())
}

var infoContributors = make(map[string]Contributor)

type Endpoint struct {
}

func NewEndpoint() *Endpoint {
	return &Endpoint{}
}

func (e *Endpoint) Handle(ctx context.Context, params actuator.ParamsScanner) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	var resultErr error = nil
	contributors := infoContributors
	for k, c := range contributors {
		cinfo, err := c.GetInfo()
		if err != nil {
			log.DefaultLogger.Errorf("[actuator][info] Error when GetInfo.Contributor:%v,error:%v", k, err)
			result[k] = err.Error()
			resultErr = err
		} else {
			result[k] = cinfo
		}
	}
	return result, resultErr
}

// AddInfoContributor register info.Contributor.It's not concurrent-safe,so please invoke it ONLY in init method.
func AddInfoContributor(name string, c Contributor) {
	if c == nil {
		return
	}
	infoContributors[name] = c
}

// AddInfoContributorFunc register info.Contributor.It's not concurrent-safe,so please invoke it ONLY in init method.
func AddInfoContributorFunc(name string, f func() (interface{}, error)) {
	AddInfoContributor(name, ContributorAdapter(f))
}
