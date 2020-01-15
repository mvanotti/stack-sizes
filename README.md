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

From the [llvm docs](http://releases.llvm.org/9.0.0/docs/CodeGenerator.html#emitting-function-stack-size-information):

> A section containing metadata on function stack sizes will be emitted when `TargetLoweringObjectFile::StackSizesSection` is not `null`, and `TargetOptions::EmitStackSizeSection` is set (`-stack-size-section`). The section will contain an array of pairs of function symbol values (pointer size) and stack sizes (unsigned `LEB128`). The stack size values only include the space allocated in the function prologue. Functions with dynamic stack allocations are not included

Which is triggered by using the [`-stack-size-section` in `llc`](https://llvm.org/docs/CommandGuide/llc.html#cmdoption-llc-stack-size-section):

> Emit the `.stack_sizes` section which contains stack size metadata. The section contains an array of pairs of function symbol values (pointer size) and stack sizes (unsigned `LEB128`). The stack size values only include the space allocated in the function prologue. Functions with dynamic stack allocations are not included.

You might also be interested in using the [`-Wframe-larger-than` clang argument](https://clang.llvm.org/docs/ClangCommandLineReference.html#cmdoption-clang-wframe-larger-than), which will warn you if one of your stack frames is bigger than the specified amount of bytes:

```shellsession
$ clang -Wall -Wextra -Wframe-larger-than=100 t.c -o t
t.c:4:6: warning: stack frame size of 296 bytes in function 'foo' [-Wframe-larger-than=]
void foo(uint8_t ind) {
     ^
1 warning generated.
```

