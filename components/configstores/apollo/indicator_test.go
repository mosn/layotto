package apollo

import (
	"testing"

	"github.com/layotto/components/pkg/common"

	testify "github.com/stretchr/testify/assert"
)

func TestGetHealthInitOrSuccess(t *testing.T) {
	assert := testify.New(t)

	hi := newHealthIndicator()
	v, _ := hi.Report()
	assert.Equal(v, common.INIT)
	hi.setStarted()
	h, _ := hi.Report()
	assert.Equal(h, common.UP)
}

func TestGetHealthError(t *testing.T) {
	assert := testify.New(t)

	hi := newHealthIndicator()
	hi.reportError("sub error")
	h, v := hi.Report()
	assert.Equal(h, common.DOWN)
	assert.Equal(v[reasonKey], "sub error")
}
