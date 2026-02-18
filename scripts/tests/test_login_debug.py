#!/usr/bin/env python3
"""Debug login issues"""
import socket
import time

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
print("[1] Connecting to localhost:4000...")
try:
    sock.connect(('localhost', 4000))
    print("✓ Connected")
except Exception as e:
    print(f"✗ Connection failed: {e}")
    exit(1)

print("\n[2] Reading welcome message...")
time.sleep(0.5)
sock.settimeout(2)
try:
    welcome = sock.recv(4096).decode('utf-8', errors='ignore')
    print(f"Welcome message:\n{welcome}")
except socket.timeout:
    print("✗ Timeout reading welcome")
except Exception as e:
    print(f"✗ Error: {e}")

print("\n[3] Sending username 'vex'...")
sock.sendall(b"vex\r\n")
time.sleep(0.5)
try:
    resp = sock.recv(4096).decode('utf-8', errors='ignore')
    print(f"Response:\n{resp}")
except Exception as e:
    print(f"✗ Error: {e}")

print("\n[4] Sending password 'password'...")
sock.sendall(b"password\r\n")
time.sleep(1)
try:
    resp = sock.recv(4096).decode('utf-8', errors='ignore')
    print(f"Response:\n{resp}")
except socket.timeout:
    print("✗ Timeout - server may have crashed")
except Exception as e:
    print(f"✗ Error: {e}")

print("\n[5] Trying to send command after login...")
sock.sendall(b"look\r\n")
time.sleep(0.5)
try:
    resp = sock.recv(4096).decode('utf-8', errors='ignore')
    print(f"Response:\n{resp}")
except socket.timeout:
    print("✗ Timeout - server may have crashed")
except Exception as e:
    print(f"✗ Error: {e}")

sock.close()
print("\n✓ Test complete")
