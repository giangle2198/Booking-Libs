package log

import (
	"encoding/json"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *logger = &logger{}
)

type logger struct {
	*zap.SugaredLogger
	L *zap.Logger
}

func Object(key string, value interface{}) zapcore.Field {
	marshalByte, err := json.Marshal(value)
	if err != nil {
		zap.S().Panic(err)
		return zapcore.Field{Key: key, Type: zapcore.ObjectMarshalerType, Interface: value}
	}

	masked := GetJSONMaskLogging().MaskJSON(string(marshalByte))
	return zapcore.Field{Key: key, Type: zapcore.StringType, Interface: masked, String: masked}
}

func configLogLevel(defaultEnv string) zapcore.Level {
	env := os.Getenv("ENVIROMENT")

	if env == "" {
		env = "P"
	}

	level := zapcore.ErrorLevel
	switch env {
	case "D", "d", "dev", "DEV":
		level = zapcore.DebugLevel
	case "P", "p", "PROD", "prod":
		level = zapcore.WarnLevel
	}
	return level
}

func InitZap(app, env string, maskFields map[string]string) error {
	logLevel := configLogLevel(env)
	encoderConfig := zapcore.EncoderConfig{
		MessageKey: "message",

		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,

		TimeKey:    "time",
		EncodeTime: zapcore.ISO8601TimeEncoder,

		CallerKey:    "caller",
		EncodeCaller: zapcore.ShortCallerEncoder,

		NameKey:    "app",
		EncodeName: zapcore.FullNameEncoder,
	}

	InitJSONMaskLogging(maskFields)
	zap.RegisterEncoder("custom-json", func(cfg zapcore.EncoderConfig) (zapcore.Encoder, error) {
		return NewJSONEncoder(cfg), nil
	})

	cfg := zap.Config{
		Encoding:         "custom-json",
		Level:            zap.NewAtomicLevelAt(logLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    encoderConfig,
	}

	l, err := cfg.Build()
	if err != nil {
		return err
	}

	l = l.Named(app)
	zap.ReplaceGlobals(l)
	Logger = &logger{l.Sugar(), l}
	return nil
}
