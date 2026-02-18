#!/usr/bin/env python3
"""Test Study command with learning a NEW spell"""
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

print("=== Teleporting to Library ===")
send_command(sock, "teleport 9001")

print("\n=== Current abilities ===")
response = send_command(sock, "abilities")
print(response[-300:])

print("\n=== Studying wand of blindness (spell 1004 - Shadow Veil) ===")
response = send_command(sock, "study blindness")
for line in response.split('\n'):
    if line.strip():
        print(line)

print("\n=== Updated abilities (should have Shadow Veil now) ===")
response = send_command(sock, "abilities")
print(response[-300:])

sock.close()
