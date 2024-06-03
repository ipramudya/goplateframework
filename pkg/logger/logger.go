package logger

import (
	"os"

	"github.com/goplateframework/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

type Log struct {
	conf   *config.Config
	logger *zerolog.Logger
}

var remapLogLevel = map[string]zerolog.Level{
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
	"fatal": zerolog.FatalLevel,
}

func Init(conf *config.Config) *Log {
	l := &Log{conf: conf}
	l.start()

	return l
}

func (l *Log) start() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(l.getLogLevel())

	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	l.logger = &logger
}

func (l *Log) getLogLevel() zerolog.Level {
	level, exists := remapLogLevel[l.conf.Logger.Level]
	if !exists {
		return zerolog.DebugLevel
	}
	return level
}

func (l *Log) Debug(msg string) {
	l.logger.Debug().Msg(msg)
}

func (l *Log) Debugf(template string, args ...interface{}) {
	l.logger.Debug().Msgf(template, args...)
}

func (l *Log) Info(msg string) {
	l.logger.Info().Msg(msg)
}

func (l *Log) Infof(template string, args ...interface{}) {
	l.logger.Info().Msgf(template, args...)
}

func (l *Log) Warn(msg string) {
	l.logger.Warn().Msg(msg)
}

func (l *Log) Warnf(template string, args ...interface{}) {
	l.logger.Warn().Msgf(template, args...)
}

func (l *Log) Error(msg string) {
	l.logger.Error().Msg(msg)
}

func (l *Log) Errorf(template string, args ...interface{}) {
	l.logger.Error().Msgf(template, args...)
}

func (l *Log) Fatal(msg string) {
	l.logger.Fatal().Msg(msg)
}

func (l *Log) Fatalf(template string, args ...interface{}) {
	l.logger.Fatal().Msgf(template, args...)
}
