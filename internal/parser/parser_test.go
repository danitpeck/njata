package parser

import "testing"

func TestParseInput(t *testing.T) {
    command, args := ParseInput("  SAY   hello world  ")
    if command != "say" {
        t.Fatalf("expected command 'say', got '%s'", command)
    }
    if args != "hello world" {
        t.Fatalf("expected args 'hello world', got '%s'", args)
    }
}

func TestParseInputEmpty(t *testing.T) {
    command, args := ParseInput("   ")
    if command != "" || args != "" {
        t.Fatalf("expected empty command and args, got '%s' and '%s'", command, args)
    }
}
