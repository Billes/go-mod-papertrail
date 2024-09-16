package papertrail

import (
	"fmt"
	"log"
	"log/syslog"
	"os"
)

func flatten(tags []string) string {
	result := ""
	for _, tag := range tags {
		result += fmt.Sprintf("%s,", tag)
	}
	return result[:len(result)-1]
}

type severity string

const (
	criticalSeverity severity = "CRITICAL"
	errorSeverity    severity = "ERROR"
	warningSeverity  severity = "WARNING"
	infoSeverity     severity = "INFO"
	debugSeverity    severity = "DEBUG"
)

var w *syslog.Writer
var remoteLogging bool
var errorLevel severity

func Init(url, system string) {
	remoteLogging = true
	if url != "" {
		var err error
		w, err = syslog.Dial("udp", url, syslog.LOG_SYSLOG, system)
		if err != nil {
			log.Fatalf("failed to dial syslog, not able to contact %s as %s, error was %s", url, system, err)
			remoteLogging = false
		}
		log.Printf("Will log remotely")
	} else {
		log.Println("Will log local only")
		remoteLogging = false
	}
	errorL := os.Getenv("LOG_LEVEL")
	switch errorL {
	case string(criticalSeverity):
		errorLevel = criticalSeverity
	case string(errorSeverity):
		errorLevel = errorSeverity
	case string(warningSeverity):
		errorLevel = warningSeverity
	case string(infoSeverity):
		errorLevel = infoSeverity
	case string(debugSeverity):
		errorLevel = debugSeverity
	default:
		errorLevel = infoSeverity
	}
	log.Printf("Will log from %s", errorLevel)
}

func Close() {
	if w != nil {
		w.Close()
	}
}

func Debug(tags []string, msg string, err string) {
	if errorLevel == debugSeverity {
		print(debugSeverity, tags, msg, err)
	}
}

func Info(tags []string, msg string, err string) {
	if errorLevel == debugSeverity || errorLevel == infoSeverity {
		print(infoSeverity, tags, msg, err)
	}
}

func Warning(tags []string, msg string, err string) {
	if errorLevel == debugSeverity || errorLevel == infoSeverity || errorLevel == warningSeverity {
		print(warningSeverity, tags, msg, err)
	}
}

func Error(tags []string, msg string, err string) {
	if errorLevel == debugSeverity || errorLevel == infoSeverity || errorLevel == warningSeverity || errorLevel == errorSeverity {
		print(errorSeverity, tags, msg, err)
	}
}

func Critical(tags []string, msg string, err string) {
	if errorLevel == debugSeverity || errorLevel == infoSeverity || errorLevel == warningSeverity || errorLevel == errorSeverity || errorLevel == criticalSeverity {
		print(criticalSeverity, tags, msg, err)
	}
}

func print(infoType severity, tags []string, msg string, err string) {
	flattenTags := fmt.Sprintf("[%s]", flatten(tags))
	finalMsg := fmt.Sprintf("%-8s - %-30s - %-70s ", infoType, flattenTags, msg)
	if err != "" {
		finalMsg = fmt.Sprintf("%s - %s", finalMsg, err)
	}
	log.Println(finalMsg)
	if remoteLogging {
		switch infoType {
		case infoSeverity:
			_ = w.Info(finalMsg)
		case warningSeverity:
			_ = w.Warning(finalMsg)
		case errorSeverity:
			_ = w.Err(finalMsg)
		case debugSeverity:
			_ = w.Debug(finalMsg)
		case criticalSeverity:
			_ = w.Crit(finalMsg)
		default:
			_ = w.Info(finalMsg)
		}
	}
}
