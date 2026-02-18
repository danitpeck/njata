#!/usr/bin/env python3
"""Test capitalized player names"""
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

# Skip welcome
recv_all()

# Login
print("[TEST] Capitalized player names everywhere\n")
send_cmd("vex")
time.sleep(0.3)

# Test: Who command
print("[WHO] Players online list:")
who_resp = send_cmd("who")
print(who_resp)
if "Vex" in who_resp:
    print("✓ WHO shows 'Vex'\n")
else:
    print("✗ WHO doesn't show capitalized name\n")

# Test: Look (also here list)
print("[LOOK] Looking around in Library:")
send_cmd("teleport 9001")
time.sleep(0.3)
recv_all()
look_resp = send_cmd("look")
print(look_resp)
if "Vex" in look_resp:
    print("✓ LOOK shows 'Vex' in also here list\n")

# Test: Say (should show capitalized in room broadcast)
print("\n[SAY] Testing say message broadcast:")
say_resp = send_cmd("say Hello everyone!")
print(say_resp)
print("(If another player saw this, it would show 'Vex says...')")

sock.close()
print("\n✓ Capitalization test complete!")
