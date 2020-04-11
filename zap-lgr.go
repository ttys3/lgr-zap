// zap bridge for github.com/go-pkgz/lgr v0.6.3
package lgr

import (
	"fmt"
	"io"
	stdlog "log"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var loggerMutex sync.RWMutex
var logger *zap.Logger

func InitZapLogger(dbg bool) func() {
	if Z() != nil {
		return func() {
			// no-ops
		}
	}
	// prod stackLevel := ErrorLevel
	// dev stackLevel = WarnLevel
	zapCfg := zap.NewDevelopmentConfig()
	//zapCfg := zap.NewProductionConfig()
	//zapCfg.Encoding = "console"
	if !dbg {
		zapCfg.DisableCaller = true
	}
	// if Development , stackLevel = WarnLevel
	zapCfg.Development = false
	zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	tmpLogger, _ := zapCfg.Build()
	loggerMutex.Lock()
	logger = tmpLogger.Named("[remark42]")
	logger = logger.WithOptions(zap.AddCallerSkip(1))
	//The default global logger used by zap.L() and zap.S() is a no-op logger.
	//To configure the global loggers, you must use ReplaceGlobals.
	zap.ReplaceGlobals(logger)
	loggerMutex.Unlock()
	if dbg {
		zap.S().Infof("debug enabled")
	}
	return func() {
		// flushes buffer, if any
		if err := logger.Sync(); err != nil {
			fmt.Printf("zap: Sync() failed with error: %s\n", err)
		}
	}
}

// return the default global *zap.Logger
func Z() *zap.Logger {
	loggerMutex.RLock()
	s:= logger
	loggerMutex.RUnlock()
	return s
}

// wrap the default global *zap.Logger as a go stdlib log
func NewStdLog() *stdlog.Logger {
	if Z() == nil {
		InitZapLogger(true)
	}
	return zap.NewStdLog(logger)
}

// wrap the default global *zap.Logger as a go stdlib log at specific log level
func NewStdLogAt(lvl string) *stdlog.Logger {
	if Z() == nil {
		InitZapLogger(true)
	}
	level := gopkgzlgrLevelToZapCoreLevel(lvl)
	l, err := zap.NewStdLogAt(logger, level)
	if err != nil {
		Errorf("bridge NewStdLogAt(%s) failed with error: %s", lvl, err)
		return nil
	}
	return l
}

//type gopkgzlgr interface {
//	Printf(format string, args ...interface{})
//	Print(line string)
//	Fatalf(format string, args ...interface{})
//	Setup(opts ...Option)
//}

// implement github.com/go-pkgz/rest/logger.Backend
// which required by required by backend/app/rest/api/rest.go
type GopkgzRestLogger interface {
	Logf(format string, args ...interface{})
}

var GopkgzRestLoggerBridge GopkgzRestLogger = &Logger{}

type Logger struct {}

func (z *Logger) Logf(format string, args ...interface{}) {
	Infof(format, args...)
}

// required by backend/app/rest/api/rest.go
func Default() GopkgzRestLogger {
	return GopkgzRestLoggerBridge
}

// bridged Option type for github.com/go-pkgz/lgr v0.6.3
// Option func type
type Option func(l *Logger)

// levels from github.com/go-pkgz/lgr/logger.go
var gopkgzlgrLevels = []string{"TRACE", "DEBUG", "INFO", "WARN", "ERROR", "PANIC", "FATAL"}

func gopkgzlgrLevelToZapCoreLevel(lvl string) zapcore.Level {
	switch lvl {
	case "TRACE":
		fallthrough
	case "DEBUG":
		return zapcore.DebugLevel
	case "INFO":
		return zapcore.InfoLevel
	case "WARN":
		return zapcore.WarnLevel
	case "ERROR":
		return zapcore.ErrorLevel
	case "PANIC":
		return zapcore.PanicLevel
	case "FATAL":
		return zapcore.FatalLevel
	}
	return zapcore.InfoLevel
}

func gopkgzlgrLevelToZapFunc(lvl string) func(template string, args ...interface{}) {
	switch lvl {
	case "TRACE":
		fallthrough
	case "DEBUG":
		return zap.S().Debugf
	case "INFO":
		return zap.S().Infof
	case "WARN":
		return zap.S().Warnf
	case "ERROR":
		return zap.S().Errorf
	case "PANIC":
		return zap.S().Panicf
	case "FATAL":
		return zap.S().Fatalf
	}
	return zap.S().Infof
}

// extractLevel parses messages with optional level prefix and returns level and the message with stripped level
func extractLevel(line string) (level, msg string) {
	for _, lv := range gopkgzlgrLevels {
		if strings.HasPrefix(line, lv) {
			return lv, strings.TrimSpace(line[len(lv):])
		}
		if strings.HasPrefix(line, "["+lv+"]") {
			return lv, strings.TrimSpace(line[len("["+lv+"]"):])
		}
	}
	return "INFO", line
}

// bridged Printf for github.com/go-pkgz/lgr v0.6.3
func Printf(format string, args ...interface{}) {
	// avoid %!s(MISSING)
	if len(args) == 0 {
		Print(format)
		return
	}
	lvl, msg := extractLevel(format)
	gopkgzlgrLevelToZapFunc(lvl)(msg, args)
}

// bridged Print for github.com/go-pkgz/lgr v0.6.3
func Print(line string) {
	lvl, msg := extractLevel(line)
	gopkgzlgrLevelToZapFunc(lvl)(msg)
}

// Setup default logger with options
// bridged Setup for github.com/go-pkgz/lgr v0.6.3
func Setup(opts ...Option) {
	// no-ops
}

// bridge non-format func wrap for zap sugar
// Debug in go-pkgz/lgr is not the same as in zap
//func Debug(line string) {
//	zap.S().Debug(line)
//}

func Info(line string) {
	zap.S().Info(line)
}

func Error(line string) {
	zap.S().Error(line)
}

func Warn(line string) {
	zap.S().Warn(line)
}

func Fatal(line string) {
	zap.S().Fatal(line)
}

func Panic(line string) {
	zap.S().Panic(line)
}

// bridge format func wrap for zap sugar
func Debugf(format string, args ...interface{}) {
	zap.S().Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	zap.S().Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	zap.S().Errorf(format, args...)
}

func Warnf(format string, args ...interface{}) {
	zap.S().Warnf(format, args...)
}

// bridged Fatalf for github.com/go-pkgz/lgr v0.6.3
func Fatalf(format string, args ...interface{}) {
	zap.S().Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	zap.S().Panicf(format, args...)
}

// other things

// Out sets out writer, stdout by default
func Out(w io.Writer) Option {
	return func(l *Logger) {
	}
}

// Err sets error writer, stderr by default
func Err(w io.Writer) Option {
	return func(l *Logger) {
	}
}

// Debug turn on dbg mode
func Debug(l *Logger) {
}

// Trace turn on trace + dbg mode
func Trace(l *Logger) {
}

// CallerDepth sets number of stack frame skipped for caller reporting, 0 by default
func CallerDepth(n int) Option {
	return func(l *Logger) {
	}
}

// Format sets output layout, overwrites all options for individual parts, i.e. Caller*, Msec and LevelBraces
func Format(f string) Option {
	return func(l *Logger) {
	}
}

// CallerFunc adds caller info with function name. Ignored if Format option used.
func CallerFunc(l *Logger) {
}

// CallerPkg adds caller's package name. Ignored if Format option used.
func CallerPkg(l *Logger) {
}

// LevelBraces surrounds level with [], i.e. [INFO]. Ignored if Format option used.
func LevelBraces(l *Logger) {
}

// CallerFile adds caller info with file, and line number. Ignored if Format option used.
func CallerFile(l *Logger) {
}

// Msec adds .msec to timestamp. Ignored if Format option used.
func Msec(l *Logger) {
}

// Secret sets list of substring to be hidden, i.e. replaced by "******"
// Useful to prevent passwords or other sensitive tokens to be logged.
func Secret(vals ...string) Option {
	return func(l *Logger) {
	}
}
