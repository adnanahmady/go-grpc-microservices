package applog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adnanahmady/go-grpc-microservices/config"
	"github.com/adnanahmady/go-grpc-microservices/pkg/app"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Fatal(msg string, args ...any)
	Panic(msg string, args ...any)
	NewWith(args ...any) Logger
}

var _ Logger = (*AppLogger)(nil)

type AppLogger struct {
	lgr zerolog.Logger
}
func NewAppLogger(cfg *config.Config, serviceName string) *AppLogger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	writers := getWriters(cfg)
	multiWriter := zerolog.MultiLevelWriter(writers...)
	level := getLevel(cfg)
	zerolog.SetGlobalLevel(level)

	lgr := zerolog.New(multiWriter).With().
		Str("service_name", serviceName).Timestamp().Logger()

	return &AppLogger{lgr: lgr}
}

func getWriters(cfg *config.Config) []io.Writer {
	writers := make([]io.Writer, 0, 2)

	writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout})
	if cfg.Log.WriteToFile {
		writers = append(writers, getFileWriter(cfg))
	}

	return writers
}

func getFileWriter(cfg *config.Config) io.Writer {
	return &lumberjack.Logger{
		Filename: filepath.Join(
			filepath.Join(app.GetRootDir(), cfg.Log.Dir),
			fmt.Sprintf("app-%s.log", time.Now().Format("2006-01-02")),
		),
		MaxAge:   cfg.Log.MaxAge,
		MaxSize:  cfg.Log.MaxSize,
		Compress: cfg.Log.CompressLogFiles,
	}
}

func getLevel(cfg *config.Config) zerolog.Level {
	switch cfg.Log.Level {
	case "debug":
		return zerolog.DebugLevel
	case "info", "information":
		return zerolog.InfoLevel
	case "warn", "warning":
		return zerolog.WarnLevel
	case "error", "err":
		return zerolog.ErrorLevel
	case "fatal", "critical", "crit":
		return zerolog.FatalLevel
	case "panic", "emergency":
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel
	}
}

func (l *AppLogger) NewWith(args ...any) Logger {
	return &AppLogger{lgr: l.lgr.With().Fields(args).Logger()}
}

func (l *AppLogger) Debug(msg string, args ...any) {
	write(l.lgr.Debug(), msg, args)
}

func (l *AppLogger) Info(msg string, args ...any) {
	write(l.lgr.Info(), msg, args)
}

func (l *AppLogger) Warn(msg string, args ...any) {
	write(l.lgr.Warn(), msg, args)
}

func (l *AppLogger) Error(msg string, args ...any) {
	write(l.lgr.Error(), msg, args)
}

func (l *AppLogger) Fatal(msg string, args ...any) {
	write(l.lgr.Fatal(), msg, args)
}

func (l *AppLogger) Panic(msg string, args ...any) {
	write(l.lgr.Panic(), msg, args)
}

func write(e *zerolog.Event, msg string, fields []any) {
	if len(fields) > 0 {
		if err, ok := fields[0].(error); ok {
			e = e.Err(err)
			fields = fields[1:]
		}
	}

	args := make([]any, 0, len(fields))
	count := strings.Count(msg, "%")
	if count > 0 {
		args = fields[:count]
		fields = fields[count:]
	}

	e.Fields(fields).Msgf(msg, args...)
}
