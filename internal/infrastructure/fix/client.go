package fix

import (
	"fmt"
	"os"
	"strconv"

	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgox/zaplog"
	"go.uber.org/zap/zapcore"
)

// Client holds the FIX initiator and manages its lifecycle.
type Client struct {
	Initiator   *quickfix.Initiator
	application quickfix.Application
}

// NewClient creates a new FIX Client with the specified configuration file.
func NewClient(cfgFileName string, app quickfix.Application) (*Client, error) {
	// Open configuration file
	cfg, err := os.Open(cfgFileName)
	if err != nil {
		return nil, fmt.Errorf("error opening config file(%s): %w", cfgFileName, err)
	}
	defer cfg.Close()

	// Parse settings from the configuration file
	settings, err := quickfix.ParseSettings(cfg)
	if err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	if err != nil {
		return nil, fmt.Errorf("error creating application: %w", err)
	}

	// Create message store factory
	// storeFactory := file.NewStoreFactory(settings)
	storeFactory := quickfix.NewMemoryStoreFactory()

	// logFactory, err := quickfix.NewFileLogFactory(settings)

	logLevel, err := strconv.Atoi(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = 0
	}
	logFactory, err := zaplog.NewZapLogFactory(
		settings,
		zaplog.WithConsoleLogLevel(zapcore.Level(logLevel)),
		zaplog.WithExtension(zaplog.LogExtension_Log),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating log factory: %w", err)
	}

	// logger.Fatal("logFactory: ", logFactory)
	// Create the FIX initiator
	initiator, err := quickfix.NewInitiator(app, storeFactory, settings, logFactory)
	if err != nil {
		return nil, err
	}

	// Return the newly created client
	return &Client{
		Initiator:   initiator,
		application: app,
	}, nil
}

// Start begins the FIX session managed by the initiator.
func (c *Client) Start() error {
	return c.Initiator.Start()
}

// Stop ends the FIX session managed by the initiator.
func (c *Client) Stop() {
	c.Initiator.Stop()
}
