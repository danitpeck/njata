#!/usr/bin/env python3
"""Test scattered magical items across 8 locations"""
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

print("\n" + "="*70)
print("TESTING SCATTERED ITEM LOCATIONS")
print("="*70)

recv()  # welcome
send("vex")
time.sleep(0.3)

locations = [
    (2200, "Wizard's Study", "Wand of Arcane Bolt"),
    (2201, "Docks of Immak", "Wand of Leviathan's Fire"),
    (2202, "Temple of Healing", "Scroll of Mend"),
    (2203, "Dark Grove", "Wand of Shadow Veil"),
    (2204, "Fairy Glade", "Feather of Ephemeral Step"),
    (2205, "Courier's Office", "Amulet of Path Shift"),
    (2206, "Mountain Peak", "Crystal of Winter's Whisper"),
    (2207, "Scholar's Archive", "Tome of Knowing"),
]

found_items = 0
for vnum, loc_name, item_name in locations:
    send(f"teleport {vnum}")
    resp = send("look")
    
    if item_name.lower() in resp.lower():
        print(f"[FOUND] {loc_name:20} - {item_name}")
        found_items += 1
    else:
        print(f"[MISS]  {loc_name:20} - NOT FOUND")

print("\n" + "="*70)
print(f"RESULT: Found {found_items}/8 items in scattered locations")
if found_items == 8:
    print("SUCCESS! Item distribution complete!")
else:
    print(f"Missing {8 - found_items} items")
print("="*70)

send("quit")
sock.close()
