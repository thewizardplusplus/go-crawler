package extractors

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewDelayedExtractor(test *testing.T) {
	sleeper := new(MockSleeper)
	linkExtractor := new(MockLinkExtractor)
	got := NewDelayedExtractor(100*time.Millisecond, sleeper.Sleep, linkExtractor)

	mock.AssertExpectationsForObjects(test, sleeper, linkExtractor)
	require.NotNil(test, got)
	assert.Equal(test, 100*time.Millisecond, got.minimalDelay)
	// don't use the reflect.Value.Pointer() method for this check; see details:
	// * https://golang.org/pkg/reflect/#Value.Pointer
	// * https://stackoverflow.com/a/9644797
	assert.NotNil(test, got.sleeper)
	assert.Equal(test, linkExtractor, got.linkExtractor)
}
