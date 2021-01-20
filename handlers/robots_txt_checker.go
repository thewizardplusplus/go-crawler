package handlers

import (
	"github.com/go-log/log"
	crawler "github.com/thewizardplusplus/go-crawler"
	"github.com/thewizardplusplus/go-crawler/registers"
)

// RobotsTXTHandler ...
type RobotsTXTHandler struct {
	UserAgent         string
	RobotsTXTRegister registers.RobotsTXTRegister
	LinkHandler       crawler.LinkHandler
	Logger            log.Logger
}
