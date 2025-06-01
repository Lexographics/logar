package gormlogger

import (
	"time"

	"gorm.io/gorm/logger"
	"sadk.dev/logar"
)

func New(lg logar.App, model logar.Model, category string, stackTraceSkip int) logger.Interface {
	return newInterface(
		newWriter(lg, model, category),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		}, stackTraceSkip,
	)
}

func newInterface(writer *writer, config logger.Config, stackTraceSkip int) logger.Interface {
	var (
		infoStr      = "%s\n[info] "
		warnStr      = "%s\n[warn] "
		errStr       = "%s\n[error] "
		traceStr     = "%s\n[%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s\n[%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s\n[%.3fms] [rows:%v] %s"
	)

	if config.Colorful {
		infoStr = Green + "%s\n" + Reset + Green + "[info] " + Reset
		warnStr = BlueBold + "%s\n" + Reset + Magenta + "[warn] " + Reset
		errStr = Magenta + "%s\n" + Reset + Red + "[error] " + Reset
		traceStr = Green + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
		traceWarnStr = Green + "%s " + Yellow + "%s\n" + Reset + RedBold + "[%.3fms] " + Yellow + "[rows:%v]" + Magenta + " %s" + Reset
		traceErrStr = RedBold + "%s " + MagentaBold + "%s\n" + Reset + Yellow + "[%.3fms] " + BlueBold + "[rows:%v]" + Reset + " %s"
	}

	return &loggerImpl{
		writer:         writer,
		Config:         config,
		infoStr:        infoStr,
		warnStr:        warnStr,
		errStr:         errStr,
		traceStr:       traceStr,
		traceWarnStr:   traceWarnStr,
		traceErrStr:    traceErrStr,
		stackTraceSkip: stackTraceSkip,
	}
}

const (
	ResetTag       = "[reset]"
	RedTag         = "[red]"
	GreenTag       = "[green]"
	YellowTag      = "[yellow]"
	BlueTag        = "[blue]"
	MagentaTag     = "[magenta]"
	CyanTag        = "[cyan]"
	WhiteTag       = "[white]"
	BlueBoldTag    = "[blue-bold]"
	MagentaBoldTag = "[magenta-bold]"
	RedBoldTag     = "[red-bold]"
	YellowBoldTag  = "[yellow-bold]"
)

const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)
