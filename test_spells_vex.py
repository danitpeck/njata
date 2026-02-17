#!/usr/bin/env python3
"""
Spell test using an existing character that's already created.
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
print("=== Connected to NJATA ===\n")

# Login with existing character (Vex from test-client)
print("Logging in as 'Vex' (existing character)...")
sock.send(b"Vex\r\n")
time.sleep(1.5)  # Wait for login to complete
response = sock.recv(4096).decode('utf-8', errors='ignore')
print(f"Response:\n{response[:400]}\n")

# Clear input buffer
sock.settimeout(0.1)
try:
    while True:
        sock.recv(1024)
except socket.timeout:
    pass
sock.settimeout(None)

# Now test commands
print("=" * 50)
print("SPELL SYSTEM TESTS")
print("=" * 50)

print("\n[TEST 1] Spellbook")
print("-" * 50)
response = send_command(sock, "spellbook")
print(response)

print("\n[TEST 2] Cast magic missile")
print("-" * 50)
response = send_command(sock, "cast magic missile")
print(response)

print("\n[TEST 3] Cast magic missile again (should show cooldown)")
print("-" * 50)
response = send_command(sock, "cast magic missile")
print(response)

print("\n[TEST 4] Study wand (should learn a new spell)")
print("-" * 50)
response = send_command(sock, "study wand")
print(response)

print("\n[TEST 5] Spellbook (should show learned spell)")
print("-" * 50)
response = send_command(sock, "spellbook")
print(response)

print("\n[TEST 6] Cast learned spell")
print("-" * 50)
spell = "fireball"  # Study wand teaches one of the spells
response = send_command(sock, f"cast {spell}")
print(response)

sock.close()
print("\n" + "=" * 50)
print("Tests Complete")
print("=" * 50)
