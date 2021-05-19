package apollo

import (
	"github.com/layotto/layotto/pkg/actuator/health"
	testify "github.com/stretchr/testify/assert"
	"testing"
)

func TestGetHealthInitOrSuccess(t *testing.T) {
	assert := testify.New(t)

	hi := newHealthIndicator()
	h := hi.Report()
	assert.Equal(h.Status, health.INIT)
	hi.setStarted()
	h = hi.Report()
	assert.Equal(h.Status, health.UP)
}

func TestGetHealthError(t *testing.T) {
	assert := testify.New(t)

	hi := newHealthIndicator()
	hi.reportError("sub error")
	h := hi.Report()
	assert.Equal(h.Status, health.DOWN)
	assert.Equal(h.Details[reasonKey], "sub error")
}
