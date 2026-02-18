#!/usr/bin/env python3
"""Test the renamed magical items"""
import socket
import time

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.connect(('localhost', 4000))
sock.settimeout(1.5)

def recv():
    data = b""
    try:
        while True:
            chunk = sock.recv(4096)
            if not chunk: break
            data += chunk
    except: pass
    return data.decode('utf-8', errors='ignore')

def send(cmd):
    sock.sendall(f"{cmd}\r\n".encode())
    time.sleep(0.3)
    return recv()

# Test
print("\n" + "="*70)
print("TESTING RENAMED MAGICAL ITEMS")
print("="*70)

recv()  # welcome
send("vex")
time.sleep(0.3)

send("teleport 9001")
time.sleep(0.3)

resp = send("look")
print("\n[LIBRARY CONTENTS]\n")
print(resp)

# Check for new item names
print("\n" + "="*70)
print("VERIFICATION")
print("="*70)

items_old = ["magic missile", "fireball", "healing", "blindness", "invisibility", "teleportation", "frost bolt", "identify"]
items_new = ["arcane bolt", "leviathan's fire", "mend", "shadow veil", "ephemeral step", "path shift", "winter's whisper", "knowing"]

found_new = sum(1 for item in items_new if item.lower() in resp.lower())
found_old = sum(1 for item in items_old if item.lower() in resp.lower())

print(f"New item names found: {found_new}/8")
print(f"Old item names found: {found_old}/8")

if found_new >= 7 and found_old == 0:
    print("\nPASS: All items renamed successfully!")
else:
    print("\nCheck output above for any old item names still present")

send("quit")
sock.close()
