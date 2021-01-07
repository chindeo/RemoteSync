package logging

import (
	"github.com/snowlyg/RemoteSync/utils"
	"path/filepath"
	"sync"
)

var remoteLogger *Logger
var locLogger *Logger
var usertypeLogger *Logger
var commonLogger *Logger

//var WorkDir string

func GetRemoteLogger() *Logger {
	var single sync.Once
	single.Do(
		func() {
			workDir := getWorkDir()
			remoteLogger = NewLogger(&Options{
				Rolling:     DAILY,
				TimesFormat: TIMESECOND,
			}, filepath.Join(workDir, "./logs/rest.log"))
			remoteLogger.SetLogPrefix("log_prefix")
		})

	return remoteLogger
}

func GetLocLogger() *Logger {
	var single sync.Once
	single.Do(
		func() {
			workDir := getWorkDir()
			locLogger = NewLogger(&Options{
				Rolling:     DAILY,
				TimesFormat: TIMESECOND,
			}, filepath.Join(workDir, "./logs/loc.log"))
			locLogger.SetLogPrefix("log_prefix")
		})

	return locLogger
}

func GetUserTypeLogger() *Logger {
	var single sync.Once
	single.Do(
		func() {
			workDir := getWorkDir()
			usertypeLogger = NewLogger(&Options{
				Rolling:     DAILY,
				TimesFormat: TIMESECOND,
			}, filepath.Join(workDir, "./logs/device.log"))
			usertypeLogger.SetLogPrefix("log_prefix")
		})

	return usertypeLogger
}

func GetCommonLogger() *Logger {
	var single sync.Once
	single.Do(
		func() {
			workDir := getWorkDir()
			commonLogger = NewLogger(&Options{
				Rolling:     DAILY,
				TimesFormat: TIMESECOND,
			}, filepath.Join(workDir, "./logs/common.log"))
			commonLogger.SetLogPrefix("log_prefix")
		})

	return commonLogger
}

func getWorkDir() string {
	return utils.Config.Outdir
}
