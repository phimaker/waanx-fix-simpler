package service

import (
	"context"
	"fmt"
	"time"

	"github.com/phimaker/waanx-fix-simpler/internal/logger"
	// "github.com/phimaker/waanx-fix-simpler/internal/domain"

	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/fix44/securitylist"
	"github.com/quickfixgo/fix44/securitylistrequest"
	"github.com/quickfixgo/quickfix"
)

type SecurityListService interface {
	RouterService

	SecurityListRequest(ctx context.Context, sessionID quickfix.SessionID) (string, error)
}

type securityListServiceImpl struct {
}

func NewSecurityListService() SecurityListService {
	return &securityListServiceImpl{}
}

func (srv *securityListServiceImpl) RegisterRouters(route func(beginString string, msgType string, router quickfix.MessageRoute)) {
	route(securitylist.Route(srv.OnSecurityList))
}

func (srv *securityListServiceImpl) OnSecurityList(msg securitylist.SecurityList, sessionID quickfix.SessionID) quickfix.MessageRejectError {
	logger.Infof("SecurityList: %v", msg)
	groups, err := msg.GetNoRelatedSym()
	if err != nil {
		logger.Errorf("Error getting NoMDEntries group: %v", err)
		return err
	}

	if groups.Len() == 0 {
		logger.Error("No MDEntries found")
		return quickfix.NewMessageRejectError("No MDEntries found", 0, nil)
	}

	for i := 0; i < groups.Len(); i++ {
		group := groups.Get(i)

		symbol := useExactValueIgnoreError(group.GetSymbol)
		securityId := useExactValueIgnoreError(group.GetSecurityID)
		logger.Infof("[SYMBOL:%s]", symbol)
		logger.Infof("[SecurityID:%s]", securityId)
	}
	return nil
}

func (srv *securityListServiceImpl) SecurityListRequest(ctx context.Context, sessionID quickfix.SessionID) (string, error) {
	return srv.sendSecurityListRequest(ctx, sessionID)
}

func (srv *securityListServiceImpl) sendSecurityListRequest(ctx context.Context, sessionID quickfix.SessionID) (string, error) {
	t := time.Now()
	reqID := fmt.Sprintf("SLR-%d", t.UnixNano())
	req := securitylistrequest.New(
		field.NewSecurityReqID(reqID),
		field.NewSecurityListRequestType(enum.SecurityListRequestType_ALL_SECURITIES),
	)
	logger.Infof("Request: %v\n", req.ToMessage())

	if err := quickfix.SendToTarget(req, sessionID); err != nil {
		return "", fmt.Errorf("Error sending market data request: %v", err)
	}

	return reqID, nil
}
