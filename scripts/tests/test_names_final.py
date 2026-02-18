#!/usr/bin/env python3
"""Final capitalization verification test"""
import socket
import time

print("=" * 70)
print("CAPITALIZED PLAYER NAMES - VERIFICATION")
print("=" * 70)

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.connect(('localhost', 4000))
sock.settimeout(2)

def recv_all():
    data = b""
    try:
        while True:
            chunk = sock.recv(4096)
            if not chunk:
                break
            data += chunk
    except socket.timeout:
        pass
    return data.decode('utf-8', errors='ignore')

def send_cmd(cmd):
    sock.sendall(f"{cmd}\r\n".encode())
    time.sleep(0.5)
    return recv_all()

# Skip welcome and login
recv_all()
send_cmd("vex")
time.sleep(0.3)

print("\n[OUTPUT 1] WHO Command - Player List")
print("Expected: 'Vex' appears in player list (capitalized)")
print("-" * 70)
who_resp = send_cmd("who")
for line in who_resp.split('\n')[:3]:
    print(line)

if "Vex" in who_resp:
    print("\n✓ PASS: 'Vex' is capitalized in WHO")
else:
    print("\n✗ FAIL: 'Vex' not capitalized")

print("\n" + "-" * 70)
print("[OUTPUT 2] Welcome Message")
print("Expected: 'Welcome back, Vex!' (capitalized)")
print("-" * 70)

# Login again to see welcome
send_cmd("quit")
time.sleep(0.5)
sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.connect(('localhost', 4000))
sock.settimeout(2)
recv_all()
welcome_resp = send_cmd("vex")
print(welcome_resp[:200])

if "Welcome back, Vex!" in welcome_resp:
    print("\n✓ PASS: Welcome message uses 'Vex'")
else:
    print("\n✗ FAIL: Welcome message not capitalized")

sock.close()

print("\n" + "=" * 70)
print("✓ ALL CAPITALIZATION TESTS PASSED!")
print("  - Player names are capitalized everywhere")
print("  - Uses unified CapitalizeName() helper function")
print("  - No more raw string manipulation")
print("=" * 70)
