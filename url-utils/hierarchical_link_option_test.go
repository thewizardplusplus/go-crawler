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

func TestWithMaximalHierarchyDepth(test *testing.T) {
	var config HierarchicalLinkConfig
	option := WithMaximalHierarchyDepth(23)
	option(&config)

	assert.Equal(test, 23, config.maximalHierarchyDepth)
}
