package logging

import (
	"github.com/snowlyg/RemoteSync/utils"
	"path/filepath"
)

var remoteLogger *Logger
var locLogger *Logger
var usertypeLogger *Logger
var commonLogger *Logger

//var WorkDir string

func GetRemoteLogger() *Logger {
	if remoteLogger != nil {
		return remoteLogger
	}
	workDir := getWorkDir()
	remoteLogger = NewLogger(&Options{
		Rolling:     DAILY,
		TimesFormat: TIMESECOND,
	}, filepath.Join(workDir, "./logs/remote.log"))
	remoteLogger.SetLogPrefix("log_prefix")
	return remoteLogger
}

func GetLocLogger() *Logger {
	if remoteLogger != nil {
		return locLogger
	}
	workDir := getWorkDir()
	locLogger = NewLogger(&Options{
		Rolling:     DAILY,
		TimesFormat: TIMESECOND,
	}, filepath.Join(workDir, "./logs/loc.log"))
	locLogger.SetLogPrefix("log_prefix")
	return locLogger
}

func GetUserTypeLogger() *Logger {
	if usertypeLogger != nil {
		return usertypeLogger
	}
	workDir := getWorkDir()
	usertypeLogger = NewLogger(&Options{
		Rolling:     DAILY,
		TimesFormat: TIMESECOND,
	}, filepath.Join(workDir, "./logs/user_type.log"))
	usertypeLogger.SetLogPrefix("log_prefix")
	return usertypeLogger
}

func GetCommonLogger() *Logger {
	if commonLogger != nil {
		return commonLogger
	}
	workDir := getWorkDir()
	commonLogger = NewLogger(&Options{
		Rolling:     DAILY,
		TimesFormat: TIMESECOND,
	}, filepath.Join(workDir, "./logs/common.log"))
	commonLogger.SetLogPrefix("log_prefix")
	return commonLogger
}

func getWorkDir() string {
	return utils.Config.Outdir
}
