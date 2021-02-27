package sink

import (
	"log"

	"github.com/panjf2000/ants/v2"
)

type standardLogger struct {
}

func (s standardLogger) Printf(format string, args ...interface{}) {
	log.Printf(format, args...)
}

// Logger is the logging interface.
type Logger interface {
	ants.Logger
}
