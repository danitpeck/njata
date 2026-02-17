#!/usr/bin/env python3
"""Test combat system with Warrior kit and Slash maneuver"""
import socket
import time

def send(sock, msg):
    sock.sendall((msg + "\n").encode())
    time.sleep(0.2)

def recv_all(sock, timeout=1.0):
    sock.settimeout(timeout)
    data = b""
    try:
        while True:
            chunk = sock.recv(4096)
            if not chunk:
                break
            data += chunk
    except socket.timeout:
        pass
    return data.decode(errors='ignore')

def main():
    print("=== Testing Combat System ===\n")
    
    # Connect
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect(("localhost", 4000))
    print("Connected to server")
    
    # Receive welcome
    welcome = recv_all(sock, 0.5)
    print(f"Server says: {welcome[:100]}...")
    
    # Create character with Warrior's Kit
    print("\n--- Creating Warrior character ---")
    send(sock, "TestWarrior")
    time.sleep(0.5)
    recv_all(sock, 0.5)
    
    # Race selection (1 = Human)
    send(sock, "1")
    time.sleep(0.5)
    recv_all(sock, 0.5)
    
    # Sex selection (1 = Male)
    send(sock, "1")
    time.sleep(0.5)
    recv_all(sock, 0.5)
    
    # Confirm sex
    send(sock, "yes")
    time.sleep(0.5)
    recv_all(sock, 0.5)
    
    # Age selection (2 = Youth)
    send(sock, "2")
    time.sleep(0.5)
    recv_all(sock, 0.5)
    
    # Kit selection (2 = Warrior's Kit - has Slash at 10%)
    send(sock, "2")
    time.sleep(0.5)
    resp = recv_all(sock, 0.5)
    print(f"Kit selection response: {resp[-200:]}")
    
    # Confirm kit
    send(sock, "yes")
    time.sleep(2.0)  # Give time for character creation to complete
    resp = recv_all(sock, 2.0)
    print(f"Character created: {resp[-300:]}")
    
    # Check stats
    print("\n--- Checking stats ---")
    send(sock, "stats")
    time.sleep(0.3)
    resp = recv_all(sock)
    print(f"Stats: {resp}")
    
    # Check spellbook (should show Slash maneuver)
    print("\n--- Checking spellbook ---")
    send(sock, "spellbook")
    time.sleep(0.3)
    resp = recv_all(sock)
    print(f"Spellbook: {resp}")
    
    # Teleport to goblin cave (20000 has a sentry mob)
    print("\n--- Teleporting to goblin cave ---")
    send(sock, "makekeeper TestWarrior")
    time.sleep(0.3)
    recv_all(sock)
    
    send(sock, "teleport 20000")
    time.sleep(0.3)
    resp = recv_all(sock)
    print(f"Teleport response: {resp}")
    
    # Look at the room
    print("\n--- Looking at room ---")
    send(sock, "look")
    time.sleep(0.3)
    resp = recv_all(sock)
    print(f"Room: {resp}")
    
    # Try to slash the sentry
    print("\n--- Attempting to slash sentry ---")
    send(sock, "slash sentry")
    time.sleep(0.5)
    resp = recv_all(sock)
    print(f"Slash result: {resp}")
    
    # Slash again to see cooldown
    print("\n--- Attempting second slash (testing cooldown) ---")
    send(sock, "slash sentry")
    time.sleep(0.3)
    resp = recv_all(sock)
    print(f"Second slash: {resp}")
    
    # Wait for cooldown (2 seconds)
    print("\n--- Waiting 2 seconds for cooldown ---")
    time.sleep(2.5)
    
    # Slash again
    print("\n--- Third slash after cooldown ---")
    send(sock, "slash sentry")
    time.sleep(0.3)
    resp = recv_all(sock)
    print(f"Third slash: {resp}")
    
    # Keep slashing until mob dies
    print("\n--- Slashing until mob dies ---")
    for i in range(10):
        time.sleep(2.2)  # Wait for cooldown
        send(sock, "slash sentry")
        time.sleep(0.3)
        resp = recv_all(sock)
        print(f"Slash {i+4}: {resp[-200:]}")
        if "defeated" in resp.lower() or "don't see" in resp.lower():
            print("Mob defeated or not found!")
            break
    
    # Look at room again
    print("\n--- Final room check ---")
    send(sock, "look")
    time.sleep(0.3)
    resp = recv_all(sock)
    print(f"Final room: {resp}")
    
    # Check spellbook for proficiency increase
    print("\n--- Final spellbook check ---")
    send(sock, "spellbook")
    time.sleep(0.3)
    resp = recv_all(sock)
    print(f"Final spellbook: {resp}")
    
    # Quit
    send(sock, "quit")
    sock.close()
    print("\n=== Test Complete ===")

if __name__ == "__main__":
    main()
