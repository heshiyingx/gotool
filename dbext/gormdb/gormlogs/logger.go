package gormlogs

import (
	"context"
	"errors"
	"fmt"
	"go.opentelemetry.io/otel/trace"
	logger2 "gorm.io/gorm/logger"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/gorm/utils"
)

// ErrRecordNotFound record not found error
//var ErrRecordNotFound = errors.New("record not found")

// Colors
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

//// LogLevel log level
//type LogLevel int
//
//const (
//	// Silent silent log level
//	Silent LogLevel = iota + 1
//	// Error error log level
//	Error
//	// Warn warn log level
//	Warn
//	// Info info log level
//	Info
//)

// Writer log writer interface
type Writer interface {
	Printf(string, ...interface{})
}

// Config logger config
type Config struct {
	SlowThreshold             time.Duration
	Colorful                  bool
	IgnoreRecordNotFoundError bool
	ParameterizedQueries      bool
	LogLevel                  logger2.LogLevel
}



var (
	// Discard logger will print any log to io.Discard
	Discard = New(log.New(io.Discard, "", log.LstdFlags), Config{})
	// Default Default logger
	Default = New(log.New(os.Stdout, "\r\n", log.LstdFlags), Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger2.Warn,
		IgnoreRecordNotFoundError: false,
		Colorful:                  true,
	})
	// Recorder logger records running SQL into a recorder instance
	Recorder = traceRecorder{Interface: Default, BeginAt: time.Now()}
)

// New initialize logger
func New(writer Writer, config Config) logger2.Interface {
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

	return &logger{
		Writer:       writer,
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

type logger struct {
	Writer
	Config
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// LogMode log mode
func (l *logger) LogMode(level logger2.LogLevel) logger2.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l *logger) Info(ctx context.Context, msg string, data ...interface{}) {

	if l.LogLevel >= logger2.Info {
		span := trace.SpanFromContext(ctx).SpanContext()
		hasTraceID := span.HasTraceID()
		if hasTraceID {
			msg = msg + "; traceID:" + span.TraceID().String()
		}
		l.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn print warn messages
func (l *logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger2.Warn {
		span := trace.SpanFromContext(ctx).SpanContext()
		hasTraceID := span.HasTraceID()
		if hasTraceID {
			msg = msg + "; traceID:" + span.TraceID().String()
		}
		l.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error print error messages
func (l *logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger2.Error {
		span := trace.SpanFromContext(ctx).SpanContext()
		hasTraceID := span.HasTraceID()
		if hasTraceID {
			msg = msg + "; traceID:" + span.TraceID().String()
		}
		l.Printf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace print sql message
//
//nolint:cyclop
func (l *logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger2.Silent {
		return
	}
	span := trace.SpanFromContext(ctx).SpanContext()
	hasTraceID := span.HasTraceID()
	traceIDStr := ""
	if hasTraceID {
		traceIDStr = "; traceID:" + span.TraceID().String()
	}

	elapsed := time.Since(begin)
	switch {
	case err != nil && l.LogLevel >= logger2.Error && (!errors.Is(err, logger2.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceErrStr+traceIDStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceErrStr+traceIDStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger2.Warn:
		sql, rows := fc()
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		if rows == -1 {
			l.Printf(l.traceWarnStr+traceIDStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceWarnStr+traceIDStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case l.LogLevel == logger2.Info:
		sql, rows := fc()
		if rows == -1 {
			l.Printf(l.traceStr+traceIDStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			l.Printf(l.traceStr+traceIDStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

// ParamsFilter filter params
func (l *logger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.Config.ParameterizedQueries {
		return sql, nil
	}
	return sql, params
}

type traceRecorder struct {
	logger2.Interface
	BeginAt      time.Time
	SQL          string
	RowsAffected int64
	Err          error
}

// New trace recorder
func (l *traceRecorder) New() *traceRecorder {
	return &traceRecorder{Interface: l.Interface, BeginAt: time.Now()}
}

// Trace implement logger interface
func (l *traceRecorder) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	l.BeginAt = begin
	l.SQL, l.RowsAffected = fc()
	l.Err = err
}
