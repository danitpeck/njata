#!/usr/bin/env python3
"""
Simple cast command tester for NJATA MVP spells system.
Tests the new spell casting and studying mechanics.
"""

import socket
import time
import sys

def send_command(sock, cmd):
    """Send a command and receive the response."""
    sock.send((cmd + '\r\n').encode())
    time.sleep(0.2)  # Give server time to process
    
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

def test_spells():
    """Test cast and study commands."""
    print("=== NJATA MVP SPELL SYSTEM TESTER ===\n")
    
    # Connect to server
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect(('localhost', 4000))
    
    # Skip banner
    time.sleep(0.5)
    sock.recv(4096)
    
    # Login with existing character
    print("[TEST] Logging in as SpellTest...")
    sock.send(b"SpellTest\r\n")
    time.sleep(0.5)
    response = sock.recv(4096).decode('utf-8', errors='ignore')
    print(f"Response: {response[:100]}...")
    
    # If new character, skip creation (just press enter for defaults)
    if "SELECT YOUR RACE" in response:
        print("[TEST] Creating new character (selecting defaults)...")
        # Select human
        sock.send(b"1\r\n")
        time.sleep(0.5)
        sock.recv(4096)
        # Confirm race
        sock.send(b"y\r\n")
        time.sleep(0.5)
        sock.recv(4096)
        # Select Scholar's Kit (has spells)
        sock.send(b"1\r\n")
        time.sleep(0.5)
        sock.recv(4096)
        # Confirm kit
        sock.send(b"y\r\n")
        time.sleep(0.5)
        sock.recv(4096)
        # Age
        sock.send(b"Adult\r\n")
        time.sleep(0.5)
        sock.recv(4096)
        # Confirm age
        sock.send(b"y\r\n")
        time.sleep(0.5)
        sock.recv(4096)
        # Sex
        sock.send(b"Male\r\n")
        time.sleep(0.5)
        sock.recv(4096)
        # Confirm sex
        sock.send(b"y\r\n")
        time.sleep(0.5)
        sock.recv(4096)
        # Final confirm
        sock.send(b"y\r\n")
        time.sleep(0.5)
        sock.recv(4096)
    
    time.sleep(1)
    
    # Test 1: Check spellbook
    print("\n[TEST 1] Spellbook command...")
    response = send_command(sock, "spellbook")
    if "Magic Missile" in response:
        print("✓ Spellbook shows learned spells")
        print(f"  Found: Magic Missile\n  Response:\n{response[:200]}")
    else:
        print("❌ Spellbook didn't show expected spells")
        print(f"  Response: {response[:200]}")
    
    # Test 2: Cast Magic Missile
    print("\n[TEST 2] Cast Magic Missile...")
    response = send_command(sock, "cast magic missile")
    if "magical missile" in response.lower() or "you hurl" in response.lower() or "proficiency" in response:
        print("✓ Cast command executed")
        print(f"  Response: {response[:200]}")
    else:
        print("❌ Cast command failed")
        print(f"  Response: {response[:200]}")
    
    # Test 3: Try casting same spell again (should hit cooldown)
    print("\n[TEST 3] Cooldown check...")
    response = send_command(sock, "cast magic missile")
    if "cooldown" in response.lower() or "seconds remaining" in response.lower():
        print("✓ Cooldown detected")
        print(f"  Response: {response[:200]}")
    else:
        print("⚠️  Cooldown mechanism may not be working")
        print(f"  Response: {response[:200]}")
    
    # Test 4: Study command
    print("\n[TEST 4] Study command (looking for wand)...")
    response = send_command(sock, "study wand")
    if "learn" in response.lower() or "carefully study" in response.lower() or "spell" in response.lower():
        print("✓ Study command worked")
        print(f"  Response: {response[:200]}")
    else:
        print("❌ Study command didn't work as expected")
        print(f"  Response: {response[:200]}")
    
    # Test 5: Spellbook after study
    print("\n[TEST 5] Check spellbook after study...")
    response = send_command(sock, "spellbook")
    spell_count = response.count("[")
    if spell_count > 1:
        print(f"✓ Spellbook now shows {spell_count} spells")
        print(f"  Response: {response[:300]}")
    else:
        print("⚠️  Spellbook might not be showing new spell")
        print(f"  Response: {response[:300]}")
    
    # Cleanup
    sock.close()
    print("\n=== Test Complete ===\n")

if __name__ == '__main__':
    try:
        test_spells()
    except Exception as e:
        print(f"Error: {e}")
        sys.exit(1)
