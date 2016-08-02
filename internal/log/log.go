package log

import (
	"fmt"
	"github.com/ttacon/chalk"
	"strings"
	"os"
)

var logctx []string

// push logging state
func PushState(s string) {
	logctx = append(logctx, s)
}

// pop logging state
func PopState() {
	logctx = logctx[:len(logctx)-1]
}

func Infof(s string, v ...interface{}) {
	PushState(s)
	fmt.Printf(chalk.Green.Color(strings.Join(logctx, ": ")), v...)
	PopState()
}

func Infoln(s string) {
	PushState(s)
	fmt.Println(chalk.Green.Color(strings.Join(logctx, ": ")))
	PopState()
}

func Warnf(s string, v ...interface{}) {
	PushState(s)
	fmt.Printf(chalk.Yellow.Color(strings.Join(logctx, ": ")), v...)
	PopState()
}

func Warnln(s string) {
	PushState(s)
	fmt.Println(chalk.Yellow.Color(strings.Join(logctx, ": ")))
	PopState()
}

func Fatal(s string) {
	PushState(s)
	fmt.Print(chalk.Red.Color(strings.Join(logctx, ": ")))
	os.Exit(1)
}

func Fatalf(s string, v ...interface{}) {
	PushState(s)
	fmt.Printf(chalk.Red.Color(strings.Join(logctx, ": ")), v...)
	os.Exit(1)
}
