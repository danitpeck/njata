#!/usr/bin/env python3
"""Test spell targeting in combat"""
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

# Teleport to test arena
send("teleport 90000")
time.sleep(0.3)
resp = recv()
print("Teleported to test arena")
print(resp[-200:])

# Spawn a test mob
print("\n=== Spawning test mob (vnum 90001) ===")
send("spawn 90001")
time.sleep(0.3)
resp = recv()
print(resp.strip())

# Test different casting syntaxes
print("\n=== Test 1: Cast without target (should fail) ===")
send("cast arcane bolt")
time.sleep(0.3)
resp = recv()
print(resp)

print("\n=== Test 2: Cast with target keyword ===")
send("cast arcane bolt dummy")
time.sleep(0.3)
resp = recv()
print(resp)

print("\n=== Test 3: Cast Leviathan's Fire (multi-word spell) ===")
time.sleep(5.5)  # Wait for cooldown
send("spawn 90001")  # Spawn another target
time.sleep(0.3)
recv()
send("cast leviathan's fire dummy")
time.sleep(0.3)
resp = recv()
print(resp)

print("\n=== Test 4: Slash for comparison ===")
time.sleep(2.5)
send("spawn 90001")  # Spawn another target
time.sleep(0.3)
recv()
send("slash dummy")
time.sleep(0.3)
resp = recv()
print(resp)

print("\n=== Test 5: Check abilities ===")
send("abilities")
time.sleep(0.3)
resp = recv()
print(resp)

send("quit")
sock.close()
print("\n=== Done ===")
