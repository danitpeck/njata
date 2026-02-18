#!/usr/bin/env python3
"""
Test admin commands: teleport, setadmin, demote
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

# Get banner and login as admin
time.sleep(0.5)
sock.recv(4096)
sock.send(b"Vex\r\n")
time.sleep(1.5)
response = sock.recv(4096).decode('utf-8', errors='ignore')

# Clear input buffer
sock.settimeout(0.1)
try:
    while True:
        sock.recv(1024)
except socket.timeout:
    pass
sock.settimeout(None)

print("=" * 60)
print("KEEPER COMMANDS TEST")
print("=" * 60)

# Test 1: Teleport
print("\n[TEST 1] Teleport to spell items room (9001)")
print("-" * 60)
response = send_command(sock, "teleport 9001")
print(response)

# Test 2: Check location
print("\n[TEST 2] Look around (should be in spell items room)")
print("-" * 60)
response = send_command(sock, "look")
print(response[:300])

# Test 3: Try study command in spell items room
print("\n[TEST 3] Study fireball wand")
print("-" * 60)
response = send_command(sock, "study fireball")
print(response)

# Test 4: Check spellbook
print("\n[TEST 4] Spellbook")
print("-" * 60)
response = send_command(sock, "spellbook")
print(response[:300])

# Test 5: Cast fireball
print("\n[TEST 5] Cast fireball")
print("-" * 60)
response = send_command(sock, "cast fireball")
print(response)

sock.close()
print("\n" + "=" * 60)
print("Admin Commands Test Complete")
print("=" * 60)
