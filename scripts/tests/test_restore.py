#!/usr/bin/env python3
"""Test the keeper restore command"""
import socket
import time


def send_command(sock, command):
    """Send a command and get the response."""
    sock.sendall((command + "\n").encode('utf-8'))
    time.sleep(0.2)
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

    print("Testing restore command...")
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect((host, port))

    # Login as Vex (keeper)
    send_command(sock, "vex")
    send_command(sock, "test")

    # Check initial stats
    print("\n=== Initial Stats ===")
    response = send_command(sock, "stats")
    for line in response.split('\n')[-10:]:
        if line.strip():
            print(line)

    # Teleport and damage herself
    print("\n=== Teleporting to test arena ===")
    send_command(sock, "teleport 90000")
    
    print("\n=== Spawning mob ===")
    send_command(sock, "spawn 90001")
    
    print("\n=== Attacking mob (taking damage) ===")
    response = send_command(sock, "slash dummy")
    for line in response.split('\n'):
        if line.strip() and ('strike' in line.lower() or 'slash' in line.lower()):
            print(line)
    
    # Check damaged stats
    print("\n=== Damaged Stats ===")
    response = send_command(sock, "stats")
    for line in response.split('\n')[-10:]:
        if line.strip():
            print(line)

    # Restore command
    print("\n=== Using RESTORE command ===")
    response = send_command(sock, "restore")
    print(response.strip())

    # Check restored stats
    print("\n=== Restored Stats ===")
    response = send_command(sock, "stats")
    for line in response.split('\n')[-10:]:
        if line.strip():
            print(line)

    sock.close()
    print("\n=== Done ===")


if __name__ == "__main__":
    main()
