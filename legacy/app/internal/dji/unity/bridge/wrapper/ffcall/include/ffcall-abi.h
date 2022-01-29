/*
 * Copyright 2017-2019 Bruno Haible <bruno@clisp.org>
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 2 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

/* Define some canonical CPU and ABI indicators.
   References:
     - host-cpu-c-abi.m4 from gnulib
     - https://sourceforge.net/p/predef/wiki/Architectures/
     - GCC source code: definitions of macro TARGET_CPU_CPP_BUILTINS
     - clang source code: defineMacro invocations in Basic/Targets.cpp,
       especially in getTargetDefines methods.
   Limitation: Unlike host-cpu-c-abi.m4, this preprocessor-based approach
   can not reliably distinguish __arm__ and __armhf__.
 */

#ifndef __i386__
#if defined(__i386__) /* GCC, clang */ || defined(__i386) /* Sun C */ || defined(_M_IX86) /* MSVC */
#define __i386__ 1
#endif
#endif

#ifndef __m68k__
#if defined(__m68k__) /* GCC */
#define __m68k__ 1
#endif
#endif

/* On mips, there are three ABIs:
   - 32 or o32: It defines _MIPS_SIM == _ABIO32 and _MIPS_SZLONG == 32.
   - n32: It defines _MIPS_SIM == _ABIN32 and _MIPS_SZLONG == 32.
   - 64: It defines _MIPS_SZLONG == 64.
 */
/* Note: When __mipsn32__ or __mips64__ is defined, __mips__ may or may not be
   defined as well. To test for the MIPS o32 ABI, use
     #if defined(__mips__) && !defined(__mipsn32__) && !defined(__mips64__)
 */
/* To distinguish little-endian and big-endian arm, use the preprocessor
   defines _MIPSEB vs. _MIPSEL. */
#ifndef __mips__
#if defined(__mips) /* GCC, clang, IRIX cc */ /* Note: GCC, clang also define __mips__. */
#define __mips__ 1
#endif
#endif
#ifndef __mipsn32__
#if defined(__mips__) && (_MIPS_SIM == _ABIN32)
#define __mipsn32__ 1
#endif
#endif
#ifndef __mips64__
#if defined(__mips__) && defined(_MIPS_SZLONG) && (_MIPS_SZLONG == 64)
#define __mips64__ 1
#endif
#endif

/* Note: When __sparc64__ is defined, __sparc__ may or may not be defined as
   well. To test for the SPARC 32-bit ABI, use
     #if defined(__sparc__) && !defined(__sparc64__)
 */
#ifndef __sparc__
#if defined(__sparc) /* GCC, clang, Sun C */ /* Note: GCC, clang also define __sparc__. */
#define __sparc__ 1
#endif
#endif
#ifndef __sparc64__
#if defined(__sparcv9) /* GCC/Solaris, Sun C */ || defined(__arch64__) /* GCC/Linux */
#define __sparc64__ 1
#endif
#endif

#ifndef __alpha__
#if defined(__alpha) /* GCC, DEC C */ /* Note: GCC also defines __alpha__. */
#define __alpha__ 1
#endif
#endif

/* On hppa, the C compiler may be generating 32-bit code or 64-bit code.
   In the latter case, it defines _LP64 and __LP64__.
 */
/* Note: When __hppa64__ is defined, __hppa__ may or may not be defined as well.
   To test for the HP-PA 32-bit ABI, use
     #if defined(__hppa__) && !defined(__hppa64__)
 */
#ifndef __hppa__
#if defined(__hppa) /* GCC, HP C */ /* Note: GCC also defines __hppa__. */
#define __hppa__ 1
#endif
#endif
#ifndef __hppa64__
#if defined(__hppa__) && defined(__LP64__)
#define __hppa64__ 1
#endif
#endif

/* Distinguish arm which passes floating-point arguments and return values
   in integer registers (r0, r1, ...) - this is gcc -mfloat-abi=soft or
   gcc -mfloat-abi=softfp - from arm which passes them in float registers
   (s0, s1, ...) and double registers (d0, d1, ...) - this is
   gcc -mfloat-abi=hard. GCC 4.6 or newer sets the preprocessor defines
   __ARM_PCS (for the first case) and __ARM_PCS_VFP (for the second case),
   but older GCC does not. */
/* Note: When __armhf__ is defined, __arm__ may or may not be defined as well.
   To test for the ARM ABI that does not use floating-point registers for
   parameter passing, use
     #if defined(__arm__) && !defined(__armhf__)
 */
/* To distinguish little-endian and big-endian arm, use the preprocessor
   defines __ARMEL__ vs. __ARMEB__. */
#ifndef __arm__
#if defined(__arm__) /* GCC, clang */ || defined(_M_ARM) /* MSVC */
#define __arm__ 1
#endif
#endif
#ifndef __armhf__
#if defined(__arm__) && defined(__ARM_PCS_VFP) /* GCC */
#define __armhf__ 1
#endif
#endif

/* On arm64 systems, the C compiler may be generating code in one of these ABIs:
   - aarch64 instruction set, 64-bit pointers, 64-bit 'long': arm64.
   - aarch64 instruction set, 32-bit pointers, 32-bit 'long': arm64-ilp32.
   - 32-bit instruction set, 32-bit pointers, 32-bit 'long': arm or armhf
     (see above).
 */
