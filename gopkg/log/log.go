package log

import (
	"os"
	"strings"

	"github.com/lab46/monorepo/gopkg/errors"
	"go.uber.org/zap"
)

// logging library using uber zap

var (
	logger       *zap.Logger
	sugared      *zap.SugaredLogger
	currentLevel level
	fileOutput   string
)

type level int

// level of log
const (
	DebugLevel level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Log level
const (
	DebugLevelString = "debug"
	InfoLevelString  = "info"
	WarnLevelString  = "warn"
	ErrorLevelString = "error"
	FatalLevelString = "fatal"
)

func init() {
	SetLevel(InfoLevel)
}

func newZapConfig() zap.Config {
	config := zap.NewProductionConfig()
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.ErrorOutputPaths = []string{"stderr"}
	return config
}

func setLevel(config *zap.Config, l level) {
	switch l {
	case DebugLevel:
		config.Level.SetLevel(zap.DebugLevel)
	case InfoLevel:
		config.Level.SetLevel(zap.InfoLevel)
	case WarnLevel:
		config.Level.SetLevel(zap.WarnLevel)
	case ErrorLevel:
		config.Level.SetLevel(zap.ErrorLevel)
	case FatalLevel:
		config.Level.SetLevel(zap.FatalLevel)
	default:
		config.Level.SetLevel(zap.InfoLevel)
	}
}

func setLogger(log *zap.Logger) {
	logger = log
	logger.Sync()
	sugared = logger.Sugar()
	sugared.Sync()
}

// SetOutputToFile function
func SetOutputToFile(filename string) error {
	// attempted to create log file, if not created then the logger will return error
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if _, err = os.Create(filename); err != nil {
			return err
		}
	}
	config := newZapConfig()
	config.OutputPaths = []string{"stdout", filename}
	setLevel(&config, currentLevel)
	fileLogger, err := config.Build()
	if err != nil {
		return err
	}
	setLogger(fileLogger)
	fileOutput = filename
	return nil
}

// SetLevel will set level to logger and create a new logger based on level
func SetLevel(l level) {
	// set output to file if write to file exists
	if fileOutput != "" {
		currentLevel = l
		SetOutputToFile(fileOutput)
		return
	}
	config := newZapConfig()
	setLevel(&config, l)
	logger, _ := config.Build()
	setLogger(logger)
	currentLevel = l
}

// GetLevel return log level in string
func GetLevel() string {
	return levelToString(currentLevel)
}

// SetLevelString set level from string level
func SetLevelString(l string) {
	SetLevel(stringToLevel(l))
}

func stringToLevel(s string) level {
	switch strings.ToLower(s) {
	case DebugLevelString:
		return DebugLevel
	case InfoLevelString:
		return InfoLevel
	case WarnLevelString:
		return WarnLevel
	case ErrorLevelString:
		return ErrorLevel
	case FatalLevelString:
		return FatalLevel
	default:
		// TODO: make this more informative when happened
		return InfoLevel
	}
}

func levelToString(l level) string {
	switch l {
	case DebugLevel:
		return DebugLevelString
	case InfoLevel:
		return InfoLevelString
	case WarnLevel:
		return WarnLevelString
	case ErrorLevel:
		return ErrorLevelString
	case FatalLevel:
		return FatalLevelString
	default:
		return InfoLevelString
	}
}

// Debug log
func Debug(args ...interface{}) {
	sugared.Debug(args...)
}

// Debugf log
func Debugf(format string, args ...interface{}) {
	sugared.Debugf(format, args...)
}

// Debugw log
func Debugw(msg string, keyAndValues ...interface{}) {
	sugared.Debugw(msg, keyAndValues)
}

// Print log
func Print(args ...interface{}) {
	sugared.Info(args...)
}

// Println log
func Println(args ...interface{}) {
	sugared.Info(args...)
}

// Printf log
func Printf(format string, args ...interface{}) {
	sugared.Infof(format, args...)
}

// Printw log
func Printw(msg string, keyAndValues ...interface{}) {
	sugared.Infow(msg, keyAndValues...)
}

// Info log
func Info(args ...interface{}) {
	sugared.Info(args...)
}

// Infof log
func Infof(format string, args ...interface{}) {
	sugared.Infof(format, args...)
}

// Infow log
func Infow(msg string, keyAndValues ...interface{}) {
	sugared.Infow(msg, keyAndValues...)
}

// Warn log
func Warn(args ...interface{}) {
	sugared.Warn(args...)
}

// Warnf log
func Warnf(format string, args ...interface{}) {
	sugared.Warnf(format, args...)
}

// Warnw log
func Warnw(msg string, keyAndValues ...interface{}) {
	sugared.Warnw(msg, keyAndValues...)
}

// Error log
func Error(args ...interface{}) {
	sugared.Error(args...)
}

// Errorf log
func Errorf(format string, args ...interface{}) {
	sugared.Errorf(format, args...)
}

// Errorw log
func Errorw(msg string, keyAndValues ...interface{}) {
	sugared.Errorw(msg, keyAndValues...)
}

// Errors log log error detail from Errs
func Errors(err error) {
	var (
		errFields = make(errors.Fields)
		file      string
		line      int
	)
	switch err.(type) {
	case *errors.Errs:
		errs := err.(*errors.Errs)
		errFields = errs.GetFields()
		file, line = errs.GetFileAndLine()
		if len(errFields) == 0 {
			errFields = make(errors.Fields)
		}
	}
	if line != 0 {
		errFields["err_file"] = formatFilePath(file)
		errFields["err_line"] = line
	}
	intf := errFields.ToArrayInterface()
	sugared.With(intf...).Error(err.Error())
}

// Fatal log
func Fatal(args ...interface{}) {
	sugared.Fatal(args...)
}

// Fatalf log
func Fatalf(format string, args ...interface{}) {
	sugared.Fatalf(format, args...)
}

// Fatalw log
func Fatalw(format string, keyAndValues ...interface{}) {
	sugared.Fatalw(format, keyAndValues...)
}

// With log
func With(args ...interface{}) *zap.SugaredLogger {
	return sugared.With(args...)
}

func formatFilePath(f string) string {
	slash := strings.LastIndex(f, "/")
	return f[slash+1:]
}
