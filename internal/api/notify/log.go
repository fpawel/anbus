package notify

import (
	"github.com/powerman/structlog"
)

func (x msg) Log(log *structlog.Logger) func(interface{}, ...interface{}) {
	return log.Debug
}
