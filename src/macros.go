package src

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
)

func GetMacroList() []string {
	files, _ := os.ReadDir("macros/")
	var macros []string

	for _, file := range files {
		macros = append(macros, file.Name())
	}

	return macros
}

type Interpreter struct {
	Macro   string
	Actions []Action
}

type Action struct {
	Command string
	Args    []interface{}
}

func (interpreter *Interpreter) CompileToActions() []Action {
	// gets first line of the file
	file, _ := os.ReadFile(interpreter.Macro)

	var actions []Action

	for i, line := range strings.Split(string(file), "\n") {
		action := Action{}
		splitLine := strings.Split(line, "|")

		if len(splitLine) < 1 {
			continue
		}

        if splitLine[0] == "" || strings.Contains(splitLine[0], "#") || i == 0 {
            continue
        }

		action.Command = splitLine[0]

        switch action.Command {
        case "keyboard":
            switch splitLine[1] {
            case "press":
                action.Args = append(action.Args, splitLine[1], splitLine[2])
            case "release":
                action.Args = append(action.Args, splitLine[1], splitLine[2])
            case "type":
                action.Args = append(action.Args, splitLine[1])
            }
        case "delay":
            dur, _ := time.ParseDuration(strings.Trim(splitLine[1], "\n"))
            action.Args = append(action.Args, dur)
        case "mouse":
            switch splitLine[1] {
            case "set":
                x, _ := strconv.Atoi(splitLine[2])
                y, _ := strconv.Atoi(splitLine[3])
                action.Args = append(action.Args, splitLine[1], x, y)
            case "click":
                action.Args = append(action.Args, splitLine[1], splitLine[2])
            case "down":
                action.Args = append(action.Args, splitLine[1], splitLine[2])
            case "up":
                action.Args = append(action.Args, splitLine[1], splitLine[2])
            case "scroll":
                x, _ := strconv.Atoi(splitLine[2])
                y, _ := strconv.Atoi(splitLine[3])
                action.Args = append(action.Args, splitLine[1], x, y)
            }

        }

		actions = append(actions, action)
	}

	interpreter.Actions = actions
	return actions
}

func (interpreter *Interpreter) Run() {
	for _, action := range interpreter.Actions {
		switch action.Command {
		case "keyboard":
			switch action.Args[0] {
			case "press":
                robotgo.KeyToggle(action.Args[1].(string), "down")
			case "release":
                robotgo.KeyToggle(action.Args[1].(string), "up")
			case "type":
                robotgo.TypeStr(action.Args[1].(string))
			}
		case "delay":
			if dur, ok := action.Args[0].(time.Duration); ok {
				time.Sleep(dur)
			}
		case "mouse":
			switch action.Args[0] {
			case "set":
                robotgo.Move(action.Args[1].(int), action.Args[2].(int))
			case "click":
                robotgo.Click(action.Args[1].(string), false)
			case "down":
                robotgo.Toggle(action.Args[1].(string), "down")
			case "up":
                robotgo.Toggle(action.Args[1].(string), "up")
			case "scroll":
                robotgo.Scroll(action.Args[1].(int), action.Args[2].(int))
			}
		}
	}
}

