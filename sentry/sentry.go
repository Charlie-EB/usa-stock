package sentry

import (
	"fmt"
	"log"
	"m/utils"
	"time"

	"github.com/getsentry/sentry-go"
)

func Setup() {

  env, err := utils.GetEnv()
	if err != nil {
		log.Fatalf("failed to get env: %w", err)
	}

  err = sentry.Init(sentry.ClientOptions{
    Dsn: env["SENTRY_DSN"],
    Debug: true,
  })
  if err != nil {
    log.Fatalf("sentry.Init: %s", err)
  }
  // Flush buffered events before the program terminates.
  defer sentry.Flush(2 * time.Second)

  sentry.CaptureMessage("Server started and connected to Sentry a-okay")
}


func Notify (err error, context string) {
  if err == nil {
    return 
  }

  // Add context info to the message if provided
	message := fmt.Sprintf("Error: %v", err)
	if context != "" {
		message = fmt.Sprintf("%s | Context: %s", message, context)
	}

	sentry.CaptureMessage(message)
	sentry.Flush(2 * time.Second)
}