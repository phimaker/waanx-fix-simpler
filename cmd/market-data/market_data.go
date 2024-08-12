package marketdata

import (
	"context"
	"os"
	"os/signal"

	"github.com/phimaker/waanx-fix-simpler/internal/config"
	"github.com/phimaker/waanx-fix-simpler/internal/logger"
	"github.com/phimaker/waanx-fix-simpler/internal/service"
	"github.com/quickfixgo/quickfix"
	"github.com/spf13/cobra"
)

const (
	usage = "marketdata"
	short = "Starts the market data service."
	long  = "Starts the market data service."
)

var (
	// Cmd is the executor command.
	Cmd = &cobra.Command{
		Use:     usage,
		Short:   short,
		Long:    long,
		Aliases: []string{"md", "market-data", "marketdata"},
		Example: "waanx-adapter marketdata -c config.yaml",
		RunE:    execute,
	}

	configPath string
	sessionID  quickfix.SessionID
)

func init() {
	Cmd.Flags().StringVarP(&configPath, "config", "c", "config.yaml", "path to the configuration file")
}

func execute(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)
	defer cancel()
	logger.InitLogger()
	cfg := config.GetConfig()

	heartbeatSrv := service.NewHeartbeatService()
	securityListSrv := service.NewSecurityListService()

	fixSrv, err := service.NewFIXService(cfg.Fix, heartbeatSrv, securityListSrv)
	if err != nil {
		logger.Fatalf("error creating FIX service: %w", err)
	}

	fixSrv.RegisterRouters(ctx)

	go fixSrv.Start(ctx)
	defer fixSrv.Stop()

	for {
		select {
		case sID := <-fixSrv.OnLoggedOn():
			sessionID = sID
			logger.Infof("Logged on: %s", sessionID)
			slrId, err := securityListSrv.SecurityListRequest(ctx, sessionID)
			if err != nil {
				logger.Errorf("Error sending security list request: %v", err)
			}
			logger.Infof("Sent SecurityListRequest with ID: %s", slrId)
		case <-ctx.Done():
			logger.Info("Shutting down market data service")
			return nil
		}
	}

}
