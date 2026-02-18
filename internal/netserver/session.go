package netserver

import (
    "bufio"
    "net"
    "strings"
    "sync"

    mudtext "njata/internal/text"
)

type Session struct {
    conn           net.Conn
    reader         *bufio.Reader
    writer         *bufio.Writer
    writeMu        sync.Mutex
    stateMu        sync.Mutex
    closed         bool
    disconnectOnce sync.Once
    disconnected   chan struct{}
}

func NewSession(conn net.Conn) *Session {
    return &Session{
        conn:         conn,
        reader:       bufio.NewReader(conn),
        writer:       bufio.NewWriter(conn),
        disconnected: make(chan struct{}),
    }
}

func (s *Session) Write(message string) {
    if s.IsDisconnectRequested() {
        return
    }

    s.writeMu.Lock()
    defer s.writeMu.Unlock()

    if s.IsDisconnectRequested() {
        return
    }

    translated := mudtext.TranslateSmaugColors(message)
    _, _ = s.writer.WriteString(translated)
    _ = s.writer.Flush()
}

func (s *Session) WriteLine(text string) {
    s.Write(text + "\r\n")
}

func (s *Session) ReadLine() (string, error) {
    line, err := s.reader.ReadString('\n')
    if err != nil && len(line) == 0 {
        return "", err
    }

    line = strings.TrimRight(line, "\r\n")
    return line, err
}

func (s *Session) RequestDisconnect(reason string) {
    s.disconnectOnce.Do(func() {
        s.stateMu.Lock()
        s.closed = true
        s.stateMu.Unlock()

        close(s.disconnected)
        _ = s.conn.Close()
    })
}

func (s *Session) IsDisconnectRequested() bool {
    select {
    case <-s.disconnected:
        return true
    default:
        return false
    }
}

func (s *Session) Close() {
    s.RequestDisconnect("closed")
}
