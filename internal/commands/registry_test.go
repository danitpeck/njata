package commands

import (
	"testing"
)

func TestStringPrefix(t *testing.T) {
	tests := []struct {
		str      string
		target   string
		expected bool
	}{
		{"look", "look", true},
		{"l", "look", true},
		{"lo", "look", true},
		{"loo", "look", true},
		{"look", "look", true},
		{"lookx", "look", false},
		{"x", "look", false},
		{"L", "look", true},      // case insensitive
		{"LO", "LOOK", true},     // case insensitive
		{"", "look", true},       // empty string is prefix of anything
		{"say", "look", false},
	}

	for _, tt := range tests {
		result := StringPrefix(tt.str, tt.target)
		if result != tt.expected {
			t.Errorf("StringPrefix(%q, %q) = %v, want %v", tt.str, tt.target, result, tt.expected)
		}
	}
}

func TestRegistryPrefixMatching(t *testing.T) {
	reg := NewRegistry()
	
	// Track which commands were called
	var called []string
	
	reg.Register("look", func(ctx Context, args string) {
		called = append(called, "look")
	})
	reg.Register("say", func(ctx Context, args string) {
		called = append(called, "say")
	})
	reg.Register("stats", func(ctx Context, args string) {
		called = append(called, "stats")
	})
	
	// Create a dummy context
	ctx := Context{}
	
	// Test exact match
	called = []string{}
	reg.Execute(ctx, "look", "")
	if len(called) != 1 || called[0] != "look" {
		t.Errorf("Exact match failed: got %v", called)
	}
	
	// Test prefix match - single letter
	called = []string{}
	reg.Execute(ctx, "l", "")
	if len(called) != 1 || called[0] != "look" {
		t.Errorf("Prefix 'l' should match 'look': got %v", called)
	}
	
	// Test prefix match - "lo"
	called = []string{}
	reg.Execute(ctx, "lo", "")
	if len(called) != 1 || called[0] != "look" {
		t.Errorf("Prefix 'lo' should match 'look': got %v", called)
	}
	
	// Test prefix match - "loo"
	called = []string{}
	reg.Execute(ctx, "loo", "")
	if len(called) != 1 || called[0] != "look" {
		t.Errorf("Prefix 'loo' should match 'look': got %v", called)
	}
	
	// Test "sa" matches "say"
	called = []string{}
	reg.Execute(ctx, "sa", "")
	if len(called) != 1 || called[0] != "say" {
		t.Errorf("Prefix 'sa' should match 'say': got %v", called)
	}
	
	// Test "st" matches "stats"
	called = []string{}
	reg.Execute(ctx, "st", "")
	if len(called) != 1 || called[0] != "stats" {
		t.Errorf("Prefix 'st' should match 'stats': got %v", called)
	}
	
	// Test case insensitivity
	called = []string{}
	reg.Execute(ctx, "LOOK", "")
	if len(called) != 1 || called[0] != "look" {
		t.Errorf("Case insensitive 'LOOK' should match 'look': got %v", called)
	}
	
	// Test non-matching command
	called = []string{}
	result := reg.Execute(ctx, "xyz", "")
	if result || len(called) != 0 {
		t.Errorf("Non-matching command 'xyz' should return false")
	}
}
