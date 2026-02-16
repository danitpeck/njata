package parser

import "strings"

func ParseInput(line string) (string, string) {
    trimmed := strings.TrimSpace(line)
    if trimmed == "" {
        return "", ""
    }

    fields := strings.Fields(trimmed)
    if len(fields) == 0 {
        return "", ""
    }

    command := strings.ToLower(fields[0])
    if len(fields) == 1 {
        return command, ""
    }

    index := strings.Index(trimmed, fields[1])
    if index == -1 {
        return command, ""
    }

    args := strings.TrimSpace(trimmed[index:])
    return command, args
}
