package service

import "github.com/quickfixgo/quickfix"

type RouterService interface {
	// Router() (string, string, quickfix.MessageRoute)
	RegisterRouters(fn func(beginString string, msgType string, router quickfix.MessageRoute))
}
