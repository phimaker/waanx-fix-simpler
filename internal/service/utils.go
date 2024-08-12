package service

import (
	"github.com/phimaker/waanx-fix-simpler/internal/logger"
	"github.com/quickfixgo/quickfix"
)

func useExactValueIgnoreError[T any](fn func() (T, quickfix.MessageRejectError)) T {
	v, err := fn()
	if err != nil {
		logger.Debugf("Error getting value: %v", err)
	}
	return v
}
