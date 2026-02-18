import socket
import time
import sys

def send_command(sock, cmd):
    try:
        sock.sendall((cmd + '\n').encode())
        time.sleep(0.1)
        return sock.recv(4096).decode()
    except Exception as e:
        print(f"Error: {e}")
        return ""

try:
    # Connect to server
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect(('localhost', 4000))
    print("Connected to server")
    
    # Read banner
    banner = sock.recv(1024).decode()
    print("Received banner:", banner[:50])
    time.sleep(0.1)
    
    # Login
    response = send_command(sock, 'Vex')
    print("Login response length:", len(response))
    
    # Send save command
    time.sleep(0.2)
    response = send_command(sock, 'save')
    print("Save command response:", response[:100] if response else "No response")
    
    # Check if player file was written
    import os
    if os.path.exists('players/Vex.json'):
        print("[OK] Player save file created")
        with open('players/Vex.json', 'r') as f:
            data = f.read()
            if 'Vex' in data and 'is_keeper' in data:
                print("[OK] Player file contains correct data")
    else:
        print("[FAIL] Player save file not found")
    
    sock.close()
except Exception as e:
    print(f"Fatal error: {e}")
    import traceback
    traceback.print_exc()

