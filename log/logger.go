// +build  !debug

/*
 * Copyright (c) 2018 QLC Chain Team
 *
 * This software is released under the MIT License.
 * https://opensource.org/licenses/MIT
 */

package log

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"github.com/qlcchain/go-qlc/common/util"
	"github.com/qlcchain/go-qlc/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	logfile = "qlc.log"
)

var (
	once      sync.Once
	lumlog    lumberjack.Logger
	logger, _ = zap.NewDevelopment()
)

func InitLog(config *config.Config) error {
	var initErr error
	once.Do(func() {
		logFolder := config.LogDir()
		err := util.CreateDirIfNotExist(logFolder)
		if err != nil {
			initErr = err
		}
		logfile, _ := filepath.Abs(filepath.Join(logFolder, logfile))
		lumlog = lumberjack.Logger{
			Filename:   logfile,
			MaxSize:    10, // megabytesÒ
			MaxBackups: 10,
			MaxAge:     28, // days
			Compress:   true,
			LocalTime:  true,
		}
		var logCfg zap.Config
		err = json.Unmarshal([]byte(`{
		"level": "error",
		"outputPaths": ["stdout"],
		"errorOutputPaths": ["stderr"],
		"encoding": "json",
		"encoderConfig": {
			"messageKey": "message",
			"levelKey": "level",
			"levelEncoder": "lowercase"
		}
	}`), &logCfg)
		if err != nil {
			initErr = err
			fmt.Println(err)
		}
		err = logCfg.Level.UnmarshalText([]byte(config.LogLevel))
		if err != nil {
			initErr = err
			fmt.Println(err)
		}
		logCfg.EncoderConfig = zap.NewProductionEncoderConfig()
		logger, _ = logCfg.Build(zap.Hooks(lumberjackZapHook))
	})

	return initErr
}

//NewLogger create logger by name
func NewLogger(name string) *zap.SugaredLogger {
	return logger.Sugar().Named(name)
}

func lumberjackZapHook(e zapcore.Entry) error {
	_, err := lumlog.Write([]byte(fmt.Sprintf("%s %s [%s] %s %s\n", e.Time.Format(time.RFC3339Nano), e.Level.CapitalString(), e.LoggerName, e.Caller.TrimmedPath(), e.Message)))
	return err
}
