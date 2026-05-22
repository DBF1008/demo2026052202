package sys_log_hook

import (
	"go.uber.org/zap/zapcore"
)

func ZapLogHandler(entry zapcore.Entry) error {

	go func(paramEntry zapcore.Entry) {

	}(entry)
	return nil
}
