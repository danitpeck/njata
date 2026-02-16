package text

import "strings"

var colorMap = map[byte]string{
    'x': "\x1b[0m",
    'w': "\x1b[37m",
    'W': "\x1b[97m",
    'r': "\x1b[31m",
    'R': "\x1b[91m",
    'g': "\x1b[32m",
    'G': "\x1b[92m",
    'b': "\x1b[34m",
    'B': "\x1b[94m",
    'c': "\x1b[36m",
    'C': "\x1b[96m",
    'y': "\x1b[33m",
    'Y': "\x1b[93m",
    'p': "\x1b[35m",
    'P': "\x1b[95m",
    'z': "\x1b[90m",
    'Z': "\x1b[37m",
    'o': "\x1b[33m",
    'O': "\x1b[93m",
}

// TranslateSmaugColors converts SMAUG-style & codes to ANSI sequences.
// Unknown codes are stripped. "&&" becomes a literal "&".
func TranslateSmaugColors(input string) string {
    if input == "" {
        return input
    }

    var builder strings.Builder
    builder.Grow(len(input) + 8)

    for i := 0; i < len(input); i++ {
        if input[i] != '&' {
            builder.WriteByte(input[i])
            continue
        }

        if i+1 >= len(input) {
            continue
        }

        next := input[i+1]
        if next == '&' {
            builder.WriteByte('&')
            i++
            continue
        }

        if seq, ok := colorMap[next]; ok {
            builder.WriteString(seq)
            i++
            continue
        }

        i++
    }

    return builder.String()
}
