#!/usr/bin/env python3
"""Test newbie area and new chat/help functionality."""
import socket
import time

def send_command(sock, command):
    """Send a command and get the response."""
    sock.sendall((command + "\n").encode())
    time.sleep(0.2)
    response = b""
    sock.settimeout(1)
    try:
        while True:
            data = sock.recv(1024)
            if not data:
                break
            response += data
    except socket.timeout:
        pass
    return response.decode('utf-8', errors='ignore')

def test_newbie_area():
    """Test the newbie area changes."""
    print("=" * 60)
    print("TESTING NEWBIE AREA AND NEW FEATURES")
    print("=" * 60)
    
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    sock.connect(('localhost', 4000))
    
    initial = send_command(sock, "")
    print("\n[CONNECT] Initial prompt received")
    
    # Create character
    print("\n[CREATING CHARACTER]")
    send_command(sock, "TestChar")
    print("  → Name: TestChar")
    
    time.sleep(0.3)
    send_command(sock, "human")
    print("  → Race: human")
    
    time.sleep(0.3)
    send_command(sock, "scholar")
    print("  → Kit: scholar")
    
    time.sleep(0.3)
    send_command(sock, "male")
    print("  → Sex: male")
    
    time.sleep(0.5)
    out = send_command(sock, "")
    
    # Verify in newbie area
    if "8101" in out or "Newbie" in out or "Welcome" in out:
        print("  ✓ Successfully in newbie area")
    
    # Test room descriptions
    print("\n[TESTING ROOM DESCRIPTIONS]")
    
    # Move around and check for OLD content
    bad_content = [
        "automap", "compass", "[Desc]", "scroll",
        "sing", "sayto", "whisper", "yell",
        "[Kill]", "[Consider]",
        "ask", "answer", "channels"
    ]
    
    for _ in range(3):
        send_command(sock, "north")
        time.sleep(0.2)
    
    room_output = send_command(sock, "look")
    
    found_bad = []
    for bad in bad_content:
        if bad.lower() in room_output.lower():
            found_bad.append(bad)
    
    if found_bad:
        print(f"  ✗ Found outdated content: {found_bad}")
    else:
        print("  ✓ No outdated command references found")
    
    # Test chat command
    print("\n[TESTING CHAT COMMAND]")
    chat_out = send_command(sock, "chat Hello from chat!")
    if "Chat" in chat_out or "Hello from chat" in chat_out:
        print("  ✓ Chat command works")
    else:
        print("  ✗ Chat command may have issues")
    
    # Test help rules
    print("\n[TESTING HELP RULES]")
    help_out = send_command(sock, "help rules")
    if "rules" in help_out.lower() and ("botting" in help_out.lower() or "responsible" in help_out.lower()):
        print("  ✓ Help RULES displays correctly")
    else:
        print("  ✗ Help RULES may not be loading properly")
        print("  Response:", help_out[:200])
    
    sock.close()
    print("\n" + "=" * 60)
    print("TEST COMPLETE")
    print("=" * 60)

if __name__ == "__main__":
    test_newbie_area()
