#!/usr/bin/env python3
"""Test new abilities command"""
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
recv()

print("=== Testing 'abilities' command ===\n")
send("abilities")
resp = recv()
print(resp)

print("\n=== Testing 'score' command (should not exist) ===\n")
send("score")
resp = recv()
print(resp)

print("\n=== Testing 'stats' command (should still work) ===\n")
send("stats")
resp = recv()
print(resp)

send("quit")
sock.close()
