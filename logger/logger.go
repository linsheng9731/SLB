package logger

import (
	"flag"
	"github.com/linsheng9731/slb/config"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
)

var Access = NewLogger("./access.log", config.GlobalConfig)
var Server = NewLogger("./server.log", config.GlobalConfig)

type Logger struct {
	logger *log.Logger
	cfg    *config.Configuration
}

// If this is running with test option set out of logger
// to stdout to avoid generating a new log file
func NewLogger(filename string, cfg *config.Configuration) *Logger {
	var file *os.File
	if v := flag.Lookup("test.v"); v == nil {
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			file.Close()
			log.Fatal("Open log file error!")
		}
		file = f
	} else {
		file = os.Stdout
	}
	l := log.New(file, "[info]", log.Lshortfile)
	l.SetOutput(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    cfg.LogSize, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   //days
		Compress:   true, // disabled by default
	})

	return &Logger{l, cfg}
}

func (l *Logger) Debug(msg ...interface{}) {
	if config.GlobalConfig.LogLevel == "debug" {
		l.logger.SetPrefix("[debug]")
		l.logger.Println(msg)
	}
}

func (l *Logger) Info(msg ...interface{}) {
	Level := config.GlobalConfig.LogLevel
	if Level == "debug" || Level == "info" {
		l.logger.SetPrefix("[info]")
		l.logger.Println(msg)
	}
}

func (l *Logger) Warn(msg ...interface{}) {
	Level := config.GlobalConfig.LogLevel
	if Level == "debug" || Level == "info" || Level == "warn" {
		l.logger.SetPrefix("[warn]")
		l.logger.Println(msg)
	}
}

func (l *Logger) Error(msg ...interface{}) {
	Level := config.GlobalConfig.LogLevel
	if Level != "fatal" {
		l.logger.SetPrefix("[error]")
		l.logger.Println(msg)
	}
}

func (l *Logger) Fatal(msg ...interface{}) {
	Level := config.GlobalConfig.LogLevel
	if Level == "fatal" {
		l.logger.SetPrefix("[fatal]")
		l.logger.Fatalln(msg)
	}
}
