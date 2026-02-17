#!/usr/bin/env python3
"""
Spell system test that waits for character creation to complete.
"""

import socket
import time

def send_command(sock, cmd):
    """Send a command and receive the response."""
    sock.send((cmd + '\r\n').encode())
    time.sleep(0.3)
    
    response = b''
    sock.settimeout(0.5)
    try:
        while True:
            chunk = sock.recv(1024)
            if not chunk:
                break
            response += chunk
    except socket.timeout:
        pass
    sock.settimeout(None)
    
    return response.decode('utf-8', errors='ignore')

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.connect(('localhost', 4000))

# Get banner
time.sleep(0.5)
banner = sock.recv(4096).decode('utf-8', errors='ignore')
print("=== BANNER ===")
print(banner[:300])

# Login
print("\n=== LOGIN ===")
sock.send(b"SpellTester2\r\n")
time.sleep(0.5)
response = sock.recv(4096).decode('utf-8', errors='ignore')
print(f"Login response: {response[:200]}")

# Character creation - full race list check
if "SELECT YOUR RACE" in response:
    print("\n=== CHARACTER CREATION - Selecting Race 1 (Human) ===")
    sock.send(b"1\r\n")
    time.sleep(0.5)
    
    # Class selection
    response = sock.recv(4096).decode('utf-8', errors='ignore')
    print(f"Class response:\n{response[:200]}")
    
    if "SELECT YOUR CLASS" in response or "class" in response.lower():
        print("\n=== Selecting Class 1 (Warrior) ===")
        sock.send(b"1\r\n")
        time.sleep(0.5)
        
        # Stats/confirm
        response = sock.recv(4096).decode('utf-8', errors='ignore')
        print(f"Confirm response:\n{response[:200]}")
        
        # Confirm character
        sock.send(b"y\r\n")
        time.sleep(2)  # Wait for character to fully load
        
        # Clear received data
        try:
            while True:
                sock.recv(1024, socket.MSG_DONTWAIT)
        except:
            pass

# Now test commands
print("\n=== TESTING SPELL COMMANDS ===")

print("\n[TEST 1] Spellbook")
response = send_command(sock, "spellbook")
print(f"Response: {response[:300]}")
print(f"Has 'Magic Missile': {'Magic Missile' in response}")

print("\n[TEST 2] Cast magic missile")
response = send_command(sock, "cast magic missile")
print(f"Response: {response[:300]}")

print("\n[TEST 3] Study wand")
response = send_command(sock, "study wand")
print(f"Response: {response[:300]}")

print("\n[TEST 4] Spellbook (check for new spell)")
response = send_command(sock, "spellbook")
print(f"Response: {response[:400]}")

sock.close()
print("\n=== Complete ===")
