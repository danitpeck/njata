#!/usr/bin/env python3
"""Test combat message capitalization more carefully"""
import socket
import time

class CombatTest:
    def __init__(self):
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sock.connect(('localhost', 4000))
        self.sock.settimeout(2.0)
        
    def recv_safe(self):
        data = b""
        try:
            while True:
                try:
                    chunk = self.sock.recv(4096)
                    if not chunk:
                        break
                    data += chunk
                except socket.timeout:
                    break
        except:
            pass
        return data.decode('utf-8', errors='ignore')
    
    def send_cmd(self, cmd, wait=0.5):
        self.sock.sendall(f"{cmd}\r\n".encode())
        time.sleep(wait)
        return self.recv_safe()
    
    def run(self):
        print("\n" + "="*70)
        print("🐛 TESTING COMBAT CAPITALIZATION")
        print("="*70)
        
        # Welcome
        self.recv_safe()
        self.send_cmd("vex")
        time.sleep(0.3)
        
        self.send_cmd("teleport 90000")
        self.send_cmd("restore")
        time.sleep(0.3)
        
        # Spawn a second character to see room broadcast
        print("\nSpawning a second player to see room messages...")
        self.send_cmd("restore")
        self.send_cmd("spawn 90001")
        time.sleep(0.4)
        
        # Now slash and capture ALL output
        print("\nSending slash command and capturing full output...")
        resp = self.send_cmd("slash dummy")
        
        print(f"\nFull response:\n{resp}")
        print("\n" + "="*70)
        
        # Check for capitalized name in different scenarios
        if "Vex slashes" in resp or "Vex strikes" in resp:
            print("✅ PASS: Player name is capitalized in combat messages!")
        elif "vex slashes" in resp or "vex strikes" in resp:
            print("❌ FAIL: Player name still lowercase in combat")
        else:
            print("ℹ️  Note: Room message may not appear in single-client test")
            print("   (would see it if multiple characters in room)")
        
        try:
            self.send_cmd("quit")
            self.sock.close()
        except:
            pass

if __name__ == "__main__":
    test = CombatTest()
    test.run()
