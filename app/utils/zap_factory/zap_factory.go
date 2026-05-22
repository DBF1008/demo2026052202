package zap_factory

import (
	"ginskeleton/app/global/variable"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"time"
)

func CreateZapFactory(entry func(zapcore.Entry) error) *zap.Logger {

	appDebug := variable.ConfigYml.GetBool("AppDebug")

	if appDebug == true {
		if logger, err := zap.NewDevelopment(zap.Hooks(entry)); err == nil {
			return logger
		} else {
			log.Fatal("创建zap日志包失败，详情：" + err.Error())
		}
	}

	encoderConfig := zap.NewProductionEncoderConfig()

	timePrecision := variable.ConfigYml.GetString("Logs.TimePrecision")
	var recordTimeFormat string
	switch timePrecision {
	case "second":
		recordTimeFormat = "2006-01-02 15:04:05"
	case "millisecond":
		recordTimeFormat = "2006-01-02 15:04:05.000"
	default:
		recordTimeFormat = "2006-01-02 15:04:05"

	}
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(recordTimeFormat))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.TimeKey = "created_at"

	var encoder zapcore.Encoder
	switch variable.ConfigYml.GetString("Logs.TextFormat") {
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	default:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	fileName := variable.BasePath + variable.ConfigYml.GetString("Logs.ginskeletonLogName")
	lumberJackLogger := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    variable.ConfigYml.GetInt("Logs.MaxSize"),
		MaxBackups: variable.ConfigYml.GetInt("Logs.MaxBackups"),
		MaxAge:     variable.ConfigYml.GetInt("Logs.MaxAge"),
		Compress:   variable.ConfigYml.GetBool("Logs.Compress"),
	}
	writer := zapcore.AddSync(lumberJackLogger)

	zapCore := zapcore.NewCore(encoder, writer, zap.InfoLevel)
	return zap.New(zapCore, zap.AddCaller(), zap.Hooks(entry), zap.AddStacktrace(zap.WarnLevel))
}
