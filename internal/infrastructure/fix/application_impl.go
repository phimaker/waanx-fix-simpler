package fix

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/phimaker/waanx-fix-simpler/internal/logger"
	"github.com/quickfixgo/enum"
	"github.com/quickfixgo/field"
	"github.com/quickfixgo/quickfix"
)

type FixApplication interface {
	quickfix.Application

	AddRouter(beginString string, msgType string, router quickfix.MessageRoute)
}

type fixApplicationOpt func(*fixApplicationImpl)

type fixApplicationImpl struct {
	username string
	password string
	router   *quickfix.MessageRouter

	logonHandler func(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError
}

func NewApplication(opts ...fixApplicationOpt) (FixApplication, error) {
	e := &fixApplicationImpl{
		username: "",
		password: "",
		router:   quickfix.NewMessageRouter(),
		logonHandler: func(msg *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {
			return nil
		},
	}

	for _, opt := range opts {
		opt(e)
	}

	return e, nil
}

func WithUsername(username string) fixApplicationOpt {
	return func(c *fixApplicationImpl) {
		c.username = username
	}
}

func WithPassword(password string) fixApplicationOpt {
	return func(c *fixApplicationImpl) {
		c.password = password
	}
}

// OnCreate implemented as part of Application interface
func (e *fixApplicationImpl) OnCreate(sessionID quickfix.SessionID) {
	logger.Infof("[ON_CREATE]: %s", sessionID.String())
}

// OnLogon implemented as part of Application interface
func (e *fixApplicationImpl) OnLogon(sessionID quickfix.SessionID) {
	logger.Infof("[LOGGED_ON]: %s", sessionID.String())
	e.logonHandler(&quickfix.Message{}, sessionID)
}

// OnLogout implemented as part of Application interface
func (e *fixApplicationImpl) OnLogout(sessionID quickfix.SessionID) {
	logger.Warnf("[LOGGED_OUT]: %s", sessionID.String())
}

func generateRawData() (string, error) {
	// Generate timestamp in milliseconds
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	// Generate nonce with 32 bytes of random data
	nonce := make([]byte, 32)
	_, err := rand.Read(nonce)
	if err != nil {
		return "", fmt.Errorf("failed to generate nonce: %v", err)
	}

	// Base64 encode the nonce
	encodedNonce := base64.StdEncoding.EncodeToString(nonce)

	// Combine timestamp and nonce with a period separator
	rawData := fmt.Sprintf("%d.%s", timestamp, encodedNonce)

	return rawData, nil
}

func generatePassword(rawData, clientSecret string) string {
	// Concatenate RawData and client secret
	combined := rawData + clientSecret

	// Compute SHA256 hash
	hash := sha256.Sum256([]byte(combined))

	// Base64 encode the hash
	encodedHash := base64.StdEncoding.EncodeToString(hash[:])

	return encodedHash
}

func generatewaanxAppSig(rawData, appSecret string) string {
	// Concatenate RawData and application secret
	combined := rawData + appSecret

	// Compute SHA256 hash
	hash := sha256.Sum256([]byte(combined))

	// Base64 encode the hash
	encodedHash := base64.StdEncoding.EncodeToString(hash[:])

	return encodedHash
}

// FromAdmin implemented as part of Application interface
func (e *fixApplicationImpl) FromAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	logger.Infof("[FROM_ADMIN] %s", msg.String())
	return nil
}

// ToAdmin implemented as part of Application interface
func (e *fixApplicationImpl) ToAdmin(msg *quickfix.Message, sessionID quickfix.SessionID) {
	logger.Infof("[TO_ADMIN] %s", msg.String())

	msgType, err := msg.MsgType()
	if err != nil {
		logger.Error(fmt.Sprintf("MsgType not found in message: %v", err))
		return
	}

	if msgType == "" {
		logger.Error("MsgType is empty")
		return
	}

	switch enum.MsgType(msgType) {

	case enum.MsgType_LOGON:
		if e.username != "" {
			msg.Header.Set(field.NewUsername(e.username))
		}

		// Generate Password
		if e.password != "" {
			// Generate RawData
			rawData, err := generateRawData()
			if err != nil {
				fmt.Printf("Error generating RawData: %v\n", err)
				return
			}

			password := generatePassword(rawData, e.password)
			msg.Header.Set(field.NewPassword(password))
			msg.Body.Set(field.NewRawData(rawData))
			msg.Body.Set(field.NewRawDataLength(len(rawData)))
		}

	case enum.MsgType_HEARTBEAT:
		logger.Infof("Sending Heartbeat to %s", sessionID.String())

	case enum.MsgType_TEST_REQUEST:
		logger.Infof("Sending TestRequest to %s", sessionID.String())
		// msg.SetField(field.NewTestReqID("test"))

	default:
		logger.Infof("MsgType %s not handled", msgType)
		return
	}
}

// ToApp implemented as part of Application interface
func (e *fixApplicationImpl) ToApp(msg *quickfix.Message, sessionID quickfix.SessionID) (err error) {
	logger.Infof("[TO_APP] %s", msg.String())
	return
}

// FromApp implemented as part of Application interface. This is the callback for all Application level messages from the counter party.
func (e *fixApplicationImpl) FromApp(msg *quickfix.Message, sessionID quickfix.SessionID) (reject quickfix.MessageRejectError) {
	msgType, err := msg.MsgType()
	if err != nil {
		logger.Error(fmt.Sprintf("MsgType not found in message: %v", err))
		return nil
	}
	if msgType != "y" && msgType != "W" {
		logger.Infof("[FROM_APP] %s", msg.String())
	} else {
		logger.Debugf("[FROM_APP] %s", msg.String())
	}
	return e.router.Route(msg, sessionID)
}

func (e *fixApplicationImpl) AddRouter(beginString string, msgType string, router quickfix.MessageRoute) {
	logger.Infof("Adding router for %s %#v", msgType, router)
	if msgType == "A" {
		e.logonHandler = router
	} else {
		e.router.AddRoute(beginString, msgType, router)
	}
}
