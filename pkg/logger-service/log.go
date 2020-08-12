package logger_service

import "github.com/sirupsen/logrus"

var self ILogger = nil

type Log struct {
	lg *logrus.Logger
}

func (l Log) Debug(message string, fields map[string]interface{}) {
	l.lg.WithFields(fields).Debug(message)
}

func (l Log) Info(message string, fields map[string]interface{}) {
	l.lg.WithFields(fields).Info(message)
}

func (l Log) Warning(message string, fields map[string]interface{}) {
	l.lg.WithFields(fields).Warning(message)
}

func (l Log) Error(message string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{}, 1)
	}

	fields["original_error"] = err

	l.lg.WithFields(fields).Error(message)
}

func (l Log) Fatal(message string, err error, fields map[string]interface{}) {
	if fields == nil {
		fields = make(map[string]interface{}, 1)
	}

	fields["original_error"] = err

	l.lg.WithFields(fields).Fatal(message)
}

func MustNewLogger(conf Config) ILogger {
	if self != nil {
		logrus.Fatal("logger already created")
	}

	lg := logrus.New()
	lg.Formatter = &logrus.JSONFormatter{
		DisableTimestamp: false,
	}

	if conf.IsVerbose {
		lg.Level = logrus.DebugLevel
	}

	self = &Log{
		lg: lg,
	}

	return self
}

func GetLogger() ILogger {
	if self == nil {
		self = &Log{
			lg: logrus.New(),
		}
	}

	return self
}
