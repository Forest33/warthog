// Package logger wrapper for zerolog
package logger

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// ZeroConfig logger settings
type ZeroConfig struct {
	Level             string `envconfig:"LEVEL"`
	TimeFieldFormat   string `envconfig:"TIME_FIELD_FORMAT"`
	PrettyPrint       bool   `envconfig:"PRETTY_PRINT"`
	RedirectStdLogger bool   `envconfig:"REDIRECT_STD_LOGGER"`
	DisableSampling   bool   `envconfig:"DISABLE_SAMPLING"`
	ErrorStack        bool   `envconfig:"ERROR_STACK"`
	ShowCaller        bool   `envconfig:"SHOW_CALLER"`
}

// Zerolog object capable of interacting with Zerolog
type Zerolog struct {
	zero              zerolog.Logger
	zeroErr           zerolog.Logger
	level             string
	prettyPrint       bool
	redirectSTDLogger bool
	rootInitialized   bool
	showCaller        bool
}

// NewDefaultZerolog creates Zerolog with default settings
func NewDefaultZerolog() *Zerolog {
	zerolog.SetGlobalLevel(getZerologLevel(defaultZeroCfg.Level))
	zerolog.DisableSampling(true)
	zerolog.TimeFieldFormat = defaultZeroCfg.TimeFieldFormat
	if defaultZeroCfg.ErrorStack {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	}

	var logger Zerolog
	logger.level = defaultZeroCfg.Level
	logger.prettyPrint = defaultZeroCfg.PrettyPrint
	logger.showCaller = defaultZeroCfg.ShowCaller
	logger.compileLogger()

	return &logger
}

var defaultZeroCfg = ZeroConfig{
	Level:           "debug",
	TimeFieldFormat: time.RFC3339,
	PrettyPrint:     true,
	ErrorStack:      false,
	ShowCaller:      false,
}

// NewZerolog creates a new Zerolog
func NewZerolog(config ZeroConfig) *Zerolog {
	zerolog.SetGlobalLevel(getZerologLevel(config.Level))
	zerolog.DisableSampling(config.DisableSampling)
	zerolog.TimeFieldFormat = config.TimeFieldFormat
	if config.ErrorStack {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	}

	var logger Zerolog
	logger.level = config.Level
	logger.prettyPrint = config.PrettyPrint
	logger.redirectSTDLogger = config.RedirectStdLogger
	logger.showCaller = config.ShowCaller
	logger.compileLogger()

	return &logger
}

// Debug starts a new message with debug level
func (l *Zerolog) Debug() *zerolog.Event {
	return l.zero.Debug()
}

// Info starts a new message with info level
func (l *Zerolog) Info() *zerolog.Event {
	return l.zero.Info()
}

// Error starts a new message with error level
func (l *Zerolog) Error() *zerolog.Event {
	return l.zeroErr.Error()
}

// Warn starts a new message with warn level
func (l *Zerolog) Warn() *zerolog.Event {
	return l.zeroErr.Warn()
}

// Panic starts a new message with panic level
func (l *Zerolog) Panic() *zerolog.Event {
	return l.zeroErr.Panic()
}

// With creates a child logger with the field added to its context
func (l *Zerolog) With() zerolog.Context {
	return l.zero.With()
}

// Fatal sends the event with fatal level
func (l *Zerolog) Fatal(v ...interface{}) {
	l.zeroErr.Fatal().Msgf("%v", v)
}

// Fatalf sends the event with formatted msg with fatal level
func (l *Zerolog) Fatalf(format string, v ...interface{}) {
	l.zeroErr.Fatal().Msgf(format, v)
}

// Print sends the event with debug level
func (l *Zerolog) Print(v ...interface{}) {
	l.zero.Debug().Msgf("%v", v)
}

// Printf sends the event with formatted msg with debug level
func (l *Zerolog) Printf(format string, v ...interface{}) {
	l.zero.Debug().Msgf(format, v)
}

func (l *Zerolog) initRootLogger() {
	l.zero = zerolog.New(os.Stdout).Level(getZerologLevel(l.level))
	l.zeroErr = zerolog.New(os.Stderr).Level(getZerologLevel(l.level))
	l.rootInitialized = true
}

func (l *Zerolog) compileLogger() {
	if !l.rootInitialized {
		l.initRootLogger()
	}

	if l.redirectSTDLogger {
		l.setLogOutputToZerolog()
	}

	l.initDefaultFields()

	if l.prettyPrint {
		l.addPrettyPrint()
	}
}

func (l *Zerolog) initDefaultFields() {
	l.zero = l.zero.With().Timestamp().Logger()
	l.zeroErr = l.zeroErr.With().Timestamp().Logger()
	if l.showCaller {
		l.zero = l.zero.With().Caller().Logger()
		l.zeroErr = l.zero.With().Caller().Logger()
	}
}

func (l *Zerolog) addPrettyPrint() {
	prettyStdout := zerolog.ConsoleWriter{Out: os.Stdout}
	prettyStderr := zerolog.ConsoleWriter{Out: os.Stderr}

	l.zero = l.zero.Output(prettyStdout)
	l.zeroErr = l.zeroErr.Output(prettyStderr)
}

func (l *Zerolog) setLogOutputToZerolog() {
	log.SetFlags(0)
	log.SetOutput(l.zero)
}

// SubLogger duplicates the current logger
func (l Zerolog) SubLogger(zero zerolog.Logger) *Zerolog {
	l.zero = zero.Output(os.Stdout)
	l.zeroErr = zero.Output(os.Stderr)

	if l.prettyPrint {
		l.addPrettyPrint()
	}

	return &l
}

func getZerologLevel(lvl string) zerolog.Level {
	switch strings.ToLower(lvl) {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	case "panic":
		return zerolog.PanicLevel
	case "disabled":
		return zerolog.Disabled
	}
	return zerolog.NoLevel
}
