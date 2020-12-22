package logging

import (
	"github.com/snowlyg/RemoteSync/utils"
	"path/filepath"
)

var Dbug *Logger
var Err *Logger
var Norm *Logger

func init() {
	Dbug = NewLogger(&Options{
		Rolling:     DAILY,
		TimesFormat: TIMESECOND,
	}, filepath.Join(utils.LogDir(), "./debug.log"))
	Dbug.SetLogPrefix("log_prefix")

	Err = NewLogger(&Options{
		Rolling:     DAILY,
		TimesFormat: TIMESECOND,
	}, filepath.Join(utils.LogDir(), "./error.log"))
	Err.SetLogPrefix("log_prefix")

	Norm = NewLogger(&Options{
		Rolling:     DAILY,
		TimesFormat: TIMESECOND,
	}, filepath.Join(utils.LogDir(), "./info.log"))
	Norm.SetLogPrefix("log_prefix")
}
