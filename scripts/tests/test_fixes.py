#!/usr/bin/env python3
"""Quick test for capitalization and quit fixes"""
import socket
import time

class QuickTest:
    def __init__(self):
        self.sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self.sock.connect(('localhost', 4000))
        self.sock.settimeout(1.5)
        
    def recv_safe(self):
        data = b""
        try:
            while True:
                try:
                    chunk = self.sock.recv(2048)
                    if not chunk:
                        break
                    data += chunk
                except socket.timeout:
                    break
        except:
            pass
        return data.decode('utf-8', errors='ignore')
    
    def send_cmd(self, cmd, wait=0.4):
        self.sock.sendall(f"{cmd}\r\n".encode())
        time.sleep(wait)
        return self.recv_safe()
    
    def run(self):
        print("\n" + "="*70)
        print("🐛 TESTING FIXES")
        print("="*70)
        
        # Welcome
        self.recv_safe()
        self.send_cmd("vex")
        
        print("\n[TEST 1] Combat message capitalization")
        self.send_cmd("teleport 90000")
        self.send_cmd("restore")
        self.send_cmd("spawn 90001")
        
        resp = self.send_cmd("slash dummy")
        print(f"Response:\n{resp}")
        
        if "Vex slashes" in resp:
            print("✅ PASS: Combat message has capitalized 'Vex'")
        elif "vex slashes" in resp:
            print("❌ FAIL: Combat message still lowercase")
        else:
            print("⚠️  Check response above")
        
        print("\n[TEST 2] Quit command")
        resp = self.send_cmd("quit")
        print(f"Response: {resp.strip()}")
        
        if "Goodbye" in resp:
            print("✅ PASS: Quit command works and responds")
        else:
            print("❌ FAIL: Quit command missing or not recognized")
        
        try:
            self.sock.close()
        except:
            pass

if __name__ == "__main__":
    test = QuickTest()
    test.run()
