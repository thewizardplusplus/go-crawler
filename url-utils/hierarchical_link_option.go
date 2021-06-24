package urlutils

// HierarchicalLinkConfig ...
type HierarchicalLinkConfig struct {
	sanitizeBaseLink      LinkSanitizing
	maximalHierarchyDepth int
}

// HierarchicalLinkOption ...
type HierarchicalLinkOption func(config *HierarchicalLinkConfig)
