#!/usr/bin/env python3
"""Simple combat test with existing character"""
import socket
import time

def send(sock, msg):
    sock.sendall((msg + "\n").encode())
    time.sleep(0.3)

def recv_all(sock, timeout=1.0):
    sock.settimeout(timeout)
    data = b""
    try:
        while True:
            chunk = sock.recv(4096)
            if not chunk:
                break
            data += chunk
    except socket.timeout:
        pass
    return data.decode(errors='ignore')

def main():
    print("=== Testing Combat with Vex ===\n")
    
    # Connect
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect(("localhost", 4000))
    print("Connected")
    
    # Welcome
    recv_all(sock, 0.5)
    
    # Login as Vex (existing character)
    send(sock, "Vex")
    time.sleep(0.5)
    resp = recv_all(sock, 1.0)
    print(f"\nLogged in: {resp[-200:]}")
    
    # Give Vex the Slash skill manually via spellbook check
    print("\n--- Checking spellbook ---")
    send(sock, "spellbook")
    time.sleep(0.3)
    resp = recv_all(sock)
    print(f"Spellbook: {resp}")
    
    # Teleport to goblin caves
    print("\n--- Teleporting to goblin cave ---")
    send(sock, "teleport 20000")
    time.sleep(0.5)
    resp = recv_all(sock)
    print(f"Teleport: {resp[-300:]}")
    
    # Look
    print("\n--- Looking ---")
    send(sock, "look")
    time.sleep(0.3)
    resp = recv_all(sock)
    print(f"Room: {resp}")
    
    # Try slash (will fail if not learned)
    print("\n--- Attempting slash ---")
    send(sock, "slash sentry")
    time.sleep(0.3)
    resp = recv_all(sock)
    print(f"Slash result: {resp}")
    
    # Quit
    send(sock, "quit")
    sock.close()
    print("\n=== Done ===")

if __name__ == "__main__":
    main()
