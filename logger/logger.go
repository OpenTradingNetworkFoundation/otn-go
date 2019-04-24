package slog

import (
	"fmt"
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ComponentConfig specifies separate logger configuration
type ComponentConfig struct {
	Preset          string `json:"preset"`
	Level           string `json:"level"`
	JSON            bool   `json:"json"`
	TimestampFormat string `json:"timestamp_format"`
}

// Config provides structured logger configuration
type Config struct {
	ComponentConfigs map[string]*ComponentConfig `json:"components"`
}

// Component represents single component logging interface
type Component struct {
	logger      *zap.Logger
	sugarLogger *zap.SugaredLogger
}

// L returns strict typed logger which is very fast
func (c *Component) L() *zap.Logger {
	return c.logger
}

// S returns sugared logger which provides less verbose API
func (c *Component) S() *zap.SugaredLogger {
	return c.sugarLogger
}

// SLog represents structured logging interface
type SLog struct {
	components map[string]*Component
}

// C returns logger for specified component
func (s *SLog) C(component string) *Component {
	return s.components[component]
}

const (
	tsFormatUnix    = "unix"
	tsFormatISO8601 = "iso8601"
)

var presets = map[string]*ComponentConfig{
	"developer": &ComponentConfig{
		Level:           "debug",
		JSON:            false,
		TimestampFormat: tsFormatISO8601,
	},
	"production": &ComponentConfig{
		Level:           "info",
		JSON:            true,
		TimestampFormat: tsFormatUnix,
	},
}

func newComponent(c *ComponentConfig) *Component {
	cfg := c
	if cfg.Preset != "" {
		var ok bool
		cfg, ok = presets[cfg.Preset]
		if !ok {
			panic(fmt.Sprintf("Uknown logger preset %s", cfg.Preset))
		}
	}
	var level zap.AtomicLevel
	if err := level.UnmarshalText([]byte(c.Level)); err != nil {
		panic(fmt.Sprintf("Bad logger level %s", c.Level))
	}

	encoderProducer := zapcore.NewConsoleEncoder
	if cfg.JSON {
		encoderProducer = zapcore.NewJSONEncoder
	}
	encoderCfg := zap.NewDevelopmentEncoderConfig()
	switch cfg.TimestampFormat {
	case tsFormatISO8601:
		encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	default:
		encoderCfg.EncodeTime = zapcore.EpochMillisTimeEncoder
	}
	encoderCfg.CallerKey = "caller"
	encoderCfg.EncodeCaller = zapcore.ShortCallerEncoder
	//encoderCfg.InitialFields = map[string]interface{}{"test", os.GetPID() }
	core := zapcore.NewCore(
		encoderProducer(encoderCfg),
		zapcore.Lock(os.Stdout),
		level,
	)
	logger := zap.New(core)
	return &Component{
		logger:      logger,
		sugarLogger: logger.Sugar(),
	}
}

// NewSLog produces new instance of SLog based on specified configuration
func NewSLog(cfg *Config) *SLog {

	log.Printf("%v", presets)
	sl := &SLog{
		components: make(map[string]*Component),
	}
	for key, val := range cfg.ComponentConfigs {
		sl.components[key] = newComponent(val)
	}
	return sl
}