/* Note: When __arm64_ilp32__ is defined, __arm64__ may or may not be defined as
   well. To test for the arm64 64-bit ABI, use
     #if defined(__arm64__) && !defined(__arm64_ilp32__)
 */
/* To distinguish little-endian and big-endian arm64, use the preprocessor
   defines __AARCH64EL__ vs. __AARCH64EB__. */
#ifndef __arm64__
#if defined(__aarch64__) /* GCC, clang */ || defined(_M_ARM64) /* MSVC */
#define __arm64__ 1
#endif
#endif
#ifndef __arm64_ilp32__
#if defined(__arm64__) && (defined(__ILP32__) || defined (_ILP32))
#define __arm64_ilp32__ 1
#endif
#endif

/* On powerpc and powerpc64, different ABIs are in use on AIX vs. Mac OS X vs.
   Linux,*BSD. To distinguish them, use the OS dependent defines
     #if defined(_AIX)
     #if (defined(__MACH__) && defined(__APPLE__))
     #if !(defined(_AIX) || (defined(__MACH__) && defined(__APPLE__)))
 */
/* On powerpc64, there are two ABIs on Linux: The AIX compatible one and the
   ELFv2 one. The latter defines _CALL_ELF=2.
 */
/* Note: When __powerpc64__ is defined, __powerpc__ may or may not be defined as
   well. To test for the SPARC 32-bit ABI, use
     #if defined(__powerpc__) && !defined(__powerpc64__)
   Note: When __powerpc64_elfv2__ is defined, __powerpc64__ may or may not be
   defined as well. To test for the SPARC 32-bit ABI, use
     #if defined(__powerpc64__) && !defined(__powerpc64_elfv2__)
 */
#ifndef __powerpc__
#if defined(_ARCH_PPC) /* GCC, XLC */ /* Note: On AIX, Linux also __powerpc__ is defined; whereas on Mac OS X also __ppc__ is defined. On AIX also _IBMR2 is defined. */
#define __powerpc__ 1
#endif
#endif
#ifndef __powerpc64__
#if defined(_ARCH_PPC64) /* GCC, XLC */ /* Note: On Linux, also __powerpc64__ is defined. */
#define __powerpc64__ 1
#endif
#endif
#ifndef __powerpc64_elfv2__
#if defined(__powerpc64__) && defined(_CALL_ELF) && _CALL_ELF == 2
#define __powerpc64_elfv2__ 1
#endif
#endif

/* On ia64 on HP-UX, the C compiler may be generating 64-bit code or 32-bit
   code. In the latter case, it defines _ILP32.
 */
/* Note: When __ia64_ilp32__ is defined, __ia64__ may or may not be defined as
   well. To test for the ia64 64-bit ABI, use
     #if defined(__ia64__) && !defined(__ia64_ilp32__)
 */
#ifndef __ia64__
#if defined(__ia64__) /* GCC, HP C */ /* Note: GCC, HP C also define __ia64. */
#define __ia64__ 1
#endif
#endif
#ifndef __ia64_ilp32__
#if defined(__ia64__) && defined(_ILP32)
#define __ia64_ilp32__ 1
#endif
#endif

/* On x86_64 systems, the C compiler may be generating code in one of these ABIs:
   - 64-bit instruction set, 64-bit pointers, 64-bit 'long': x86_64.
   - 64-bit instruction set, 64-bit pointers, 32-bit 'long': x86_64
     with native Windows (mingw, MSVC).
   - 64-bit instruction set, 32-bit pointers, 32-bit 'long': x86_64-x32.
   - 32-bit instruction set, 32-bit pointers, 32-bit 'long': i386 (see above). */
/* Note: When __x86_64_x32__ is defined, __x86_64__ may or may not be defined as
   well. To test for the x86_64 64-bit ABI, use
     #if defined(__x86_64__) && !defined(__x86_64_x32__)
 */
#ifndef __x86_64__
#if (defined(__x86_64__) || defined(__amd64__)) /* GCC, clang, Sun C */ || (defined(_M_X64) || defined(_M_AMD64)) /* MSVC */
#define __x86_64__ 1
#endif
#endif
#ifndef __x86_64_x32__
#if defined(__x86_64__) && (defined(__ILP32__) || defined(_ILP32))
#define __x86_64_x32__ 1
#endif
#endif

/* Note: When __s390x__ is defined, __s390__ may or may not be defined as well.
   To test for the S/390 31-bit ABI, use
     #if defined(__s390__) && !defined(__s390x__)
 */
#ifndef __s390__
#if defined(__s390__) /* GCC, clang */
#define __s390__ 1
#endif
#endif
#ifndef __s390x__
#if defined(__s390x__) /* GCC, clang */
#define __s390x__ 1
#endif
#endif

#ifndef __riscv32__
#if defined(__riscv) && __riscv_xlen == 32 && !defined(__LP64__) /* GCC */
#define __riscv32__ 1
#endif
#endif

#ifndef __riscv64__
#if defined(__riscv) && __riscv_xlen == 64 && defined(__LP64__) /* GCC */
#define __riscv64__ 1
#endif
#endif
