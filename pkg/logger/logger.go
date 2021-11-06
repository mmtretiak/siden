package logger

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type LogType int

const (
	DevelopmentLogType = iota + 1
	ProductionLogType
)

func New(logType LogType) (Logger, error) {
	var log *zap.Logger
	var err error

	switch logType {
	case DevelopmentLogType:
		log, err = zap.NewDevelopment()
		if err != nil {
			return nil, fmt.Errorf("failed to create logger, reason: %s", err.Error())
		}
	case ProductionLogType:
		log, err = zap.NewProduction()
		if err != nil {
			return nil, fmt.Errorf("failed to create logger, reason: %s", err.Error())
		}
	default:
		return nil, errors.New("unknown logger type")
	}

	logger := &logger{
		log: log.Sugar(),
	}

	return logger, nil
}

type Logger interface {
	Infof(msg string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
}

type logger struct {
	log *zap.SugaredLogger
}

func (l *logger) Infof(msg string, args ...interface{}) {
	l.log.Infof(msg, args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}
