package commands

import (
    "sort"
    "strings"
)

type Handler func(Context, string)

type Registry struct {
    handlers map[string]Handler
    ordered []string // maintain registration order for consistent prefix matching
}

func NewRegistry() *Registry {
    return &Registry{
        handlers: map[string]Handler{},
        ordered:  []string{},
    }
}

// StringPrefix checks if str is a prefix match for target (case-insensitive)
func StringPrefix(str, target string) bool {
    return len(str) <= len(target) &&
           strings.EqualFold(str, target[:len(str)])
}

func (r *Registry) Register(name string, handler Handler) {
    lower := strings.ToLower(name)
    r.handlers[lower] = handler
    r.ordered = append(r.ordered, lower)
}

func (r *Registry) Execute(ctx Context, command string, args string) bool {
    lower := strings.ToLower(command)
    
    // First check for exact match
    if handler, ok := r.handlers[lower]; ok {
        handler(ctx, args)
        return true
    }
    
    // Then check for prefix match (in registration order)
    for _, name := range r.ordered {
        if StringPrefix(lower, name) {
            handler := r.handlers[name]
            handler(ctx, args)
            return true
        }
    }
    
    return false
}

func (r *Registry) List() []string {
    names := make([]string, 0, len(r.handlers))
    for name := range r.handlers {
        names = append(names, name)
    }

    sort.Strings(names)
    return names
}
