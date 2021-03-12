package registers

import (
	"sync"
)

// LinkGenerator ...
type LinkGenerator interface {
	GenerateLinks(baseLink string) ([]string, error)
}

// SitemapRegister ...
type SitemapRegister struct {
	linkGenerator LinkGenerator

	registeredSitemaps *sync.Map
}
