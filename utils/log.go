package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

// LogEventType represents the type of event to log
type LogEventType string

const (
    INFO    LogEventType = "INFO"
    WARNING LogEventType = "WARNING"
    ERROR   LogEventType = "ERROR"
)

// LogMessage represents a log message
type LogMessage struct {
    EventType   LogEventType
    Text        string
    ContainerID string
}

// Logger struct to hold the log file and channel
type Logger struct {
    file    *os.File
    logChan chan LogMessage
    done    chan struct{}
}

// NewLogger creates a new logger instance
func NewLogger(filePath string) (*Logger, error) {
    file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, err
    }
    logger := &Logger{
        file:    file,
        logChan: make(chan LogMessage, 100),
        done:    make(chan struct{}),
    }
    go logger.listen()
    return logger, nil
}

// listen listens for log messages and writes them to the file
func (l *Logger) listen() {
    for {
        select {
        case msg := <-l.logChan:
            timestamp := time.Now().Format("2006/01/02 15:04:05")
            logMessage := fmt.Sprintf("%s [%s] [%s] %s\n", timestamp, msg.EventType, msg.ContainerID, msg.Text)
            if _, err := l.file.WriteString(logMessage); err != nil {
                log.Printf("Failed to write to log file: %v", err)
            }
        case <-l.done:
            return
        }
    }
}

// Log logs a message with the given event type, container ID, and text
func (l *Logger) Log(eventType LogEventType, containerID, text string) {
    l.logChan <- LogMessage{EventType: eventType, ContainerID: containerID, Text: text}
}

// Close closes the log file and stops the listener
func (l *Logger) Close() {
    close(l.done)
    close(l.logChan)
    if err := l.file.Close(); err != nil {
        log.Printf("Failed to close log file: %v", err)
    }
}

func main() {
    logger, err := NewLogger("application.log")
    if err != nil {
        log.Fatalf("Failed to create logger: %v", err)
    }
    defer logger.Close()

    logger.Log(INFO, "container123", "Application started")
    logger.Log(WARNING, "container123", "This is a warning message")
    logger.Log(ERROR, "container123", "This is an error message")
}