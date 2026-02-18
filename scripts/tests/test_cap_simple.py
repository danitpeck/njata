#!/usr/bin/env python3
"""Direct combat test"""
import socket
import time

sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
sock.connect(('localhost', 4000))
sock.settimeout(1.0)

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
recv()  # welcome
send("vex")
send("teleport 90000")
send("restore")
send("spawn 90001")

# THE TEST
resp = send("slash dummy")
print(resp)

if "Vex slashes" in resp:
    print("\nPASS: CAPITALIZATION FIX WORKS!")
elif "Vex strikes" in resp:
    print("\nPASS: CAPITALIZATION FIX WORKS!")
else:
    print("\nCheck output above for 'vex' vs 'Vex'")

send("quit")
sock.close()
