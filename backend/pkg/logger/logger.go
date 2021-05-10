package logger

import (
	"os"

	"github.com/inconshreveable/log15"
	"github.com/omeroid/wdc/backend/pkg/config"
)

// InitLogger :
func InitLogger(c config.Logger) {
	stackHandler := log15.CallerStackHandler("%+v", log15.StderrHandler)
	if c.LogJSON {
		s := log15.StreamHandler(os.Stderr, log15.JsonFormatEx(false, true))
		stackHandler = log15.CallerStackHandler("%+v", s)
	}

	lvlFilterHandler := log15.LvlFilterHandler(log15.LvlInfo, stackHandler)
	if c.Debug {
		lvlFilterHandler = log15.LvlFilterHandler(log15.LvlDebug, stackHandler)
	}
	log15.Root().SetHandler(lvlFilterHandler)
}
