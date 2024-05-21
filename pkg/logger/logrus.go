package logger

import (
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

var once sync.Once

func NewLogger() {
	once.Do(func() {
		log := logrus.New()
		log.Out = os.Stdout

		log.SetFormatter(&logrus.TextFormatter{TimestampFormat: time.RFC3339Nano, FullTimestamp: false})
		log.SetLevel(logrus.InfoLevel)
	})
}
