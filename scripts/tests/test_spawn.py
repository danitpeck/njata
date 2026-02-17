#!/usr/bin/env python3
"""Test the spawn command for keepers"""
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

    print("Testing spawn command...")
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect((host, port))

    # Login as Vex (keeper)
    send_command(sock, "vex")
    send_command(sock, "test")

    # Teleport to test arena
    print("\n=== Teleporting to test arena ===")
    response = send_command(sock, "teleport 90000")
    print(response.split('\n')[-5:])

    # Look to see current state of room
    print("\n=== Looking at room (before spawn) ===")
    response = send_command(sock, "look")
    lines = [line.strip() for line in response.split('\n') if line.strip()]
    for line in lines[-10:]:
        print(line)

    # Spawn a test dummy (vnum 90001)
    print("\n=== Test 1: Spawning mob vnum 90001 ===")
    response = send_command(sock, "spawn 90001")
    print(response.strip())

    # Look again to verify mob spawned
    print("\n=== Looking at room (after spawn) ===")
    response = send_command(sock, "look")
    lines = [line.strip() for line in response.split('\n') if line.strip()]
    for line in lines[-10:]:
        print(line)

    # Test invalid vnum
    print("\n=== Test 2: Spawning invalid vnum ===")
    response = send_command(sock, "spawn 99999")
    print(response.strip())

    # Test invalid syntax
    print("\n=== Test 3: Spawn without arguments ===")
    response = send_command(sock, "spawn")
    print(response.strip())

    # Test combat with spawned mob
    print("\n=== Test 4: Combat with spawned mob ===")
    response = send_command(sock, "slash dummy")
    print(response.strip())

    sock.close()
    print("\n=== Done ===")


if __name__ == "__main__":
    main()
