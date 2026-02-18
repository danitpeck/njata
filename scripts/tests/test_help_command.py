#!/usr/bin/env python3
"""Test the help command"""
import socket
import time

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

# Login
recv_all()
send_cmd("vex")
time.sleep(0.5)

# Test help with no args
print("[TEST 1] help (no args)")
resp = send_cmd("help")
print(resp)

# Test help for a spell
print("\n[TEST 2] help arcane")
resp = send_cmd("help arcane")
print(resp)

# Test help for Slash
print("\n[TEST 3] help slash")
resp = send_cmd("help slash")
print(resp)

# Test help for Leviathan's Fire
print("\n[TEST 4] help fire")
resp = send_cmd("help fire")
print(resp)

# Test help for Mend
print("\n[TEST 5] help mend")
resp = send_cmd("help mend")
print(resp)

sock.close()
