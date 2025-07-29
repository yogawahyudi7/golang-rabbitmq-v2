package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"

	"golang-rabbitmq-v2/internal/config"
)

type Logger struct {
	*logrus.Logger
	fileLogger *lumberjack.Logger
}

func New(cfg *config.LoggerConfig) *Logger {
	log := logrus.New()

	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	if cfg.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	logger := &Logger{
		Logger: log,
	}

	if cfg.FileEnabled {
		if err := logger.setupFileLogging(cfg); err != nil {
			log.WithError(err).Error("Failed to setup file logging, using console only")
			log.SetOutput(os.Stdout)
		}
	} else {
		log.SetOutput(os.Stdout)
	}

	return logger
}

func (l *Logger) setupFileLogging(cfg *config.LoggerConfig) error {
	if err := os.MkdirAll(cfg.LogDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	filename := l.generateDailyLogFilename(cfg.LogDir)

	l.fileLogger = &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
		LocalTime:  true,
	}

	multiWriter := io.MultiWriter(os.Stdout, l.fileLogger)
	l.Logger.SetOutput(multiWriter)

	l.Logger.WithFields(logrus.Fields{
		"log_file":     filename,
		"max_size_mb":  cfg.MaxSize,
		"max_backups":  cfg.MaxBackups,
		"max_age_days": cfg.MaxAge,
		"compress":     cfg.Compress,
	}).Info("File logging enabled")

	go l.watchDateChange(cfg)

	return nil
}

func (l *Logger) generateDailyLogFilename(logDir string) string {
	today := time.Now().Format("2006-01-02")
	return filepath.Join(logDir, fmt.Sprintf("app-%s.log", today))
}

func (l *Logger) watchDateChange(cfg *config.LoggerConfig) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	currentDate := time.Now().Format("2006-01-02")

	for range ticker.C {
		newDate := time.Now().Format("2006-01-02")
		if newDate != currentDate {
			l.Logger.Info("Date changed, switching to new log file")
			
			if err := l.switchToNewLogFile(cfg, newDate); err != nil {
				l.Logger.WithError(err).Error("Failed to switch to new log file")
				continue
			}
			
			currentDate = newDate
			l.Logger.WithField("new_log_file", l.fileLogger.Filename).Info("Switched to new daily log file")
		}
	}
}

func (l *Logger) switchToNewLogFile(cfg *config.LoggerConfig, newDate string) error {
	if l.fileLogger != nil {
		l.fileLogger.Close()
	}

	newFilename := filepath.Join(cfg.LogDir, fmt.Sprintf("app-%s.log", newDate))
	
	l.fileLogger = &lumberjack.Logger{
		Filename:   newFilename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
		LocalTime:  true,
	}

	multiWriter := io.MultiWriter(os.Stdout, l.fileLogger)
	l.Logger.SetOutput(multiWriter)

	return nil
}

func (l *Logger) WithFields(fields map[string]interface{}) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

func (l *Logger) WithContext(component string) *logrus.Entry {
	return l.Logger.WithField("component", component)
}

func (l *Logger) Close() error {
	if l.fileLogger != nil {
		return l.fileLogger.Close()
	}
	return nil
}