package snowflake

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetInstanceWM(t *testing.T) {
	machine := int64(1)
	machineRoom := int64(1)
	s := GetInstanceWM(machine, machineRoom)
	assert.NotNil(t, s)
}

func TestSingleton_NextID(t *testing.T) {
	s := GetInstanceWM(10, 10)
	id, err := s.NextID()
	if err != nil {
		t.Errorf("error: %v", err)
	}
	logger.Printf("id : %d", id)
	assert.NotNil(t, id)

}
