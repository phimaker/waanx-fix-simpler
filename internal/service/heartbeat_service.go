package service

import (
	"github.com/phimaker/waanx-fix-simpler/internal/logger"
	"github.com/quickfixgo/fix44/heartbeat"
	"github.com/quickfixgo/quickfix"
)

type HeartbeatService interface {
	RouterService
	OnHeartbeat(msg heartbeat.Heartbeat, sessionID quickfix.SessionID) quickfix.MessageRejectError
}

type heartbeatServiceImpl struct {
}

func NewHeartbeatService() HeartbeatService {
	return &heartbeatServiceImpl{}
}

func (s *heartbeatServiceImpl) RegisterRouters(route func(beginString string, msgType string, router quickfix.MessageRoute)) {
	route(heartbeat.Route(s.OnHeartbeat))
}

func (s *heartbeatServiceImpl) OnHeartbeat(msg heartbeat.Heartbeat, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	logger.Infof("Received heartbeat")
	return nil
}
