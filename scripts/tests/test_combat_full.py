#!/usr/bin/env python3
"""Full combat test - slash until mob dies"""
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
    print("=== Full Combat Test ===\n")
    
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect(("localhost", 4000))
    
    recv_all(sock, 0.5)
    send(sock, "Vex")
    time.sleep(0.5)
    recv_all(sock, 1.0)
    
    # Teleport to goblin cave
    send(sock, "teleport 20000")
    time.sleep(0.5)
    resp = recv_all(sock)
    print(f"Teleported: {resp[-200:]}")
    
    # Look at mob
    send(sock, "look")
    time.sleep(0.3)
    resp = recv_all(sock)
    print(f"Room: {resp}")
    
    # Slash until mob dies
    print("\n--- Slashing sentry until death ---")
    for i in range(15):
        time.sleep(2.2)  # Cooldown is 2s
        send(sock, "slash sentry")
        time.sleep(0.4)
        resp = recv_all(sock)
        
        # Print only the relevant lines
        lines = resp.strip().split('\n')
        for line in lines:
            if 'slash' in line.lower() or 'damage' in line.lower() or 'HP remaining' in line.lower() or 'defeated' in line.lower() or "don't see" in line.lower() or 'Proficiency' in line:
                print(f"  {line.strip()}")
        
        if 'defeated' in resp.lower() or "don't see" in resp.lower():
            print(f"\n✓ Combat ended after {i+1} slashes")
            break
    
    # Check final spellbook
    print("\n--- Final spellbook ---")
    send(sock, "spellbook")
    time.sleep(0.3)
    resp = recv_all(sock)
    for line in resp.split('\n'):
        if '9002' in line or 'Slash' in line or 'Proficiency' in line:
            print(f"  {line.strip()}")
    
    # Look at room (mob should be gone)
    print("\n--- Final room check ---")
    send(sock, "look")
    time.sleep(0.3)
    resp = recv_all(sock)
    print(f"Room after combat: {resp}")
    
    send(sock, "quit")
    sock.close()
    print("\n=== Test Complete ===")

if __name__ == "__main__":
    main()
