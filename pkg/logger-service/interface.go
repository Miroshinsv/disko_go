package logger_service

type ILogger interface {
	Debug(message string, fields map[string]interface{})
	Info(message string, fields map[string]interface{})
	Warning(message string, fields map[string]interface{})
	Error(message string, err error, fields map[string]interface{})
	Fatal(message string, err error, fields map[string]interface{})
}
