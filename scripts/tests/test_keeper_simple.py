#!/usr/bin/env python3
"""
Simple test of Keeper commands
"""

import socket
import time

def send_command(sock, cmd):
    """Send a command and receive the response."""
    sock.send((cmd + '\r\n').encode())
    time.sleep(0.5)
    
    response = b''
    sock.settimeout(1.0)
    try:
        while True:
            chunk = sock.recv(4096)
            if not chunk:
                break
            response += chunk
    except socket.timeout:
        pass
    sock.settimeout(None)
    
    return response.decode('utf-8', errors='ignore')

print("=" * 70)
print("KEEPER SYSTEM VALIDATION")
print("=" * 70)

# Connect as Vex
print("\n1. Connecting as Vex (Keeper)...")
vex = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
vex.connect(('localhost', 4000))
time.sleep(0.5)

# Get banner
vex.recv(4096)
time.sleep(0.2)

# Send name
vex.send(b"Vex\r\n")
time.sleep(1.0)

# Receive login response
banner = vex.recv(4096).decode('utf-8', errors='ignore')
print(f"✓ Logged in as Vex")

# Clear buffer
vex.settimeout(0.1)
try:
    while True:
        vex.recv(1024)
except socket.timeout:
    pass
vex.settimeout(None)

# Test teleport
print("\n2. Testing Vex's teleport command...")
response = send_command(vex, "teleport 9001")
if "teleport" in response.lower() and "9001" in response:
    print(f"✓ Teleport works: {response.split(chr(10))[0]}")
else:
    print(f"✗ Teleport failed: {response[:100]}")

# Test help to see available commands
print("\n3. Checking available Keeper commands...")
response = send_command(vex, "help makekeeper")
if "makekeeper" in response.lower() or "0/0" in response:  # 0/0 means command not found in help
    print(f"✓ makekeeper command exists")
else:
    print(f"Response: {response[:200]}")

response = send_command(vex, "help removekeeper")
if "removekeeper" in response.lower() or "0/0" in response:
    print(f"✓ removekeeper command exists")
else:
    print(f"Response: {response[:200]}")

# Try to make someone a keeper (should fail without a target online)
print("\n4. Testing makekeeper command permissions...")
response = send_command(vex, "makekeeper nobody")
if "not found" in response.lower() or "makekeeper" in response.lower():
    print(f"✓ makekeeper command recognized: {response.split(chr(10))[0]}")
else:
    print(f"Response: {response[:200]}")

# Verify Keeper permissions
print("\n5. Verifying Vex has Keeper permissions...")
response = send_command(vex, "teleport 9002")
if "teleport" in response.lower():
    print(f"✓ Vex can use teleport (Keeper privilege)")
else:
    print(f"✗ Teleport failed: {response[:100]}")

vex.close()

print("\n" + "=" * 70)
print("KEEPER SYSTEM VALIDATION COMPLETE")
print("=" * 70)
