#!/usr/bin/env python3
"""Quick combat test"""
import socket
import time

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.connect(("localhost", 4000))
sock.settimeout(2.0)

def recv():
    try:
        return sock.recv(8192).decode(errors='ignore')
    except:
        return ""

def send(msg):
    sock.sendall((msg + "\n").encode())
    time.sleep(0.5)

# Login
recv()
send("Vex")
print(recv())

# Teleport
send("teleport 20000")
print(recv())

# Slash once
send("slash sentry")
resp = recv()
print("\nFirst slash:")
print(resp)

# Wait and slash again
time.sleep(2.5)
send("slash sentry")
resp = recv()
print("\nSecond slash:")
print(resp)

# Check spellbook
send("spellbook")
resp = recv()
print("\nSpellbook:")
for line in resp.split('\n'):
    if 'Slash' in line or '9002' in line or 'Proficiency' in line:
        print(line)

send("quit")
sock.close()
print("\nDone!")
