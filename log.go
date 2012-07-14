package goku

import (
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

type defaultLogger struct {
    logger    *log.Logger
    LOG_LEVEL int
}

func (l *defaultLogger) LogLevel() int {
    return l.LOG_LEVEL
}

func (l *defaultLogger) Log(args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_LOG {
        l.logger.Print(args...)
    }
}

func (l *defaultLogger) Logln(args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_LOG {
        l.logger.Println(args...)
    }
}

func (l *defaultLogger) Logf(format string, args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_LOG {
        l.logger.Printf(format, args...)
    }
}

func (l *defaultLogger) Notice(args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_NOTICE {
        v := make([]interface{}, 1, len(args)+1)
        v[0] = "[NOTICE]"
        v = append(v, args...)
        // l.logger.Print("[NOTICE]")
        l.logger.Print(v...)
    }
}

func (l *defaultLogger) Noticeln(args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_NOTICE {
        v := make([]interface{}, 1, len(args)+1)
        v[0] = "[NOTICE]"
        v = append(v, args...)
        // l.logger.Print("[NOTICE] ")
        l.logger.Println(v...)
    }
}

func (l *defaultLogger) Noticef(format string, args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_NOTICE {
        l.logger.Printf("[NOTICE] "+format, args...)
    }
}

func (l *defaultLogger) Warn(args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_WARN {
        v := make([]interface{}, 1, len(args)+1)
        v[0] = "[WARN]"
        v = append(v, args...)
        // l.logger.Print("[WARN] ")
        l.logger.Print(v...)
    }
}

func (l *defaultLogger) Warnln(args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_WARN {
        v := make([]interface{}, 1, len(args)+1)
        v[0] = "[WARN]"
        v = append(v, args...)
        // l.logger.Print("[WARN] ")
        l.logger.Println(v...)
    }
}

func (l *defaultLogger) Warnf(format string, args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_WARN {
        l.logger.Printf("[WARN] "+format, args...)
    }
}

func (l *defaultLogger) Error(args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_ERROR {
        v := make([]interface{}, 1, len(args)+1)
        v[0] = "[ERROR]"
        v = append(v, args...)
        // l.logger.Print("[ERROR] ")
        l.logger.Print(v...)
    }
}

func (l *defaultLogger) Errorln(args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_ERROR {
        v := make([]interface{}, 1, len(args)+1)
        v[0] = "[ERROR]"
        v = append(v, args...)
        // l.logger.Print("[ERROR] ")
        l.logger.Println(v...)
    }
}

func (l *defaultLogger) Errorf(format string, args ...interface{}) {
    if l.LOG_LEVEL >= LOG_LEVEL_ERROR {
        l.logger.Printf("[ERROR] "+format, args...)
    }
}

var dlogger logger = &defaultLogger{
    logger:    log.New(os.Stdout, "", log.LstdFlags),
    LOG_LEVEL: LOG_LEVEL_LOG,
}

func Logger() logger {
    return dlogger
}
