package utils

import (
	"cloud.google.com/go/errorreporting"
	"context"
	"errors"
	"fmt"
	"heypay-cash-in-server/constants"
	"heypay-cash-in-server/settings"
	. "heypay-cash-in-server/utils/types"
	"io/ioutil"
	"runtime"
	"time"

	log "github.com/sirupsen/logrus"
)

// LogSeverity possible values
const (
	DEFAULT   = "default"
	DEBUG     = "debug"
	INFO      = "info"
	NOTICE    = "notice"
	WARNING   = "warning"
	ERROR     = "error"
	CRITICAL  = "critical"
	ALERT     = "alert"
	EMERGENCY = "emergency"
)

var errorClient *errorreporting.Client

func init() {
	if settings.V.Env == "LOCAL" {
		formatter := &log.TextFormatter{
			FullTimestamp: true,
		}
		log.SetFormatter(formatter)
	} else {
		log.SetFormatter(&log.JSONFormatter{})
		// GCP Error Reporting setup
		var err error
		errorClient, err = errorreporting.NewClient(context.Background(), settings.V.Project.ID, errorreporting.Config{
			ServiceName: "cash-in",
			OnError: func(err error) {
				log.Printf("Could not log error: %v", err)
			},
		})
		if err != nil {
			log.Fatal(err)
		}
	}
}

func createFields(baseFields, additionalFields Map) log.Fields {
	fields := log.Fields{
		"version": settings.V.Version,
	}
	for key, entry := range baseFields {
		fields[key] = entry
	}
	for key, entry := range additionalFields {
		fields[key] = entry
	}
	return fields
}

func getStackTrace() (string, int, string) {
	file := ""
	var skip = 2
	for {
		_, file, _, ok := runtime.Caller(skip)
		skip++
		if !ok || !IsUtil(file) || skip > 7 {
			break
		}
	}
	pc, file, line, ok := runtime.Caller(skip - 1)
	if !ok {
		return "undefined.go", -1, "Undefined"
	}
	function := runtime.FuncForPC(pc)
	return TruncateFileName(file), line, TruncateFunctionName(function.Name())
}

func DisableLogging() {
	log.SetOutput(ioutil.Discard)
}

func LogEntryWithId(ctx context.Context, severity string, metadata Map) *log.Entry {
	fileName, line, functionName := getStackTrace()
	ip := GetIpAddress(ctx)
	loggerID := GetLoggerID(ctx)
	if metadata == nil {
		metadata = Map{
			"step": GetCurrentStepLabel(ctx),
		}
	} else {
		metadata["step"] = GetCurrentStepLabel(ctx)
	}
	return log.WithFields(createFields(Map{
		"filename": fileName,
		"function": functionName,
		"ip":       ip,
		"line":     line,
		"loggerId": loggerID,
		"severity": severity,
	}, metadata))
}

func LogLastStep(ctx context.Context, now int64) {
	currentStepLabel := ctx.Value(constants.ContextKeyCurrentStepLabel)
	if currentStepLabel != nil {
		startTime := ctx.Value(constants.ContextKeyStartTime).(int64)
		currentStepTime := ctx.Value(constants.ContextKeyCurrentStepTime).(int64)
		elapsedTimeFromPreviousStep := (now - currentStepTime) / 1000000
		Debugmf(ctx, Map{
			"elapsedTimeFromPreviousStep": elapsedTimeFromPreviousStep,
			"step":                        currentStepLabel,
			"totalElapsedTime":            (now - startTime) / 1000000,
		}, "Step %s took %d ms", currentStepLabel, elapsedTimeFromPreviousStep)
	}
}

func ChangeStep(ctx *context.Context, step, code string) {
	now := time.Now().UnixNano()
	LogLastStep(*ctx, now)
	*ctx = context.WithValue(*ctx, constants.ContextKeyCurrentStepCode, code)
	*ctx = context.WithValue(*ctx, constants.ContextKeyCurrentStepLabel, step)
	*ctx = context.WithValue(*ctx, constants.ContextKeyCurrentStepTime, now)
}

func Debug(ctx context.Context, metadata Map, message string) {
	LogEntryWithId(ctx, DEBUG, metadata).Info(message)
}

func Debugf(ctx context.Context, message string, args ...interface{}) {
	LogEntryWithId(ctx, DEBUG, nil).Infof(message, args...)
}

func Debugmf(ctx context.Context, metadata Map, message string, args ...interface{}) {
	LogEntryWithId(ctx, DEBUG, metadata).Infof(message, args...)
}

func Info(ctx context.Context, message string) {
	LogEntryWithId(ctx, INFO, nil).Info(message)
}

func Infof(ctx context.Context, message string, args ...interface{}) {
	LogEntryWithId(ctx, INFO, nil).Infof(message, args...)
}

func Infomf(ctx context.Context, metadata Map, message string, args ...interface{}) {
	LogEntryWithId(ctx, INFO, metadata).Infof(message, args...)
}

func Warn(ctx context.Context, message string) {
	LogEntryWithId(ctx, WARNING, nil).Warn(message)
}

func Warnf(ctx context.Context, message string, args ...interface{}) {
	LogEntryWithId(ctx, WARNING, nil).Warnf(message, args...)
}

func Warnmf(ctx context.Context, metadata Map, message string, args ...interface{}) {
	LogEntryWithId(ctx, WARNING, metadata).Warnf(message, args...)
}

func Error(ctx context.Context, message string) {
	ReportError(ctx, message)
	LogEntryWithId(ctx, ERROR, nil).Error(message)
}

func Errorf(ctx context.Context, message string, args ...interface{}) {
	ReportError(ctx, message, args...)
	LogEntryWithId(ctx, ERROR, nil).Errorf(message, args...)
}

func Errormf(ctx context.Context, metadata Map, message string, args ...interface{}) {
	ReportError(ctx, message, args...)
	LogEntryWithId(ctx, ERROR, metadata).Errorf(message, args...)
}

func Fatal(ctx context.Context, message string) {
	ReportError(ctx, message)
	LogEntryWithId(ctx, ALERT, nil).Fatal(message)
}

func Fatalf(ctx context.Context, message string, args ...interface{}) {
	ReportError(ctx, message, args...)
	LogEntryWithId(ctx, ALERT, nil).Fatalf(message, args...)
}

func ReportError(ctx context.Context, message string, args ...interface{}) {
	if settings.V.Env == "LOCAL" {
		return
	}
	if ctx != nil {
		errorClient.Report(errorreporting.Entry{
			Error: errors.New(fmt.Sprintf(message, args...)),
			User: GetLoggerID(ctx),
		})
	} else {
		errorClient.Report(errorreporting.Entry{Error: errors.New(fmt.Sprintf(message, args...))})
	}
}
