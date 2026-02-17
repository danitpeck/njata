package main

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"strings"
	"time"
)

type TestClient struct {
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewTestClient(addr string) (*TestClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, err
	}

	return &TestClient{
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}, nil
}

func (tc *TestClient) Close() error {
	return tc.conn.Close()
}

// ReadResponse reads all available data from the server with a timeout
func (tc *TestClient) ReadResponse(timeout time.Duration) string {
	tc.conn.SetReadDeadline(time.Now().Add(timeout))
	defer tc.conn.SetReadDeadline(time.Time{})

	var response strings.Builder
	buffer := make([]byte, 4096)

	for {
		n, err := tc.conn.Read(buffer)
		if n > 0 {
			response.Write(buffer[:n])
		}
		if err != nil {
			break
		}
	}

	return response.String()
}

func (tc *TestClient) SendCommand(cmd string) error {
	_, err := tc.writer.WriteString(cmd + "\n")
	if err != nil {
		return err
	}
	return tc.writer.Flush()
}

func (tc *TestClient) CommandResponse(cmd string) string {
	tc.SendCommand(cmd)
	time.Sleep(500 * time.Millisecond) // Give server time to process
	return tc.ReadResponse(1 * time.Second)
}

func contains(response string, keywords []string) bool {
	lower := strings.ToLower(response)
	for _, kw := range keywords {
		if strings.Contains(lower, strings.ToLower(kw)) {
			return true
		}
	}
	return false
}

func main() {
	addr := "localhost:4000"

	fmt.Println("=== NJATA MUD Integration Tests ===")

	// Connect
	fmt.Println("[TEST] Connecting to server...")
	client, err := NewTestClient(addr)
	if err != nil {
		fmt.Printf("❌ Failed to connect: %v\n", err)
		return
	}
	defer client.Close()
	fmt.Println("✓ Connected")
	fmt.Println()

	// Read banner
	fmt.Println("[TEST] Receiving banner...")
	banner := client.ReadResponse(2 * time.Second)
	if strings.Contains(banner, "|") || strings.Contains(banner, "---") || len(banner) > 50 {
		fmt.Println("✓ Banner received")
	} else {
		fmt.Printf("⚠️  Unexpected banner\n")
	}
	fmt.Println()

	// Login
	fmt.Println("[TEST] Logging in as Vex...")
	response := client.CommandResponse("vex")
	if contains(response, []string{"welcome", "vex"}) {
		fmt.Println("✓ Login successful")
	} else {
		fmt.Printf("Response: %s\n", response[:min(len(response), 200)])
	}
	fmt.Println()

	// Look command
	fmt.Println("[TEST] Testing 'look' command...")
	response = client.CommandResponse("look")
	checks := []struct {
		name    string
		markers []string
	}{
		{"Room title", []string{"lyceum", "darkhaven", "academy", "the"}},
		{"Exits", []string{"exits", "north", "south", "east", "west"}},
	}
	for _, check := range checks {
		if contains(response, check.markers) {
			fmt.Printf("  ✓ %s found\n", check.name)
		} else {
			fmt.Printf("  ❌ %s missing\n", check.name)
		}
	}
	fmt.Println()

	// Stats command
	fmt.Println("[TEST] Testing 'stats' command...")
	response = client.CommandResponse("stats")
	if contains(response, []string{"str", "int", "vex"}) {
		fmt.Println("  ✓ Stats display working")
	} else {
		fmt.Printf("  ❌ Stats incomplete\n")
	}
	fmt.Println()

	// Who command
	fmt.Println("[TEST] Testing 'who' command...")
	response = client.CommandResponse("who")
	if contains(response, []string{"vex", "playing"}) {
		fmt.Println("  ✓ Who list working")
	} else {
		fmt.Printf("  ⚠️  Who response unexpected\n")
	}
	fmt.Println()

	// Say command
	fmt.Println("[TEST] Testing 'say' command...")
	response = client.CommandResponse("say test works")
	if strings.Contains(response, "test works") || strings.Contains(response, "say") {
		fmt.Println("  ✓ Say command working")
	} else {
		fmt.Printf("  ⚠️  Say response unclear\n")
	}
	fmt.Println()

	// Score command
	fmt.Println("[TEST] Testing 'score' command...")
	response = client.CommandResponse("score")
	if contains(response, []string{"experience", "level", "vex"}) {
		fmt.Println("  ✓ Score display working")
	} else {
		fmt.Printf("  ⚠️  Score incomplete\n")
	}
	fmt.Println()

	// Spellbook command
	fmt.Println("[TEST] Testing 'spellbook' command...")
	response = client.CommandResponse("spellbook")
	if contains(response, []string{"spellbook", "fireball", "mana"}) || contains(response, []string{"spell", "magic"}) {
		fmt.Println("  ✓ Spellbook display working")
	} else {
		fmt.Printf("  ⚠️  Spellbook incomplete\n")
	}
	fmt.Println()

	// Cast command
	fmt.Println("[TEST] Testing 'cast' command...")
	response = client.CommandResponse("cast fireball")
	if contains(response, []string{"cast", "fireball", "mana"}) {
		fmt.Println("  ✓ Cast command working")
	} else {
		fmt.Printf("  ⚠️  Cast response unexpected\n")
	}
	fmt.Println()

	fmt.Println("=== Tests Complete ===")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
