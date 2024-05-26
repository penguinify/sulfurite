package src

import (
	"errors"
	"strconv"
	"time"

	"github.com/go-vgo/robotgo"
)

func Interpret(ast *ASTNode, interupted <-chan bool) (bool, error) {
    if ast.Value != "root" {
        return false, errors.New("Invalid AST, Node must be root")
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
        case TokenString:
        case TokenNumber:

        default:
        }

        pos++

        select {
        case <-interupted:
            return false, nil
        default:
            continue
        }
    }   

    <- interupted
    return true, nil
}
