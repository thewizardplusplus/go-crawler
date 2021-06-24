package urlutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeBaseLink(test *testing.T) {
	var config HierarchicalLinkConfig
	option := SanitizeBaseLink(SanitizeLink)
	option(&config)

	assert.Equal(test, SanitizeLink, config.sanitizeBaseLink)
}
