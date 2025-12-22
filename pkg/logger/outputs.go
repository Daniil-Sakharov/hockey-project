package logger

import (
	"io"
	"os"

	"go.uber.org/zap/zapcore"
)

// MultiOutput создает multi-output writer
func MultiOutput(outputs ...io.Writer) zapcore.WriteSyncer {
	if len(outputs) == 0 {
		return zapcore.AddSync(os.Stdout)
	}

	if len(outputs) == 1 {
		return zapcore.AddSync(outputs[0])
	}

	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(io.MultiWriter(outputs...)),
	)
}

// CreateFileOutput создает file writer (если нужен)
func CreateFileOutput(filename string) (io.Writer, error) {
	if filename == "" {
		return nil, nil
	}

	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o600) //nolint:gosec // filename is trusted input
	if err != nil {
		return nil, err
	}

	return file, nil
}
