package utils

import (
	"../models"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/rifflock/lfshook"
	"os"
	"strings"
)

type (
	Logger struct {
		Log *logrus.Logger
	}
)

var (
	conf = models.Config{}.Init()
)

func (l *Logger) Init() {
	if l.Log == nil {
		l.Log = logrus.New()

		// Log as JSON instead of the default ASCII formatter.
		l.Log.Formatter = &logrus.JSONFormatter{}

		// Output to stderr instead of stdout, could also be a file.
		// l.Log.Out = os.Stderr

		// Only log the warning severity or above.
		switch strings.ToLower(conf.Logger.Level) {
		case "debug":
			l.Log.Level = logrus.DebugLevel
		case "info":
			l.Log.Level = logrus.InfoLevel
		case "warn":
			l.Log.Level = logrus.WarnLevel
		default:
			l.Log.Level = logrus.ErrorLevel
		}

		// Create log directory if not exist
		logDir := strings.Replace(conf.Logger.Error,
			"/"+strings.Split(conf.Logger.Error, "/")[len(strings.Split(conf.Logger.Error, "/"))-1:][0],
			"", -1)
		if _, err := os.Stat(logDir); os.IsNotExist(err) {
			err := os.MkdirAll(logDir, 0711)
			if err != nil {
				fmt.Println("Failed to create log folder " + logDir)
			}
		}

		// Hook with log files
		l.Log.Hooks.Add(lfshook.NewHook(lfshook.PathMap{
			logrus.InfoLevel:  conf.Logger.Info,
			logrus.ErrorLevel: conf.Logger.Error,
		}))
	}
}

func (l *Logger) Debug(fields map[string]interface{}, errString string) {
	l.Init()
	l.Log.WithFields(fields).Debug(errString)
}

func (l *Logger) Info(fields map[string]interface{}, errString string) {
	l.Init()
	l.Log.WithFields(fields).Info(errString)
}

func (l *Logger) Warn(fields map[string]interface{}, errString string) {
	l.Init()
	l.Log.WithFields(fields).Warn(errString)
}

func (l *Logger) Error(fields map[string]interface{}, errString string) {
	l.Init()
	l.Log.WithFields(fields).Error(errString)
}

func (l *Logger) Fatal(fields map[string]interface{}, errString string) {
	l.Init()
	l.Log.WithFields(fields).Fatal(errString)
}

func (l *Logger) Panic(fields map[string]interface{}, errString string) {
	l.Init()
	l.Log.WithFields(fields).Panic(errString)
}

/* example
func main() {
	l := new(Logger)
	l.Init()

	l.Info(map[string]interface{}{
		"animal": "walrus",
		"size":   10,
	}, "test")

	log.WithFields(log.Fields{
		"omg":    true,
		"number": 122,
	}).Warn("The group's number increased tremendously!")

	log.WithFields(log.Fields{
		"omg":    true,
		"number": 100,
	}).Fatal("The ice breaks!")

	// A common pattern is to re-use fields between logging statements by re-using
	// the logrus.Entry returned from WithFields()
	contextLogger := log.WithFields(log.Fields{
		"common": "this is a common field",
		"other":  "I also should be logged always",
	})

	contextLogger.Info("I'll be logged with common and other field")
	contextLogger.Info("Me too")
}
*/
