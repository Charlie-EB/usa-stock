package sentry

import (
	"log"
	"m/utils"
	"time"

	"github.com/getsentry/sentry-go"
)

func Setup() {

  env, err := utils.GetEnv()
	if err != nil {
		log.Fatalf("failed to get env: %v", err)
	}

  err = sentry.Init(sentry.ClientOptions{Dsn: env["SENTRY_DSN"]})
  if err != nil {
    log.Fatalf("sentry.Init: %s", err)
  }
  // Flush buffered events before the program terminates.
  defer sentry.Flush(2 * time.Second)

  sentry.CaptureMessage("It works!")
}
