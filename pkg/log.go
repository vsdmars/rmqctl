package pkg

import (
	// zap is a reflection-free, zero-allocation JSON encoder.

	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type appLogger struct {
	*zap.Logger
	atom *zap.AtomicLevel
}

func (al *appLogger) setLevel(l zapcore.Level) {
	al.atom.SetLevel(l)
}

var once sync.Once
var logger appLogger

func init() {
	initLogger()
}

func logCleanUp() {
	logger.Sync()
}

func initLogger() {
	initLogger := func() {
		// default log level set to 'info'
		atom := zap.NewAtomicLevelAt(zap.InfoLevel)

		config := zap.Config{
			Level:       atom,
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding:         "json", // console, json, toml
			EncoderConfig:    zap.NewProductionEncoderConfig(),
			OutputPaths:      []string{"stderr"},
			ErrorOutputPaths: []string{"stderr"},
		}

		mylogger, err := config.Build()
		if err != nil {
			fmt.Printf("Initialize zap logger error: %v", err)
			os.Exit(1)
		}

		logger = appLogger{mylogger, &atom}
	}

	once.Do(initLogger)
}
