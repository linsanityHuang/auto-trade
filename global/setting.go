package global

import (
	"auto-trade/pkg/setting"
	"time"
)

var (
	// BigOneSetting bigone setting
	BigOneSetting *setting.BigOne
	// Timeout timeout
	Timeout = 1 * time.Second
	// Email email setting
	EmailSetting *setting.Email
)
