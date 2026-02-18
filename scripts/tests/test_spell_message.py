#!/usr/bin/env python3
"""Quick test to verify spell combat message formatting is fixed."""
import socket
import time


def send_command(sock, command):
    """Send a command and get the response."""
    sock.sendall((command + "\n").encode('utf-8'))
    time.sleep(0.3)
    response = b''
    while True:
        try:
            sock.settimeout(0.1)
            chunk = sock.recv(4096)
            if not chunk:
                break
            response += chunk
        except socket.timeout:
            break
    return response.decode('utf-8', errors='ignore')


def main():
    host = 'localhost'
    port = 4000

    print("Testing spell combat message formatting...")
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect((host, port))

    # Login as Vex
    send_command(sock, "vex")
    send_command(sock, "test")

    # Teleport to goblin entrance
    send_command(sock, "teleport 20000")

    # Cast arcane bolt at sentry
    print("\n=== Testing: cast arcane bolt sentry ===")
    response = send_command(sock, "cast arcane bolt sentry")
    
    # Check for the problematic period before "for"
    lines = response.split('\n')
    for line in lines:
        if 'arcane bolt' in line.lower() and 'damage' in line.lower():
            print(f"Combat message: {line.strip()}")
            if '. for' in line:
                print("❌ FAIL: Double period found ('. for')")
            else:
                print("✅ PASS: Message formatted correctly")

    sock.close()


if __name__ == "__main__":
    main()
