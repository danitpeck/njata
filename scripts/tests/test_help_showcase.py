#!/usr/bin/env python3
"""Test help command with learning new spell"""
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

print("=" * 60)
print("HELP COMMAND FEATURE SHOWCASE")
print("=" * 60)

# Restore to clear any studied spells
print("\n[RESTORE] Resetting to clean state...")
send_cmd("restore")

# Study a new spell (Shadow Veil - spell 1004)
print("\n[STUDY] Learning Shadow Veil from wand of blindness...")
send_cmd("teleport 9001")
time.sleep(0.3)
recv_all()
send_cmd("study wand of blindness")
time.sleep(0.5)
study_resp = recv_all()
if "Shadow Veil" in study_resp or "1004" in study_resp or "already know" in study_resp:
    print("✓ Studied Shadow Veil (or already knew it)")
else:
    print("Study response:", study_resp[:100])

# Now get help on Shadow Veil
print("\n[HELP] Getting detailed info on Shadow Veil...")
send_cmd("help shadow")
time.sleep(0.3)
help_resp = recv_all()
print(help_resp)

# Show what happens with ambiguous search
print("\n[AMBIGUOUS] Trying to find 'spell' (ambiguous)...")
send_cmd("help spell")
time.sleep(0.3)
resp = recv_all()
print(resp)

# Show not found
print("\n[NOT FOUND] Trying to find 'fireball' (not in game)...")  
send_cmd("help fireball")
time.sleep(0.3)
resp = recv_all()
print(resp)

sock.close()
print("=" * 60)
print("✓ Help system demo complete!")
print("=" * 60)
