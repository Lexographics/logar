package gormlogger

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sadk.dev/logar/models"
)

type loggerImpl struct {
	writer *writer
	logger.Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
	stackTraceSkip                      int
}

func (l *loggerImpl) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

func (l *loggerImpl) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.writer.Printf(ctx, models.Severity_Info, l.infoStr+msg, append([]interface{}{l.fileWithLineNum()}, data...)...)
	}
}

func (l *loggerImpl) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.writer.Printf(ctx, models.Severity_Warning, l.warnStr+msg, append([]interface{}{l.fileWithLineNum()}, data...)...)
	}
}

func (l *loggerImpl) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.writer.Printf(ctx, models.Severity_Error, l.errStr+msg, append([]interface{}{l.fileWithLineNum()}, data...)...)
	}
}

func (l *loggerImpl) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.writer.Printf(ctx, models.Severity_Error, l.traceErrStr, l.fileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.writer.Printf(ctx, models.Severity_Error, l.traceErrStr, l.fileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.writer.Printf(ctx, models.Severity_Warning, l.traceWarnStr, l.fileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.writer.Printf(ctx, models.Severity_Warning, l.traceWarnStr, l.fileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == logger.Info:
		sql, rows := fc()
		if rows == -1 {
			l.writer.Printf(ctx, models.Severity_Info, l.traceStr, l.fileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.writer.Printf(ctx, models.Severity_Info, l.traceStr, l.fileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

// ParamsFilter filter params
func (l *loggerImpl) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.Config.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}

func (l loggerImpl) fileWithLineNum() string {
	skip := 0

	pcs := [13]uintptr{}
	// the third caller usually from gorm internal
	len := runtime.Callers(3, pcs[:])
	frames := runtime.CallersFrames(pcs[:len])
	for i := 0; i < len; i++ {
		// second return value is "more", not "ok"
		frame, _ := frames.Next()
		if (!strings.Contains(frame.File, "gorm.io") ||
			strings.HasSuffix(frame.File, "_test.go")) && !strings.HasSuffix(frame.File, ".gen.go") {

			if skip >= l.stackTraceSkip {
				return string(strconv.AppendInt(append([]byte(frame.File), ':'), int64(frame.Line), 10))
			}
			skip++
		}
	}

	return ""
}
