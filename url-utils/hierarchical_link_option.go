package urlutils

// HierarchicalLinkConfig ...
type HierarchicalLinkConfig struct {
	sanitizeBaseLink      LinkSanitizing
	maximalHierarchyDepth int
}

// HierarchicalLinkOption ...
type HierarchicalLinkOption func(config *HierarchicalLinkConfig)

// SanitizeBaseLink ...
func SanitizeBaseLink(sanitize LinkSanitizing) HierarchicalLinkOption {
	return func(config *HierarchicalLinkConfig) {
		config.sanitizeBaseLink = sanitize
	}
}

// WithMaximalHierarchyDepth ...
func WithMaximalHierarchyDepth(depth int) HierarchicalLinkOption {
	return func(config *HierarchicalLinkConfig) {
		config.maximalHierarchyDepth = depth
	}
}
