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
	debugSeverity    severity = "DEBUG"
	errorSeverity    severity = "ERROR"
	infoSeverity     severity = "INFO"
	warningSeverity  severity = "WARNING"
)

var w *syslog.Writer
var localLogging bool

func Init(url, system string) {
	if url != "" {
		var err error
		w, err = syslog.Dial("udp", url, syslog.LOG_SYSLOG, system)
		if err != nil {
			log.Fatalf("failed to dial syslog, not able to contact %s as %s, error was %s", url, system, err)
			localLogging = true
		}
		localLogging = false
	} else {
		log.Println("Will log local only")
		localLogging = true
	}
}

func Close() {
	w.Close()
}

func Info(tags []string, msg string, err string) {
	print(infoSeverity, tags, msg, err)
}

func Warning(tags []string, msg string, err string) {
	print(warningSeverity, tags, msg, err)
}
func Error(tags []string, msg string, err string) {
	print(errorSeverity, tags, msg, err)
}
func Debug(tags []string, msg string, err string) {
	print(debugSeverity, tags, msg, err)
}

func Critical(tags []string, msg string, err string) {
	print(criticalSeverity, tags, msg, err)
	os.Exit(1)
}

func print(infoType severity, tags []string, msg string, err string) {
	flattenTags := fmt.Sprintf("[%s]", flatten(tags))
	finalMsg := fmt.Sprintf("%-8s -  %-30s - %-70s ", infoType, flattenTags, msg)
	if err != "" {
		finalMsg = fmt.Sprintf("%s - %s", finalMsg, err)
	}
	log.Println(finalMsg)
	if !localLogging {
		switch infoType {
		case infoSeverity:
			w.Info(finalMsg)
		case warningSeverity:
			w.Warning(finalMsg)
		case errorSeverity:
			w.Err(finalMsg)
		case debugSeverity:
			w.Debug(finalMsg)
		case criticalSeverity:
			w.Crit(finalMsg)
		default:
			w.Notice(finalMsg)
		}
	}
}
