# Синиор телемастер

## Task
«Этсамое, я и тарелку крутил, и инструкцию перечитал, и даже тапком по телику жахнул — не фурычит, ёлки-палки. Ты посмотри там, чё да как, за мной не заржавеет», — говорит вам усатый мужичок и протягивает ресивер. Поковыряйтесь в прошивке и верните телевидение.

Доступ к приставке: t-tvlink-cz7nt56c.spbctf.net/
Логин и пароль — admin.

Исходный код приставки: tvlink_fc83f20.zip

Обновление от производителя: tvlink-firmware_345b8e0.zip

## Solution

It's hash length extension attack.


Bruteforce secret length for signature on server
```bash
python3 solution.py bruteforce
```

Create simple patch that will match signature.
```bash
python3 solution.py create-patch <sescret-length-on-server>
```

Create and send zip archive with zip arhcive inside with `__main__.py` file inside.
```bash
python3 solution.py zip-in-zip <sescret-length-on-server>
```
