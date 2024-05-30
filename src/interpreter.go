package src

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-vgo/robotgo"
)

func Interpret(ast *ASTNode, interupted <-chan bool) (bool, error) {

    if !IsInSlice(HigherLevelKeywords[:], ast.Value) {
        return false, errors.New("Invalid root node, Must be a higher level keyword")
    }

    pos := 0

    for pos < len(ast.Children) {

        child := ast.Children[pos]

        switch TokenType(child.Type) {
        case TokenFunction:
            switch child.Value {
            case "mouseset":
                pos++
                x, _ := strconv.Atoi(ast.Children[pos].Value)
                pos++
                y, _ := strconv.Atoi(ast.Children[pos].Value)
                robotgo.Move(x, y)
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
                robotgo.KeyToggle(key, "down")
            case "keyup":
                pos++
                key := ast.Children[pos].Value
                robotgo.KeyToggle(key, "up")
            case "type":
                pos++
                text := ast.Children[pos].Value
                pos++
                speed, _ := strconv.Atoi(ast.Children[pos].Value)

                robotgo.TypeStr(text, speed)

            case "sleep":
                pos++
                duration, _ := strconv.Atoi(ast.Children[pos].Value)
                time.Sleep(time.Duration(duration) * time.Millisecond)
            default:
            }
        case TokenKeyword:
            switch child.Value {
            case "loop":
                pos++
                count, _ := strconv.Atoi(ast.Children[pos].Value)
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
                return true, nil
            default:
            }
        case TokenString:
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

    <-interupted
    return true, nil
}
