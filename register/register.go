package register

import (
	mapset "github.com/deckarep/golang-set"
	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/sanitizing"
)

// LinkRegister ...
type LinkRegister struct {
	sanitizeLink sanitizing.LinkSanitizing
	logger       log.Logger

	registeredLinks mapset.Set
}
