#!/usr/bin/env python3
"""
Test Keeper commands: makekeeper, removekeeper, teleport
Tests the refactored admin system with custodial terminology
"""

import socket
import time
import json

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

def connect_player(player_name):
    """Create a socket connection and login as a player."""
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect(('localhost', 4000))
    
    # Get banner
    time.sleep(0.5)
    sock.recv(4096)
    
    # Send player name
    sock.send((player_name + '\r\n').encode())
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
    
    return sock

print("=" * 70)
print("COMPREHENSIVE KEEPER SYSTEM TEST")
print("=" * 70)

# Test 1: Login as Vex (Keeper) and test teleport
print("\n[TEST 1] Login as Vex (Keeper) and test teleport command")
print("-" * 70)
vex = connect_player("Vex")
response = send_command(vex, "teleport 9001")
if "teleport" in response.lower():
    print("✓ Teleport command works for Keeper")
    print(f"  Response: {response.split(chr(10))[0]}")
else:
    print("✗ Teleport failed")

# Test 2: Vex elevates Dani to Keeper status
print("\n[TEST 2] Vex elevates Dani to Keeper (makekeeper command)")
print("-" * 70)
response = send_command(vex, "makekeeper Dani")
print(f"Response: {response}")
if "responsibility of a Keeper" in response:
    print("✓ Correct messaging: 'responsibility of a Keeper'")
else:
    print("✗ Wrong messaging received")

# Test 3: Check Dani's data file was updated
print("\n[TEST 3] Verify Dani's player file updated with is_keeper flag")
print("-" * 70)
try:
    with open('players/dani.json', 'r') as f:
        data = json.load(f)
    if 'is_keeper' in data and data['is_keeper']:
        print(f"✓ Dani.json has is_keeper: {data['is_keeper']}")
    else:
        print(f"✗ Dani.json missing is_keeper or not true: {data.get('is_keeper', 'MISSING')}")
except Exception as e:
    print(f"✗ Error reading Dani.json: {e}")

# Test 4: Login as Dani and verify Keeper status persists
print("\n[TEST 4] Login as Dani and verify Keeper status persists")
print("-" * 70)
dani = connect_player("Dani")
response = send_command(dani, "teleport 9002")
if "teleport" in response.lower():
    print("✓ Dani (newly promoted Keeper) can use teleport")
else:
    print("✗ Dani cannot use teleport (Keeper status didn't persist)")

# Test 5: Vex removes Dani's Keeper status
print("\n[TEST 5] Vex removes Dani's Keeper status (removekeeper command)")
print("-" * 70)
response = send_command(vex, "removekeeper Dani")
print(f"Response: {response}")
if "Keeper responsibilities" in response or "stripped" in response.lower():
    print("✓ Correct messaging related to Keeper responsibilities")
else:
    print("✗ Wrong messaging received")

# Test 6: Verify Dani lost Keeper status
print("\n[TEST 6] Verify Dani lost Keeper status")
print("-" * 70)
try:
    with open('players/dani.json', 'r') as f:
        data = json.load(f)
    if 'is_keeper' in data and not data['is_keeper']:
        print(f"✓ Dani.json is_keeper set to false")
    elif 'is_keeper' not in data:
        print(f"✓ Dani.json has is_keeper field (if not present, defaults to false)")
    else:
        print(f"✗ Dani.json is_keeper still true: {data['is_keeper']}")
except Exception as e:
    print(f"✗ Error reading Dani.json: {e}")

# Test 7: Verify Dani can no longer use teleport
print("\n[TEST 7] Dani attempts teleport and should fail")
print("-" * 70)
dani2 = connect_player("Dani")
response = send_command(dani2, "teleport 9003")
if "authority" in response.lower() or "permission" in response.lower():
    print(f"✓ Dani (non-Keeper) denied access with authority message")
    print(f"  Response: {response.split(chr(10))[0]}")
else:
    print(f"✗ Expected permission denial, got: {response[:100]}")

# Test 8: Verify non-Keeper cannot use makekeeper
print("\n[TEST 8] Dani attempts makekeeper and should fail")
print("-" * 70)
response = send_command(dani2, "makekeeper Vex")
if "authority" in response.lower() or "permission" in response.lower():
    print(f"✓ Non-Keeper denied access to makekeeper with authority message")
else:
    print(f"✗ Non-Keeper should not have access to makekeeper")

# Cleanup
vex.close()
dani.close()
dani2.close()

print("\n" + "=" * 70)
print("KEEPER SYSTEM TEST COMPLETE")
print("=" * 70)
