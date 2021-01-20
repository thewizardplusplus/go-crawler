package handlers

import (
	"context"
	"net/url"

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

// HandleLink ...
func (handler RobotsTXTHandler) HandleLink(link crawler.SourcedLink) {
	parsedLink, err := url.Parse(link.Link)
	if err != nil {
		handler.Logger.Logf("unable to parse the link: %s", err)
		return
	}

	robotsTXTData, err :=
		handler.RobotsTXTRegister.RegisterRobotsTXT(context.Background(), link.Link)
	if err != nil {
		handler.Logger.Logf("unable to register the robots.txt link: %s", err)
		return
	}

	group := robotsTXTData.FindGroup(handler.UserAgent)
	if !group.Test(parsedLink.Path) {
		return
	}

	handler.LinkHandler.HandleLink(link)
}
