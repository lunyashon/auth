package loger

import (
	"io"
	"log/slog"
	"os"

	"gopkg.in/natefinch/lumberjack.v2"
)

func ExecLog(env, path string) *slog.Logger {

	var (
		level  slog.Level
		writer io.Writer
	)

	switch env {
	case "local":
		level = slog.LevelDebug
		writer = os.Stdout
	case "prod":
		level = slog.LevelDebug
		writer = createFileWriter(path)
	case "dev":
		level = slog.LevelInfo
		writer = io.MultiWriter(
			os.Stdout,
			createFileWriter(path),
		)
	default:
		level = slog.LevelInfo
		writer = os.Stdout
	}

	return slog.New(
		slog.NewJSONHandler(writer, &slog.HandlerOptions{Level: level}),
	)
}

func createFileWriter(path string) io.Writer {
	if err := os.MkdirAll(path, 0755); err != nil {
		panic(err)
	}

	return &lumberjack.Logger{
		Filename:   path + "/app.log",
		MaxSize:    100,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
		LocalTime:  true,
	}
}
