package logging

import (
	"bytes"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/fatih/color"
)

type Level int

var CurrentLoggingLevel Level
var LoggingOutputReciever chan string
var ColorLogLevelLabelOnly = false
var OutputLogLevelFlag = true
var OutputPath bool = true
var OutputDateTime bool = true
var OutputArrowSuffix bool = true

const (
	BlankLevel Level = 10
	InfoLevel  Level = 3
	WarnLevel  Level = 2
	DebugLevel Level = 1
)

func init() {
	CurrentLoggingLevel = InfoLevel
}

func FlushLogs(logFilePath string, flushInitialised *chan bool) {
	f, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		Error(err.Error())
		*flushInitialised <- true
		return
	}

	*flushInitialised <- true

	LoggingOutputReciever = make(chan string)

	for d := range LoggingOutputReciever {
		_, err = fmt.Fprint(f, d)
		if err != nil {
			Error(err.Error())
			err = f.Close()
			if err != nil {
				Error(err.Error())
			}
			return
		}
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func createCallbackLabel(skip int) string {
	function, _, line, _ := runtime.Caller(skip)
	return fmt.Sprintf("(%s):%d", runtime.FuncForPC(function).Name(), line)
}

//SetLevel allows settings of the level of logging
func SetLevel(loggingLevel Level) {
	CurrentLoggingLevel = loggingLevel
}

//ColoredOutput helper to make it easy to logout with date time stamp
func ColoredOutput(colorPrinter *color.Color, stringToPrint string) {
	if LoggingOutputReciever != nil {
		LoggingOutputReciever <- stringToPrint
	}
	colorPrinter.Printf(stringToPrint)
	color.Unset()
}

func GreenOutput(stringToPrint string) {
	if LoggingOutputReciever != nil {
		LoggingOutputReciever <- stringToPrint
	}
	green := color.New(color.FgGreen)
	green.Printf(stringToPrint)
	color.Unset()
}

func YellowOutput(stringToPrint string) {
	if LoggingOutputReciever != nil {
		LoggingOutputReciever <- stringToPrint
	}
	yellow := color.New(color.FgYellow)
	yellow.Printf(stringToPrint)
	color.Unset()
}

func RedOutput(stringToPrint string) {
	if LoggingOutputReciever != nil {
		LoggingOutputReciever <- stringToPrint
	}
	red := color.New(color.FgRed)
	red.Printf(stringToPrint)
	color.Unset()
}

func WhiteOutput(stringToPrint string) {
	if LoggingOutputReciever != nil {
		LoggingOutputReciever <- stringToPrint
	}
	white := color.New(color.FgWhite)
	white.Printf(stringToPrint)
	color.Unset()
}

//Info outputs log line to console with green color text
func Info(stringToPrint string) {
	if CurrentLoggingLevel <= InfoLevel {
		if ColorLogLevelLabelOnly == false {
			GreenOutput(createOutputString(stringToPrint, "INFO", true))
		} else {
			WhiteOutput(fmt.Sprintf("%s:", GetTimeString()))
			GreenOutput(" INFO ")
			WhiteOutput(fmt.Sprintf("%s -> %s\n", createCallbackLabel(2), stringToPrint))
		}
	}
}

//InfoNnl outputs log line to console with green color text without newline
func InfoNnl(stringToPrint string) {
	if CurrentLoggingLevel <= InfoLevel {
		if ColorLogLevelLabelOnly == false {
			GreenOutput(createOutputString(stringToPrint, "INFO", false))
		} else {
			WhiteOutput(fmt.Sprintf("%s:", GetTimeString()))
			GreenOutput(" INFO ")
			WhiteOutput(fmt.Sprintf("%s -> %s", createCallbackLabel(2), stringToPrint))
		}
	}
}

//Info outputs log line to console with green color text
func InfoNoColor(stringToPrint string) {
	if CurrentLoggingLevel <= InfoLevel {
		WhiteOutput(createOutputString(stringToPrint, "INFO", true))
	}
}

//InfoNnl outputs log line to console with green color text without newline
func InfoNnlNoColor(stringToPrint string) {
	if CurrentLoggingLevel <= InfoLevel {
		WhiteOutput(createOutputString(stringToPrint, "INFO", false))
	}
}

//Warn outputs log line to console with yellow color text
func Warn(stringToPrint string) {
	if CurrentLoggingLevel <= WarnLevel {
		if ColorLogLevelLabelOnly == false {
			YellowOutput(createOutputString(stringToPrint, "WARN", true))
		} else {
			WhiteOutput(fmt.Sprintf("%s:", GetTimeString()))
			YellowOutput(" WARN ")
			WhiteOutput(fmt.Sprintf("%s -> %s\n", createCallbackLabel(2), stringToPrint))
		}
	}
}

//WarnNnl outputs log line to console with yellow color text without newline
func WarnNnl(stringToPrint string) {
	if CurrentLoggingLevel <= WarnLevel {
		if ColorLogLevelLabelOnly == false {
			YellowOutput(createOutputString(stringToPrint, "WARN", false))
		} else {
			WhiteOutput(fmt.Sprintf("%s:", GetTimeString()))
			YellowOutput(" WARN ")
			WhiteOutput(fmt.Sprintf("%s -> %s", createCallbackLabel(2), stringToPrint))
		}
	}
}

//Warn outputs log line to console with yellow color text
func WarnNoColor(stringToPrint string) {
	if CurrentLoggingLevel <= WarnLevel {
		WhiteOutput(createOutputString(stringToPrint, "WARN", true))
	}
}

//DebugNnl outputs log line to console with green color text without newline
func WarnNnlNoColor(stringToPrint string) {
	if CurrentLoggingLevel <= DebugLevel {
		WhiteOutput(createOutputString(stringToPrint, "WARN", false))
	}
}

//Debug outputs log line to console with yellow color text
func Debug(stringToPrint string) {
	if CurrentLoggingLevel <= DebugLevel {
		if ColorLogLevelLabelOnly == false {
			YellowOutput(createOutputString(stringToPrint, "DEBUG", true))
		} else {
			WhiteOutput(fmt.Sprintf("%s:", GetTimeString()))
			YellowOutput(" DEBUG ")
			WhiteOutput(fmt.Sprintf("%s -> %s\n", createCallbackLabel(2), stringToPrint))
		}
	}
}

//DebugNnl outputs log line to console with yellow color text without newline
func DebugNnl(stringToPrint string) {
	if CurrentLoggingLevel <= DebugLevel {
		if ColorLogLevelLabelOnly == false {
			YellowOutput(createOutputString(stringToPrint, "DEBUG", false))
		} else {
			WhiteOutput(fmt.Sprintf("%s:", GetTimeString()))
			YellowOutput(" DEBUG ")
			WhiteOutput(fmt.Sprintf("%s -> %s", createCallbackLabel(2), stringToPrint))
		}
	}
}

//Debug outputs log line to console with yellow color text
func DebugNoColor(stringToPrint string) {
	if CurrentLoggingLevel <= DebugLevel {
		WhiteOutput(createOutputString(stringToPrint, "DEBUG", true))
	}
}

//DebugNnl outputs log line to console with yellow color text without newline
func DebugNnlNoColor(stringToPrint string) {
	if CurrentLoggingLevel <= DebugLevel {
		WhiteOutput(createOutputString(stringToPrint, "DEBUG", false))
	}
}

//Error outputs log line to console with red color text
func Error(stringToPrint string) {
	if ColorLogLevelLabelOnly == false {
		RedOutput(createOutputString(stringToPrint, "ERROR", true))
	} else {
		WhiteOutput(fmt.Sprintf("%s:", GetTimeString()))
		RedOutput(" ERROR ")
		WhiteOutput(fmt.Sprintf("%s -> %s\n", createCallbackLabel(2), stringToPrint))
	}
}

//ErrorNnl outputs log line to console with red color text without newline
func ErrorNnl(stringToPrint string) {
	if ColorLogLevelLabelOnly == false {
		RedOutput(createOutputString(stringToPrint, "ERROR", true))
	} else {
		WhiteOutput(fmt.Sprintf("%s:", GetTimeString()))
		RedOutput(" ERROR ")
		WhiteOutput(fmt.Sprintf("%s -> %s", createCallbackLabel(2), stringToPrint))
	}
}

//Error outputs log line to console with red color text
func ErrorNoColor(stringToPrint string) {
	WhiteOutput(createOutputString(stringToPrint, "ERROR", true))
}

//ErrorNnl outputs log line to console with red color text without newline
func ErrorNnlNoColor(stringToPrint string) {
	WhiteOutput(createOutputString(stringToPrint, "ERROR", true))
}

//ErrorAndExit outputs log line to console with red color text and exits
func ErrorAndExit(stringToPrint string) {
	Error(stringToPrint)
	os.Exit(1)
}

//ErrorAndExitNnl outputs the log line to the console with red color text with no newline and exits
func ErrorAndExitNnl(stringToPrint string) {
	ErrorNnl(stringToPrint)
	os.Exit(1)
}

//ErrorAndExit outputs log line to console with red color text and exits
func ErrorAndExitNoColor(stringToPrint string) {
	ErrorNoColor(stringToPrint)
	os.Exit(1)
}

//ErrorAndExitNnl outputs the log line to the console with red color text with no newline and exits
func ErrorAndExitNnlNoColor(stringToPrint string) {
	ErrorNnlNoColor(stringToPrint)
	os.Exit(1)
}

//GetTimeString gets formatted string to timestamp log and console output
func GetTimeString() string {
	t := time.Now()
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
}

func createOutputString(stp string, lvl string, nl bool) string {
	data := make([]byte, 0)
	sb := bytes.NewBuffer(data)
	if OutputDateTime {
		sb.WriteString(fmt.Sprintf("%s: ", GetTimeString()))
	}
	if OutputLogLevelFlag {
		sb.WriteString(lvl)
	}
	if OutputPath {
		sb.WriteString(fmt.Sprintf(" %s", createCallbackLabel(3)))
	}
	if OutputArrowSuffix {
		sb.WriteString(fmt.Sprintf(" -> %s", stp))
	} else {
		sb.WriteString(stp)
	}
	if nl {
		sb.WriteString("\n")
	}
	return sb.String()
}
