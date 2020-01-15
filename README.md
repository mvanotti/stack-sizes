# stack-sizes
Utility that parses stack sizes section from elf objects and displays the preallocated stack size of each function.

## Usage

Compile a program with clang and the `-fstack-size-section` flag:

```c
#include <stdio.h>
#include <stdint.h>

void foo(uint8_t ind) {
  char a[256] = {0};

  scanf("%c\n", &(a[ind]));
}

int main(void) {
  uint8_t bar;

  scanf("%hhu", &bar);
  foo(bar);
  return 0;
}
```

```shellsession
$ clang -Wall -Wextra -fstack-size-section t.c -o t
$ go run stack-sizes.go -f t
0x401140        foo: 296 bytes
0x4011b0        main: 24 bytes
```

```c-objdump
$ objdump -d -Mintel t | grep -A 10 "<foo>:"
0000000000401140 <foo>:
  401140:       55                      push   rbp       ; 8
  401141:       48 89 e5                mov    rbp,rsp
  401144:       48 81 ec 20 01 00 00    sub    rsp,0x120 ; 288
  40114b:       40 88 f8                mov    al,dil
  40114e:       31 f6                   xor    esi,esi
  401150:       88 45 ff                mov    BYTE PTR [rbp-0x1],al
  401153:       48 8d 8d f0 fe ff ff    lea    rcx,[rbp-0x110]
  40115a:       48 89 ca                mov    rdx,rcx
  40115d:       48 89 d7                mov    rdi,rdx
  401160:       ba 00 01 00 00          mov    edx,0x100
$ objdump -d -Mintel t | grep -A 10 "<main>:"
00000000004011b0 <main>:
  4011b0:       55                      push   rbp      ; 8
  4011b1:       48 89 e5                mov    rbp,rsp
  4011b4:       48 83 ec 10             sub    rsp,0x10 ; 16
  4011b8:       c7 45 fc 00 00 00 00    mov    DWORD PTR [rbp-0x4],0x0
  4011bf:       48 bf 08 20 40 00 00    movabs rdi,0x402008
  4011c6:       00 00 00 
  4011c9:       48 8d 75 fb             lea    rsi,[rbp-0x5]
  4011cd:       b0 00                   mov    al,0x0
  4011cf:       e8 6c fe ff ff          call   401040 <__isoc99_scanf@plt>
  4011d4:       0f b6 7d fb             movzx  edi,BYTE PTR [rbp-0x5]
```
