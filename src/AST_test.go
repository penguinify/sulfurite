package src

import (
	"testing"
)

func TestGetMacroList(t *testing.T) {
	// Test with valid path
	macros := GetMacroList("./")
	if len(macros) == 0 {
		t.Errorf("Expected macros, got none")
	}

	// Test with invalid path
	macros = GetMacroList("./nonexistent")
	if len(macros) != 0 {
		t.Errorf("Expected no macros, got some")
	}
}

func TestNewParser(t *testing.T) {
	parser := NewParser("test")
	if parser == nil {
		t.Errorf("Expected a Parser, got nil")
	}
}

func TestNewLexer(t *testing.T) {
	lexer := NewLexer("test")
	if lexer == nil {
		t.Errorf("Expected a Lexer, got nil")
	}
}

func TestTokenize(t *testing.T) {
	lexer := NewLexer("test")
	tokens := lexer.Tokenize()
	if len(tokens) == 0 {
		t.Errorf("Expected tokens, got none")
	}
}


func TestReadWord(t *testing.T) {
	lexer := NewLexer("test")
	word := lexer.readWord()
	if word != "test" {
		t.Errorf("Expected 'test', got '%s'", word)
	}
}

func TestString(t *testing.T) {
	lexer := NewLexer("\"test\"")
	str := lexer.String()
	if str != "test" {
		t.Errorf("Expected 'test', got '%s'", str)
	}
}

func TestNumber(t *testing.T) {
	lexer := NewLexer("123")
	num := lexer.Number()
	if num != 123 {
		t.Errorf("Expected 123, got %d", num)
	}
}

func TestParse(t *testing.T) {
	parser := NewParser("loop 3\nmouseset 100 100")
	ast := parser.Parse()
	if ast == nil {
		t.Errorf("Expected an AST, got nil")
	}
}

func BenchmarkParse(b *testing.B) {
    parser := NewParser("loop 3")
    for i := 0; i < b.N; i++ {
        parser.Parse()
    }
}

func BenchmarkTokenize(b *testing.B) {
    lexer := NewLexer("loop 3")
    for i := 0; i < b.N; i++ {
        lexer.Tokenize()
    }
}
