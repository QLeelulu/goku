package goku

import (
    "fmt"
    "log"
    "os"
)

const (
    LOG_LEVEL_NO = iota
    LOG_LEVEL_ERROR
    LOG_LEVEL_WARN
    LOG_LEVEL_NOTICE
    LOG_LEVEL_LOG
)

var loger *log.Logger = log.New(os.Stdout, "", log.LstdFlags)

type logger interface {
    LogLevel() int
    SetLogLevel(level int)
    Log(args ...interface{})
    Logln(args ...interface{})
    Logf(format string, args ...interface{})
    Notice(args ...interface{})
    Noticeln(args ...interface{})
    Noticef(format string, args ...interface{})
    Warn(args ...interface{})
    Warnln(args ...interface{})
    Warnf(format string, args ...interface{})
    Error(args ...interface{})
    Errorln(args ...interface{})
    Errorf(format string, args ...interface{})
}

type DefaultLogger struct {
    Logger    *log.Logger
    LOG_LEVEL int
}

func (l *DefaultLogger) LogLevel() int {
    return l.LOG_LEVEL
}

func (l *DefaultLogger) SetLogLevel(level int) {
    l.LOG_LEVEL = level
}

func (l *DefaultLogger) Log(args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_LOG {
        l.Logger.Output(3, fmt.Sprint(args...))
    }
}

func (l *DefaultLogger) Logln(args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_LOG {
        l.Logger.Output(3, fmt.Sprintln(args...))
    }
}

func (l *DefaultLogger) Logf(format string, args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_LOG {
        l.Logger.Output(3, fmt.Sprintf(format, args...))
    }
}

func (l *DefaultLogger) Notice(args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_NOTICE {
        v := make([]interface{}, 1, len(args)+1)
        v[0] = "[NOTICE]"
        v = append(v, args...)
        l.Logger.Output(3, fmt.Sprint(v...))
    }
}

func (l *DefaultLogger) Noticeln(args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_NOTICE {
        v := make([]interface{}, 1, len(args)+1)
        v[0] = "[NOTICE]"
        v = append(v, args...)
        l.Logger.Output(3, fmt.Sprintln(v...))
    }
}

func (l *DefaultLogger) Noticef(format string, args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_NOTICE {
        l.Logger.Output(3, fmt.Sprintf("[NOTICE] "+format, args...))
    }
}

func (l *DefaultLogger) Warn(args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_WARN {
        v := make([]interface{}, 1, len(args)+1)
        v[0] = "[WARN]"
        v = append(v, args...)
        l.Logger.Output(3, fmt.Sprint(v...))
    }
}

func (l *DefaultLogger) Warnln(args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_WARN {
        v := make([]interface{}, 1, len(args)+1)
        v[0] = "[WARN]"
        v = append(v, args...)
        l.Logger.Output(3, fmt.Sprintln(v...))
    }
}

func (l *DefaultLogger) Warnf(format string, args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_WARN {
        l.Logger.Output(3, fmt.Sprintf("[WARN] "+format, args...))
    }
}

func (l *DefaultLogger) Error(args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_ERROR {
        v := make([]interface{}, 1, len(args)+1)
        v[0] = "[ERROR]"
        v = append(v, args...)
        l.Logger.Output(3, fmt.Sprint(v...))
    }
}

func (l *DefaultLogger) Errorln(args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_ERROR {
        v := make([]interface{}, 1, len(args)+1)
        v[0] = "[ERROR]"
        v = append(v, args...)
        l.Logger.Output(3, fmt.Sprintln(v...))
    }
}

func (l *DefaultLogger) Errorf(format string, args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_ERROR {
        l.Logger.Output(3, fmt.Sprintf("[ERROR] "+format, args...))
    }
}

var __logger logger = &DefaultLogger{
    Logger:    log.New(os.Stdout, "", log.LstdFlags),
    LOG_LEVEL: LOG_LEVEL_LOG,
}

func Logger() logger {
    return __logger
}

func SetLogger(l logger) {
    __logger = l
}

func SetLogLevel(level int) {
    __logger.SetLogLevel(level)
}
