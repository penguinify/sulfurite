package src

import (
	"os"
	"strconv"
	"unicode"

)

func GetMacroList(path string) []string {
	files, _ := os.ReadDir(path)
	var macros []string

	for _, file := range files {
		macros = append(macros, file.Name())
	}

	return macros
}


type TokenType int

/*
MACROSYNTAX.md
# TokenFunction
Possible Values:
- type <string> <ms as number>
- mouseset <number> <number>
- mouseclick <left|right|middle>
- mousemove <number> <number>
- keydown <string>
- keyup <string>
- delay <ms as number>
# Number
any positive integer, ex: delay 1000
# String
any string, ex: keydown "a".    
Must be double quoted.
# TokenKeyword
Possible Values:
- loop <number>
- forever
- end

*/

const (
    TokenFunction TokenType = iota
    TokenKeyword
    TokenString
    TokenNumber
)

var Keywords = [...]string{"loop", "end", "forever"}

var HigherLevelKeywords = [...]string{"loop", "forever", "root"}

type Token struct {
    Type TokenType
    Value string
}

type Parser struct {
    Tokens []Token
    pos int
}

func IsInSlice(array []string, value string) bool {
    for _, v := range array {
        if v == value {
            return true
        }
    }

    return false
}

type Lexer struct {
    Source string
    pos int
}

type ASTNode struct {
    Type TokenType
    Value string
    Children []*ASTNode
}

func NewParser(input string) *Parser {
    lexer := NewLexer(input)
    tokens := lexer.Tokenize()
    return &Parser{Tokens: tokens}
}

func NewLexer(source string) *Lexer {
    return &Lexer{Source: source}
}

func (l *Lexer) Tokenize() []Token {
    var tokens []Token

    for l.pos < len(l.Source) {

        char := l.Source[l.pos]

        switch {
        case unicode.IsSpace(rune(char)):
            l.pos++
            continue
        case char == '"':
            tokens = append(tokens, Token{Type: TokenString, Value: l.String()})
        case unicode.IsDigit(rune(char)):
            tokens = append(tokens, Token{Type: TokenNumber, Value: strconv.Itoa(l.Number())})
        default:
            value := l.readWord()

            if IsInSlice(Keywords[:], value) {
                tokens = append(tokens, Token{Type: TokenKeyword, Value: value})
            } else {
                tokens = append(tokens, Token{Type: TokenFunction, Value: value})
            }

        }
        l.pos++
    }


    return tokens
}



func (l *Lexer) readWord() string {
    var keyword string

    for l.pos < len(l.Source) {
        if unicode.IsSpace(rune(l.Source[l.pos])) {
            break
        }

        keyword += string(l.Source[l.pos])
        l.pos++
    }

    return keyword
}

func (l *Lexer) String() string {
    var str string

    l.pos++

    for l.pos < len(l.Source) {
        if l.Source[l.pos] == '"' {
            break
        }

        str += string(l.Source[l.pos])
        l.pos++
    }

    return str
}

func (l *Lexer) Number() int {
    var num string

    for l.pos < len(l.Source) {
        if unicode.IsSpace(rune(l.Source[l.pos])) {
            break
        }

        num += string(l.Source[l.pos])
        l.pos++
    }

    n, _ := strconv.Atoi(num)
    return n
}

func (p *Parser) Parse() *ASTNode {
    parent := &ASTNode{Type: 0, Value: "root", Children: []*ASTNode{}}

    for p.pos < len(p.Tokens) {
        token := p.Tokens[p.pos]
        
        switch token.Type {
        case TokenFunction:
            parent.Children = append(parent.Children, &ASTNode{Type: TokenFunction, Value: token.Value})
            p.pos++
        case TokenKeyword:
            keywordNode := &ASTNode{Type: TokenKeyword, Value: token.Value, Children: []*ASTNode{}}
            parent.Children = append(parent.Children, keywordNode)
            p.pos++
            switch token.Value {
            case "loop":
                
                parent.Children = append(parent.Children, &ASTNode{Type: TokenNumber, Value: p.Tokens[p.pos].Value})
                p.pos++

                keywordNode.Children = p.Parse().Children
                
                if p.pos >= len(p.Tokens) {
                    break
                }

            case "forever":
                for p.pos < len(p.Tokens) {
                    keywordNode.Children = p.Parse().Children
                }

            case "end":
                return parent
            }


        case TokenNumber:
            parent.Children = append(parent.Children, &ASTNode{Type: TokenNumber, Value: token.Value})
            p.pos++
        case TokenString:
            parent.Children = append(parent.Children, &ASTNode{Type: TokenString, Value: token.Value})
            p.pos++
        }
    }

    return parent
}

