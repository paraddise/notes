import requests
from flask import Flask, render_template, request, redirect, url_for, send_file, session, flash
from Crypto.Cipher import AES
from minio import Minio
import os, io, psycopg2, bcrypt, uuid
from datetime import timedelta
from werkzeug.utils import secure_filename
import base64, time

time.sleep(5)

app = Flask(__name__)
app.secret_key = os.environ.get('SECRET_KEY')
key = os.environ.get('APP_ENCRYPTION_KEY').encode()

# Set up MinIO client
minio_client = Minio(
    os.environ.get('MINIO_ENDPOINT'),
    access_key=os.environ.get('MINIO_ACCESS_KEY'),
    secret_key=os.environ.get('MINIO_SECRET_KEY'),
    secure=False
)

# Set up PostgreSQL connection
conn = psycopg2.connect(
        host="127.0.0.1",
        database="s3curestorage",
        user="s3curestorage",
        password="<masked>")


# AES decryption function
def decrypt_file(nonce, tag, ciphertext, key):
    cipher = AES.new(key, AES.MODE_EAX, nonce)
    plaintext = cipher.decrypt_and_verify(ciphertext, tag)

    return plaintext

# AES encryption function
def encrypt_file(file_data, key):
    cipher = AES.new(key, AES.MODE_EAX)
    ciphertext, tag = cipher.encrypt_and_digest(file_data)
    return cipher.nonce, tag, ciphertext

@app.route('/login', methods=['GET', 'POST'])
def login():
    if request.method == 'POST':
        username = request.form['username']
        password = request.form['password']

        cur = conn.cursor()
        cur.execute("SELECT password FROM users WHERE username = %s", (username,))
        user = cur.fetchone()
        cur.close()

        if user and bcrypt.checkpw(password.encode()+os.environ.get('SALT').encode(), user[0].encode()):
            session['username'] = username
            return redirect(url_for('index'))
        else:
            flash('Invalid username or password', 'error')

    return render_template('login.html')

@app.route('/logout')
def logout():
    session.pop('username', None)
    return redirect(url_for('login'))

@app.route('/')
def index():
    if 'username' not in session:
        return redirect(url_for('login'))
    # Retrieve file information from the database
    cur = conn.cursor()
    cur.execute("SELECT id, filename, description FROM files")
    files = cur.fetchall()
    cur.close()
    return render_template('index.html', files=files)


@app.route('/download/<int:file_id>')
def download(file_id):
    if 'username' not in session:
        return redirect(url_for('login'))
    # Retrieve file information from the database
    cur = conn.cursor()
    cur.execute("SELECT location, description, filename FROM files WHERE id = %s", (file_id,))
    file_info = cur.fetchone()
    cur.close()

    if file_info:
        # Download file from MinIO
        url = minio_client.presigned_get_object("s3curestorage", file_info[0], expires = timedelta(hours=2))
        r=requests.get(url)
        encrypted_data = r.content
        object_data = minio_client.stat_object('s3curestorage', file_info[0])
        nonce = base64.b64decode(object_data.metadata['x-amz-meta-nonce'])
        tag = base64.b64decode(object_data.metadata['x-amz-meta-tag'])


        # Decrypt the file
        decrypted_data = decrypt_file(nonce, tag, encrypted_data, key)

        # Provide the decrypted file as a download
        return send_file(
            io.BytesIO(decrypted_data),
            mimetype='application/octet-stream',
            as_attachment=True,
            download_name=file_info[2]
        )
    else:
        return "File not found."




@app.route('/upload', methods=['GET', 'POST'])
def upload_file():
    if 'username' not in session:
        return redirect(url_for('login'))
    if request.method == 'POST':
        # Get the uploaded file, filename, and description
        uploaded_file = request.files['file']
        #filename = request.form['filename']
        filename = secure_filename(uploaded_file.filename)
        description = request.form['description']

        # Encrypt the file with AES
        file_data = uploaded_file.read()
        nonce, tag, ciphertext = encrypt_file(file_data, key)
        object_name = str(uuid.uuid4())
        tempfile="/tmp/"+object_name
        with open (tempfile,"wb") as f:
            f.write(ciphertext)
        # Store the encrypted file in MinIO
        with open(tempfile, "rb") as f:
            minio_client.put_object(
                bucket_name='s3curestorage',
                object_name=object_name,
                data=f,
                length=len(ciphertext),
                metadata={'nonce': base64.b64encode(nonce).decode(), 'tag': base64.b64encode(tag).decode()}
            )
        os.remove(tempfile)


        # Save file info in PostgreSQL
        cur = conn.cursor()
        cur.execute("INSERT INTO files (filename, location, description) VALUES (%s, %s, %s)",
                    (filename, object_name, description))
        conn.commit()
        cur.close()

        return redirect(url_for('index'))

    return render_template('upload.html')

if __name__ == '__main__':
    app.run(host="0.0.0.0")
