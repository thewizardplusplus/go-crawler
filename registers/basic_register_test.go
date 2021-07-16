package registers

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBasicRegister(test *testing.T) {
	got := NewBasicRegister()

	assert.Equal(test, new(sync.Map), got.registeredValues)
}
