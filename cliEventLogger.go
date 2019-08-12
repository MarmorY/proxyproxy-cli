package main

import (
	"github.com/Neothorn23/proxyproxy"
	"github.com/apex/log"
)

type cliEventLogger struct {
	logger *log.Entry
}

func (l *cliEventLogger) OnProxyEvent(t proxyproxy.ProxyEventType, pc *proxyproxy.ProxyCommunication) {
	if t == proxyproxy.EventCreatingConnection {
		l.logger = l.logger.WithFields(log.Fields{"Id": pc.GetID(), "src": pc.GetClientAddr()})
	}
	switch t {

	case proxyproxy.EventCreatingConnection:
		l.logger.Info(proxyproxy.EventText(t))
	case proxyproxy.EventNtlmAuthRequestDetected:
		l.logger.Info(proxyproxy.EventText(t))
	case proxyproxy.EventProcessingRequest:
		l.logger.Infof("Processing request: %s %s", pc.GetCurrentRequest().Method, pc.GetCurrentRequest().RequestURI)
	case proxyproxy.EventConnectionClosed:
		l.logger.Info(proxyproxy.EventText(t))
	default:
		l.logger.Debugf("%s: %v", proxyproxy.EventText(t), pc)
	}

}

//NewCliEventLogger creates a new cli logger wich implements ProxyEventListener
func NewCliEventLogger(logger *log.Entry) proxyproxy.ProxyEventListener {
	return &cliEventLogger{
		logger: logger,
	}
}
