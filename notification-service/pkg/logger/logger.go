package logger

import (
    "fmt"
    "log"
    "os"
    "time"
)

// Logger предоставляет функциональность логирования
type Logger struct {
    infoLogger  *log.Logger
    warnLogger  *log.Logger
    errorLogger *log.Logger
    fatalLogger *log.Logger
}

// NewLogger создает новый логгер
func NewLogger() *Logger {
    return &Logger{
        infoLogger:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
        warnLogger:  log.New(os.Stdout, "[WARN] ", log.LstdFlags),
        errorLogger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
        fatalLogger: log.New(os.Stderr, "[FATAL] ", log.LstdFlags),
    }
}

// Info логирует информационное сообщение
func (l *Logger) Info(format string, v ...interface{}) {
    l.infoLogger.Printf(format, v...)
}

// Warn логирует предупреждение
func (l *Logger) Warn(format string, v ...interface{}) {
    l.warnLogger.Printf(format, v...)
}

// Error логирует ошибку
func (l *Logger) Error(format string, v ...interface{}) {
    l.errorLogger.Printf(format, v...)
}

// Fatal логирует фатальную ошибку и завершает программу
func (l *Logger) Fatal(format string, v ...interface{}) {
    l.fatalLogger.Printf(format, v...)
    os.Exit(1)
}