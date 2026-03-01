package utils

import (
	"log"
	"os"
	"strings"
)

var isDebug bool

func InitLogger() {
	env := strings.ToLower(os.Getenv("DEBUG"))
	isDebug = (env == "true" || env == "1")
}

func Info(format string, v ...any) {
	log.Printf("ℹ️ [INFO] "+format, v...)
}

func Debug(format string, v ...any) {
	if !isDebug {
		return
	}
	log.Printf("🛠️ [DEBUG] "+format, v...)
}

func Error(format string, v ...any) {
	log.Printf("❌ [ERROR] "+format, v...)
}