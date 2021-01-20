package checkers

import (
	"github.com/go-log/log"
	"github.com/thewizardplusplus/go-crawler/registers"
)

// RobotsTXTChecker ...
type RobotsTXTChecker struct {
	UserAgent         string
	RobotsTXTRegister registers.RobotsTXTRegister
	Logger            log.Logger
}
