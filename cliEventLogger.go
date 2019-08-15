package main

import (
	"github.com/Neothorn23/proxyproxy"
	"github.com/apex/log"
)

type cliEventLogger struct {
	logger *log.Entry
}

func (l *cliEventLogger) OnProxyEvent(e *proxyproxy.ProxyEvent) {
	if e.EventType == proxyproxy.EventCreatingConnection {
		l.logger = l.logger.WithFields(log.Fields{"Id": e.ID, "src": e.ClientHost})
	}
	switch e.EventType {
	case proxyproxy.EventCreatingConnection:
		l.logger.Info(proxyproxy.EventText(e.EventType))
	case proxyproxy.EventNtlmAuthRequestDetected:
		l.logger.Info(proxyproxy.EventText(e.EventType))
	case proxyproxy.EventProcessingRequest:
		l.logger.Infof("Processing request: %s %s", e.Method, e.RequestURI)
	case proxyproxy.EventConnectionClosed:
		l.logger.Info(proxyproxy.EventText(e.EventType))
	default:
		l.logger.Debugf("%s: %v", proxyproxy.EventText(e.EventType), e)
	}

}

//NewCliEventLogger creates a new cli logger wich implements ProxyEventListener
func NewCliEventLogger(logger *log.Entry) proxyproxy.ProxyEventListener {
	return &cliEventLogger{
		logger: logger,
	}
}
