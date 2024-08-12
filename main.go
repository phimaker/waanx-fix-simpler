package main

import (
	"log"
	"os"

	"github.com/phimaker/waanx-fix-simpler/cmd"
	"github.com/phimaker/waanx-fix-simpler/internal/logger"
)

func main() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	logger.InitLogger(
		logger.WithCommonLogPath("logs/common.log"),
		logger.WithErrorLogPath("logs/error.log"),
		logger.WithMaxSize(100),
		logger.WithMaxBackups(10),
		logger.WithMaxAge(28),
		logger.WithCompress(true),
	)
	defer logger.Sync()
	logger.Infof("Starting server hostname: %s", hostname)
	logger.WithString("hostname", hostname)

	if err := cmd.Execute(); err != nil {
		os.Exit(0)
	}

}
