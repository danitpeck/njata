package commands

import (
    "sort"
    "strings"
)

type Handler func(Context, string)

type Registry struct {
    handlers map[string]Handler
}

func NewRegistry() *Registry {
    return &Registry{handlers: map[string]Handler{}}
}

func (r *Registry) Register(name string, handler Handler) {
    r.handlers[strings.ToLower(name)] = handler
}

func (r *Registry) Execute(ctx Context, command string, args string) bool {
    handler, ok := r.handlers[strings.ToLower(command)]
    if !ok {
        return false
    }

    handler(ctx, args)
    return true
}

func (r *Registry) List() []string {
    names := make([]string, 0, len(r.handlers))
    for name := range r.handlers {
        names = append(names, name)
    }

    sort.Strings(names)
    return names
}
