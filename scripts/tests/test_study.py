#!/usr/bin/env python3
"""Test the Study command with actual magical items"""
import socket
import time


def send_command(sock, command):
    sock.sendall((command + "\n").encode('utf-8'))
    time.sleep(0.3)
    response = b''
    while True:
        try:
            sock.settimeout(0.1)
            chunk = sock.recv(4096)
            if not chunk:
                break
            response += chunk
        except socket.timeout:
            break
    return response.decode('utf-8', errors='ignore')


sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.connect(('localhost', 4000))

# Login
send_command(sock, "vex")
send_command(sock, "test")

print("=== Teleporting to Library (room 9001) ===")
response = send_command(sock, "teleport 9001")
for line in response.split('\n')[-15:]:
    if line.strip():
        print(line)

print("\n=== Looking at room (should see items) ===")
response = send_command(sock, "look")
for line in response.split('\n')[-20:]:
    if line.strip():
        print(line)

print("\n=== Checking abilities BEFORE study ===")
response = send_command(sock, "abilities")
lines = response.split('\n')
for i, line in enumerate(lines):
    if 'SPELLS' in line or 'MANEUVERS' in line:
        for j in range(i, min(i+10, len(lines))):
            if lines[j].strip():
                print(lines[j])
        break

print("\n=== Studying a wand (wand of magic missile, spell 1001) ===")
response = send_command(sock, "study wand")
for line in response.split('\n'):
    if line.strip():
        print(line)

print("\n=== Checking abilities AFTER study ===")
response = send_command(sock, "abilities")
lines = response.split('\n')
for i, line in enumerate(lines):
    if 'SPELLS' in line:
        for j in range(i, min(i+15, len(lines))):
            if lines[j].strip():
                print(lines[j])
        break

print("\n=== Looking again (wand should be gone) ===")
response = send_command(sock, "look")
for line in response.split('\n')[-15:]:
    if line.strip() and ('wand' in line.lower() or 'item' in line.lower() or 'scroll' in line.lower()):
        print(line)

print("\n=== Trying to study the same item again (should fail) ===")
response = send_command(sock, "study wand")
for line in response.split('\n'):
    if line.strip():
        print(line)
        break

sock.close()
print("\n=== Done ===")
