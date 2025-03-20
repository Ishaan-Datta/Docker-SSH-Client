package utils

// import (
// 	"log"

// 	"your-module-path/messagelogger"
// )

// func main() {
// 	// Initialize the logger with a log file path.
// 	logger, err := messagelogger.NewLogger("messages.log")
// 	if err != nil {
// 		log.Fatalf("Failed to initialize logger: %v", err)
// 	}
// 	defer logger.Close()

// 	// Log some messages.
// 	logger.LogMessage(messagelogger.Message{
// 		Type:    messagelogger.TypeA,
// 		Content: "This is a TypeA message.",
// 	})

// 	logger.LogMessage(messagelogger.Message{
// 		Type:    messagelogger.TypeB,
// 		Content: "This is a TypeB message.",
// 	})

// 	logger.LogMessage(messagelogger.Message{
// 		Type:    "UnknownType",
// 		Content: "This is an unknown type message.",
// 	})
// }
