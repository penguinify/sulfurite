package src

import (
	"errors"
	"math/rand/v2"
	"strconv"
	"time"

	"github.com/go-vgo/robotgo"

    "sulfurite/utils"
)


func ConvertNumber(n interface{}) int {
    switch n.(type) {
    case string:
        num, _ := strconv.Atoi(n.(string))
        return num
    case Range:
        return rand.IntN(n.(Range).End+1 - n.(Range).Start)
    default:
        return 0
    }
}

func Interpret(ast *ASTNode, interupted <-chan bool) (bool, error) {
    var interpreterVersion = utils.LoadConfig("config.json").MacroInterpreterVersion

    if !IsInSlice(HigherLevelKeywords[:], ast.Value.(string)) {
        return false, errors.New("Invalid root node, Must be a higher level keyword")
    } else if ast.Value == "root" {
        
        if ast.Children[0].Value != interpreterVersion {
            return false, errors.New("Invalid interpreter version! Please update or downgrade to avoid compatibility issues. (Current version: " + interpreterVersion + ", Macro version: " + ast.Children[0].Value.(string) + ")")
        }
    }

    pos := 0

    for pos < len(ast.Children) {

        child := ast.Children[pos]

        switch TokenType(child.Type) {
        case TokenFunction:
            switch child.Value {
            case "mouseset":
                pos++
                x := ConvertNumber(ast.Children[pos].Value)
                pos++
                y := ConvertNumber(ast.Children[pos].Value)
                pos++
                smooth := ConvertNumber(ast.Children[pos].Value)

                if smooth == 0 {
                    robotgo.Move(x, y)
                } else {
                    pos++
                    time := ConvertNumber(ast.Children[pos].Value)
                    pos++
                    speed := ConvertNumber(ast.Children[pos].Value)
                    robotgo.MoveSmooth(x, y, float64(time), float64(speed))
                }
            case "mousemove":
                pos++
                x := ConvertNumber(ast.Children[pos].Value)
                pos++
                y := ConvertNumber(ast.Children[pos].Value)
                pos++
                smooth := ConvertNumber(ast.Children[pos].Value)
                
                if smooth == 0 {
                    robotgo.MoveRelative(x, y)
                } else {
                    pos++
                    time := ConvertNumber(ast.Children[pos].Value)
                    pos++
                    speed := ConvertNumber(ast.Children[pos].Value)
                    robotgo.MoveSmoothRelative(x, y, float64(time), float64(speed))
                }
            case "scroll":
                pos++
                x := ConvertNumber(ast.Children[pos].Value)
                pos++
                y := ConvertNumber(ast.Children[pos].Value)
                pos++
                smooth := ConvertNumber(ast.Children[pos].Value)

                if smooth == 0 {
                    robotgo.Scroll(x, y)
                } else {
                    pos++
                    time := ConvertNumber(ast.Children[pos].Value)
                    pos++
                    speed := ConvertNumber(ast.Children[pos].Value)
                    robotgo.ScrollSmooth(x, y, time, speed)
                }
            case "drag":
                pos++
                x := ConvertNumber(ast.Children[pos].Value)
                pos++
                y := ConvertNumber(ast.Children[pos].Value)

                robotgo.DragSmooth(x, y)
            case "mousedown":
                pos++
                button := ast.Children[pos].Value
                robotgo.Click(button, true)
            case "mouseup":
                pos++
                button := ast.Children[pos].Value
                robotgo.Click(button, false)
            case "keydown":
                pos++
                key := ast.Children[pos].Value
                robotgo.KeyToggle(key.(string), "down")
            case "keyup":
                pos++
                key := ast.Children[pos].Value
                robotgo.KeyToggle(key.(string), "up")
            case "keytap":
                pos++
                key := ast.Children[pos].Value
                robotgo.KeyTap(key.(string))
            case "type":
                pos++
                text := ast.Children[pos].Value

                robotgo.TypeStr(text.(string))

            case "sleep":
                pos++
                duration := ConvertNumber(ast.Children[pos].Value)
                time.Sleep(time.Duration(duration) * time.Millisecond)
            default:
            }
        case TokenKeyword:
            switch child.Value {
            case "loop":
                pos++
                count := ConvertNumber(ast.Children[pos].Value)
                for i := 0; i < count; i++ {
                    _, err := Interpret(child, interupted)
                    if err != nil {
                        return false, err
                    }

                }
            case "forever":
                for {
                    _, err := Interpret(child, interupted)
                    if err != nil {
                        return false, err
                    }
                }
            case "end":
                if ast.Value == "loop" {
                    return true, nil
                }
            default:
            }
        case TokenString:
        case TokenRandomNumber:
        case TokenNumber:

        default:
        }

        pos++

        select {
        case <-interupted:
            return false, errors.New("Interupted")
        default:
        }

    }   

    if ast.Value == "root" {
        <-interupted
    }
    return true, nil
}
