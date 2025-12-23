package transports

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/natefinch/lumberjack"
)

type FileWriter struct {
	logger *lumberjack.Logger
	dir    string
}

func (fw *FileWriter) Write(p []byte) (n int, err error) {
	return fw.logger.Write(p)
}

func (fw *FileWriter) Close() error {
	return fw.logger.Close()
}

func NewFileWriter(dir string) *FileWriter {
	filePath := filepath.Join(dir, fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02")))
	logger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    10,   // Max size in megabytes before log is rotated
		MaxBackups: 3,    // Max number of old log files to retain
		MaxAge:     28,   // Max number of days to retain old log files
		Compress:   true, // Compress old log files
	}
	return &FileWriter{
		logger: logger,
		dir:    dir,
	}
}
