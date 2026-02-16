package area

import (
    "bufio"
    "fmt"
    "os"
    "path/filepath"
    "sort"
    "strconv"
    "strings"

    "njata/internal/game"
)

func LoadRoomsFromDir(path string) (map[int]*game.Room, int, error) {
    entries, err := os.ReadDir(path)
    if err != nil {
        return nil, 0, err
    }

    rooms := map[int]*game.Room{}
    for _, entry := range entries {
        if entry.IsDir() {
            continue
        }
        if !strings.HasSuffix(strings.ToLower(entry.Name()), ".are") {
            continue
        }

        filePath := filepath.Join(path, entry.Name())
        parsed, err := parseRoomsFromFile(filePath)
        if err != nil {
            return nil, 0, fmt.Errorf("%s: %w", entry.Name(), err)
        }

        for vnum, room := range parsed {
            rooms[vnum] = room
        }
    }

    start := findLowestVnum(rooms)
    return rooms, start, nil
}

func findLowestVnum(rooms map[int]*game.Room) int {
    if len(rooms) == 0 {
        return 0
    }

    vnums := make([]int, 0, len(rooms))
    for vnum := range rooms {
        vnums = append(vnums, vnum)
    }
    sort.Ints(vnums)
    return vnums[0]
}

func parseRoomsFromFile(path string) (map[int]*game.Room, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

    rooms := map[int]*game.Room{}
    var current *game.Room

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }

        if line == "#ROOM" {
            current = &game.Room{Exits: map[string]int{}, ExDescs: map[string]string{}}
            continue
        }

        if current == nil {
            continue
        }

        switch {
        case line == "#ENDROOM":
            if current.Vnum > 0 {
                rooms[current.Vnum] = current
            }
            current = nil
        case strings.HasPrefix(line, "Vnum"):
            vnum, err := parseIntField(line)
            if err != nil {
                return nil, err
            }
            current.Vnum = vnum
        case strings.HasPrefix(line, "Name"):
            name, err := readTildeText(fieldRemainder(line, "Name"), scanner)
            if err != nil {
                return nil, err
            }
            current.Name = name
        case strings.HasPrefix(line, "Desc"):
            desc, err := readTildeText(fieldRemainder(line, "Desc"), scanner)
            if err != nil {
                return nil, err
            }
            current.Description = desc
        case line == "#EXIT":
            direction, toRoom, err := readExit(scanner)
            if err != nil {
                return nil, err
            }
            if direction != "" && toRoom > 0 {
                current.Exits[direction] = toRoom
            }
        case line == "#EXDESC":
            keys, desc, err := readExDesc(scanner)
            if err != nil {
                return nil, err
            }
            for _, key := range keys {
                current.ExDescs[key] = desc
            }
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, err
    }

    return rooms, nil
}

func readExDesc(scanner *bufio.Scanner) ([]string, string, error) {
    var keys []string
    var desc string

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "#ENDEXDESC" {
            return keys, desc, nil
        }

        switch {
        case strings.HasPrefix(line, "ExDescKey"):
            value, err := readTildeText(fieldRemainder(line, "ExDescKey"), scanner)
            if err != nil {
                return nil, "", err
            }
            keys = normalizeKeys(value)
        case strings.HasPrefix(line, "ExDesc"):
            value, err := readTildeText(fieldRemainder(line, "ExDesc"), scanner)
            if err != nil {
                return nil, "", err
            }
            desc = value
        }
    }

    if err := scanner.Err(); err != nil {
        return nil, "", err
    }

    return keys, desc, nil
}

func normalizeKeys(value string) []string {
    fields := strings.Fields(strings.ToLower(value))
    keys := make([]string, 0, len(fields))
    for _, field := range fields {
        if field != "" {
            keys = append(keys, field)
        }
    }
    return keys
}

func readExit(scanner *bufio.Scanner) (string, int, error) {
    var direction string
    var toRoom int

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "#ENDEXIT" {
            return direction, toRoom, nil
        }

        switch {
        case strings.HasPrefix(line, "Direction"):
            value, err := readTildeText(fieldRemainder(line, "Direction"), scanner)
            if err != nil {
                return "", 0, err
            }
            direction = strings.ToLower(strings.TrimSpace(value))
        case strings.HasPrefix(line, "ToRoom"):
            parsed, err := parseIntField(line)
            if err != nil {
                return "", 0, err
            }
            toRoom = parsed
        }
    }

    if err := scanner.Err(); err != nil {
        return "", 0, err
    }

    return direction, toRoom, nil
}

func parseIntField(line string) (int, error) {
    fields := strings.Fields(line)
    if len(fields) < 2 {
        return 0, fmt.Errorf("invalid numeric field: %s", line)
    }
    value, err := strconv.Atoi(fields[1])
    if err != nil {
        return 0, fmt.Errorf("invalid number in line: %s", line)
    }
    return value, nil
}

func fieldRemainder(line string, prefix string) string {
    return strings.TrimSpace(strings.TrimPrefix(line, prefix))
}

func readTildeText(initial string, scanner *bufio.Scanner) (string, error) {
    if idx := strings.Index(initial, "~"); idx >= 0 {
        return strings.TrimRight(strings.TrimRight(initial[:idx], "\r"), " "), nil
    }

    var builder strings.Builder
    trimmed := strings.TrimRight(strings.TrimRight(initial, "\r"), " ")
    if trimmed != "" {
        appendLine(&builder, trimmed)
    }

    for scanner.Scan() {
        line := strings.TrimRight(scanner.Text(), "\r")
        if idx := strings.Index(line, "~"); idx >= 0 {
            segment := strings.TrimRight(line[:idx], " ")
            if segment != "" {
                appendLine(&builder, segment)
            }
            return builder.String(), nil
        }
        appendLine(&builder, strings.TrimRight(line, " "))
    }

    if err := scanner.Err(); err != nil {
        return builder.String(), err
    }

    return builder.String(), nil
}

func appendLine(builder *strings.Builder, line string) {
    if builder.Len() > 0 {
        builder.WriteString("\n")
    }
    builder.WriteString(line)
}
