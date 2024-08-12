package service

import (
	"context"

	"github.com/phimaker/waanx-fix-simpler/internal/config"
	"github.com/phimaker/waanx-fix-simpler/internal/infrastructure/fix"
	"github.com/phimaker/waanx-fix-simpler/internal/logger"
	"github.com/quickfixgo/fix44/logon"
	"github.com/quickfixgo/quickfix"
)

type FixService interface {
	Start(ctx context.Context)
	RegisterRouters(ctx context.Context)
	Stop()

	OnLoggedOn() <-chan quickfix.SessionID
}

type fixServiceImpl struct {
	app    fix.FixApplication
	client *fix.Client

	heartbeatSrv    HeartbeatService
	securityListSrv SecurityListService

	loggedOnCh chan quickfix.SessionID
}

func NewFIXService(
	cfg *config.Fix,
	heartbeatSrv HeartbeatService,
	securityListSrv SecurityListService,
) (FixService, error) {
	logger.Infof("Creating FIX service with config: %+v", cfg)
	app, err := fix.NewApplication(
		fix.WithUsername(cfg.Username),
		fix.WithPassword(cfg.Password),
	)
	if err != nil {
		logger.Fatalf("error creating application: %w", err)
	}

	client, err := fix.NewClient(
		cfg.ConfigPath,
		app,
	)
	if err != nil {
		logger.Fatalf("error creating client: %w", err)
	}

	return &fixServiceImpl{
		app:             app,
		client:          client,
		securityListSrv: securityListSrv,
		heartbeatSrv:    heartbeatSrv,
		loggedOnCh:      make(chan quickfix.SessionID, 10),
	}, nil
}

func (s *fixServiceImpl) RegisterRouters(ctx context.Context) {
	s.app.AddRouter(logon.Route(func(msg logon.Logon, sessionID quickfix.SessionID) quickfix.MessageRejectError {
		s.loggedOnCh <- sessionID
		return nil
	}))

	s.heartbeatSrv.RegisterRouters(s.app.AddRouter)
	s.securityListSrv.RegisterRouters(s.app.AddRouter)
}

func (s *fixServiceImpl) Start(ctx context.Context) {
	logger.Info("Starting FIX client")
	if err := s.client.Start(); err != nil {
		logger.Error("Error starting FIX client: ", err)
		return
	}
}

func (s *fixServiceImpl) Stop() {
	logger.Info("Stopping FIX client")
	s.client.Stop()
	close(s.loggedOnCh)
}

func (s *fixServiceImpl) OnLoggedOn() <-chan quickfix.SessionID {
	return s.loggedOnCh
}
