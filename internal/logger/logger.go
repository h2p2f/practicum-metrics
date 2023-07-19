// Package logger реализует логирование синглтоном. Используется в сервере и агенте раздельно.
//
// Package logger implements logging singleton. Used in the server and agent separately.
package logger

import (
	"go.uber.org/zap"
)

// Log  - объект логгера
// Log - logger object
var Log *zap.Logger = zap.NewNop()

// InitLogger - функция для инициализации логгера
// InitLogger is function for initializing logger
func InitLogger(level string) error {

	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	//used for development logging
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	Log, err = cfg.Build()
	if err != nil {
		return err
	}
	return nil
}
