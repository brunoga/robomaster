/*
 * Copyright 1995-2019 Bruno Haible <bruno@clisp.org>
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

#ifndef _VACALL_R_H
#define _VACALL_R_H

#include <stddef.h>

#include "ffcall-abi.h"


/* Use a consistent prefix for all symbols in libcallback. */
#define vacall_start                   callback_start
#define vacall_start_struct            callback_start_struct
#define vacall_arg_char                callback_arg_char
#define vacall_arg_schar               callback_arg_schar
#define vacall_arg_uchar               callback_arg_uchar
#define vacall_arg_short               callback_arg_short
#define vacall_arg_ushort              callback_arg_ushort
#define vacall_arg_int                 callback_arg_int
#define vacall_arg_uint                callback_arg_uint
#define vacall_arg_long                callback_arg_long
#define vacall_arg_ulong               callback_arg_ulong
#define vacall_arg_longlong            callback_arg_longlong
#define vacall_arg_ulonglong           callback_arg_ulonglong
#define vacall_arg_float               callback_arg_float
#define vacall_arg_double              callback_arg_double
#define vacall_arg_ptr                 callback_arg_ptr
#define vacall_arg_struct              callback_arg_struct
#define vacall_return_void             callback_return_void
#define vacall_return_char             callback_return_char
#define vacall_return_schar            callback_return_schar
#define vacall_return_uchar            callback_return_uchar
#define vacall_return_short            callback_return_short
#define vacall_return_ushort           callback_return_ushort
#define vacall_return_int              callback_return_int
#define vacall_return_uint             callback_return_uint
#define vacall_return_long             callback_return_long
#define vacall_return_ulong            callback_return_ulong
#define vacall_return_longlong         callback_return_longlong
#define vacall_return_ulonglong        callback_return_ulonglong
#define vacall_return_float            callback_return_float
#define vacall_return_double           callback_return_double
#define vacall_return_ptr              callback_return_ptr
#define vacall_return_struct           callback_return_struct
#define vacall_error_type_mismatch     callback_error_type_mismatch
#define vacall_error_struct_too_large  callback_error_struct_too_large
#define vacall_structcpy               callback_structcpy
#define vacall_struct_buffer           callback_struct_buffer


/* Determine whether the current ABI is LLP64
   ('long' = 32-bit, 'long long' = 'void*' = 64-bit). */
#if defined(__x86_64__) && defined(_WIN32) && !defined(__CYGWIN__)
#define __VA_LLP64 1
#endif

/* Determine the alignment of a type at compile time.
 */
#if defined(__GNUC__) || defined(__IBM__ALIGNOF__)
#define __VA_alignof __alignof__
#elif defined(__cplusplus)
template <class type> struct __VA_alignof_helper { char __slot1; type __slot2; };
#define __VA_alignof(type) offsetof (__VA_alignof_helper<type>, __slot2)
#elif defined(__mips__) || defined(__mipsn32__) || defined(__mips64__) /* SGI compiler */
#define __VA_alignof __builtin_alignof
#else
#define __VA_offsetof(type,ident)  ((unsigned long)&(((type*)0)->ident))
#define __VA_alignof(type)  __VA_offsetof(struct { char __slot1; type __slot2; }, __slot2)
#endif

#ifdef __cplusplus
extern "C" {
#endif

/* C builtin types.
 */
#if defined(__mipsn32__) || defined(__x86_64_x32__) || defined(__VA_LLP64)
typedef long long __vaword;
#else
typedef long __vaword;
#endif

enum __VAtype
{
  __VAvoid,
  __VAchar,
  __VAschar,
  __VAuchar,
  __VAshort,
  __VAushort,
  __VAint,
  __VAuint,
  __VAlong,
  __VAulong,
  __VAlonglong,
  __VAulonglong,
  __VAfloat,
  __VAdouble,
  __VAvoidp,
  __VAstruct
};

enum __VA_alist_flags
{

  /* how to return structs */
  /* There are basically 3 ways to return structs:
   * a. The called function returns a pointer to static data. Not reentrant.
   *    Not supported any more.
   * b. The caller passes the return structure address in a dedicated register
   *    or as a first (or last), invisible argument. The called function stores
   *    its result there.
   * c. Like b, and the called function also returns the return structure
   *    address in the return value register. (This is not very distinguishable
   *    from b.)
   * Independently of this,
   * r. small structures (<= 4 or <= 8 bytes) may be returned in the return
   *    value register(s), or
   * m. even small structures are passed in memory.
   */
  /* gcc-2.6.3 employs the following strategy:
   *   - If PCC_STATIC_STRUCT_RETURN is defined in the machine description
   *     it uses method a, else method c.
   *   - If flag_pcc_struct_return is set (either by -fpcc-struct-return or if
   *     DEFAULT_PCC_STRUCT_RETURN is defined to 1 in the machine description)
   *     it uses method m, else (either by -freg-struct-return or if
   *     DEFAULT_PCC_STRUCT_RETURN is defined to 0 in the machine description)
   *     method r.
   */
  __VA_SMALL_STRUCT_RETURN	= 1<<1,	/* r: special case for small structs */
  __VA_GCC_STRUCT_RETURN	= 1<<2,	/* consider 8 byte structs as small */
#if defined(__sparc__) && !defined(__sparc64__)
  __VA_SUNCC_STRUCT_RETURN	= 1<<3,
  __VA_SUNPROCC_STRUCT_RETURN	= 1<<4,
#else
  __VA_SUNCC_STRUCT_RETURN	= 0,
  __VA_SUNPROCC_STRUCT_RETURN	= 0,
#endif
#if defined(__i386__)
  __VA_MSVC_STRUCT_RETURN	= 1<<4,
#endif
  /* the default way to return structs */
  /* This choice here is based on the assumption that the function you are
   * going to call has been compiled with the same compiler you are using to
   * include this file.
   * If you want to call functions with another struct returning convention,
   * just  #define __VA_STRUCT_RETURN ...
   * before or after #including <vacall_r.h>.
   */
#ifndef __VA_STRUCT_RETURN
  __VA_STRUCT_RETURN		=
#if defined(__sparc__) && !defined(__sparc64__) && defined(__sun) && (defined(__SUNPRO_C) || defined(__SUNPRO_CC)) /* SUNWspro cc or CC */
				  __VA_SUNPROCC_STRUCT_RETURN,
#else
#if (defined(__i386__) && (defined(_WIN32) || defined(__CYGWIN__) || (defined(__MACH__) && defined(__APPLE__)) || defined(__FreeBSD__) || defined(__DragonFly__) || defined(__OpenBSD__))) || defined(__m68k__) || defined(__mipsn32__) || defined(__mips64__) || defined(__sparc64__) || defined(__hppa__) || defined(__hppa64__) || defined(__arm__) || defined(__armhf__) || defined(__arm64__) || defined(__powerpc64_elfv2__) || defined(__ia64__) || defined(__x86_64__) || defined(__riscv32__) || defined(__riscv64__)
				  __VA_SMALL_STRUCT_RETURN |
#endif
#if defined(__GNUC__) && !((defined(__mipsn32__) || defined(__mips64__)) && ((__GNUC__ == 3 && __GNUC_MINOR__ >= 4) || (__GNUC__ > 3)))
				  __VA_GCC_STRUCT_RETURN |
#endif
#if defined(__i386__) && defined(_WIN32) && !defined(__CYGWIN__) /* native Windows */
				  __VA_MSVC_STRUCT_RETURN |
#endif
				  0,
#endif
#endif

  /* how to return floats */
#if defined(__m68k__) || (defined(__sparc__) && !defined(__sparc64__))
  __VA_SUNCC_FLOAT_RETURN	= 1<<5,
#endif
#if defined(__m68k__)
  __VA_FREG_FLOAT_RETURN	= 1<<6,
#endif
  /* the default way to return floats */
  /* This choice here is based on the assumption that the function you are
   * going to call has been compiled with the same compiler you are using to
   * include this file.
   * If you want to call functions with another float returning convention,
   * just  #define __VA_FLOAT_RETURN ...
   * before or after #including <vacall_r.h>.
   */
#ifndef __VA_FLOAT_RETURN
#if (defined(__m68k__) || (defined(__sparc__) && !defined(__sparc64__))) && !defined(__GNUC__) && defined(__sun) && !(defined(__SUNPRO_C) || defined(__SUNPRO_CC))  /* Sun cc or CC */
  __VA_FLOAT_RETURN		= __VA_SUNCC_FLOAT_RETURN,
#elif defined(__m68k__)
  __VA_FLOAT_RETURN		= __VA_FREG_FLOAT_RETURN,
#else
  __VA_FLOAT_RETURN		= 0,
#endif
#endif

  /* how to pass structs */
#if defined(__mips__) || defined(__mipsn32__) || defined(__mips64__)
  __VA_SGICC_STRUCT_ARGS	= 1<<7,
#endif
#if defined(__powerpc__) || defined(__powerpc64__)
  __VA_AIXCC_STRUCT_ARGS	= 1<<7,
#endif
#if defined(__ia64__)
  __VA_OLDGCC_STRUCT_ARGS	= 1<<7,
#endif
  /* the default way to pass structs */
  /* This choice here is based on the assumption that the function you are
   * going to call has been compiled with the same compiler you are using to
   * include this file.
   * If you want to call functions with another structs passing convention,
   * just  #define __VA_STRUCT_ARGS ...
   * before or after #including <vacall_r.h>.
   */
#ifndef __VA_STRUCT_ARGS
#if (defined(__mips__) && !defined(__mipsn32__) && !defined(__mips64__)) && !defined(__GNUC__) /* SGI mips cc */
  __VA_STRUCT_ARGS		= __VA_SGICC_STRUCT_ARGS,
#else
#if (defined(__mipsn32__) || defined(__mips64__)) && (!defined(__GNUC__) || (__GNUC__ == 3 && __GNUC_MINOR__ >= 4) || (__GNUC__ > 3)) /* SGI mips cc or gcc >= 3.4 */
  __VA_STRUCT_ARGS		= __VA_SGICC_STRUCT_ARGS,
#else
#if defined(__powerpc__) && !defined(__powerpc64__) && defined(_AIX) && !defined(__GNUC__) /* AIX 32-bit cc, xlc */
  __VA_STRUCT_ARGS		= __VA_AIXCC_STRUCT_ARGS,
#else
#if defined(__powerpc64__) && defined(_AIX) /* AIX 64-bit cc, xlc, gcc */
  __VA_STRUCT_ARGS		= __VA_AIXCC_STRUCT_ARGS,
#else
#if defined(__ia64__) && !(defined(__GNUC__) && (__GNUC__ >= 3))
  __VA_STRUCT_ARGS		= __VA_OLDGCC_STRUCT_ARGS,
#else
  __VA_STRUCT_ARGS		= 0,
#endif
#endif
#endif
#endif
#endif
#endif

  /* how to pass floats */
  /* ANSI C compilers and GNU gcc pass floats as floats.
   * K&R C compilers pass floats as doubles. We don't support them any more.
   */
#if defined(__powerpc64__)
  __VA_AIXCC_FLOAT_ARGS         = 1<<8,      /* pass floats in the low 4 bytes of an 8-bytes word */
#endif
  /* the default way to pass floats */
  /* This choice here is based on the assumption that the function you are
   * going to call has been compiled with the same compiler you are using to
   * include this file.
   * If you want to call functions with another float passing convention,
   * just  #define __VA_FLOAT_ARGS ...
   * before or after #including <vacall_r.h>.
   */
#ifndef __VA_FLOAT_ARGS
#if defined(__powerpc64__) && defined(_AIX) && !defined(__GNUC__) /* AIX 64-bit xlc */
  __VA_FLOAT_ARGS		= __VA_AIXCC_FLOAT_ARGS,
#else
  __VA_FLOAT_ARGS		= 0,
#endif
#endif

  /* how to pass and return small integer arguments */
  __VA_ANSI_INTEGERS		= 0, /* no promotions */
  __VA_TRADITIONAL_INTEGERS	= 0, /* promote [u]char, [u]short to [u]int */
  /* Fortunately these two methods are compatible. Our macros work with both. */

  /* stack cleanup policy */
  __VA_CDECL_CLEANUP		= 0, /* caller pops args after return */
  __VA_STDCALL_CLEANUP		= 1<<9, /* callee pops args before return */
				     /* currently only supported on __i386__ */
#ifndef __VA_CLEANUP
  __VA_CLEANUP			= __VA_CDECL_CLEANUP,
#endif

  /* These are for internal use only */
#if defined(__i386__) || defined(__m68k__) || defined(__mipsn32__) || defined(__mips64__) || defined(__sparc64__) || defined(__alpha__) || defined(__hppa64__) || defined(__arm__) || defined(__armhf__) || defined(__arm64__) || defined(__powerpc__) || defined(__powerpc64__) || defined(__ia64__) || defined(__x86_64__) || (defined(__s390__) && !defined(__s390x__)) || defined(__riscv64__)
  __VA_REGISTER_STRUCT_RETURN	= 1<<10,
#endif
#if defined(__mipsn32__) || defined(__mips64__)
  __VA_REGISTER_FLOATSTRUCT_RETURN	= 1<<11,
  __VA_REGISTER_DOUBLESTRUCT_RETURN	= 1<<12,
#endif

  __VA_flag_for_broken_compilers_that_dont_like_trailing_commas
};

/*
 * Definition of the ‘va_alist’ type.
 */
struct vacall_alist
/* GNU clisp pokes in internals of the alist! */
#ifdef LISPFUN
{
  /* some va_... macros need these flags */
  int            flags;
#if defined(__i386__) || defined(__arm__) || defined(__armhf__) || (defined(__powerpc__) && !defined(__powerpc64__) && defined(__MACH__) && defined(__APPLE__))
  __vaword       filler1;
#endif
  /* temporary storage for return value */
  union {
    char                _char;
    signed char         _schar;
    unsigned char       _uchar;
    short               _short;
    unsigned short      _ushort;
    int                 _int;
    unsigned int        _uint;
    long                _long;
    unsigned long       _ulong;
    long long           _longlong;
    unsigned long long  _ulonglong;
    float               _float;
    double              _double;
    void*               _ptr;
  }              tmp;
}
#endif
;
typedef struct vacall_alist * va_alist;


/*
 * Definition of the va_start_xxx macros.
 */
#define __VA_START_FLAGS  \
  __VA_STRUCT_RETURN | __VA_FLOAT_RETURN | __VA_STRUCT_ARGS | __VA_FLOAT_ARGS | __VA_CLEANUP

extern void vacall_start (va_alist /* LIST */, int /* RETTYPE */, int /* FLAGS */);

#define va_start_void(LIST)	 vacall_start(LIST,__VAvoid,     __VA_START_FLAGS)
#define va_start_char(LIST)	 vacall_start(LIST,__VAchar,     __VA_START_FLAGS)
#define va_start_schar(LIST)	 vacall_start(LIST,__VAschar,    __VA_START_FLAGS)
#define va_start_uchar(LIST)	 vacall_start(LIST,__VAuchar,    __VA_START_FLAGS)
#define va_start_short(LIST)	 vacall_start(LIST,__VAshort,    __VA_START_FLAGS)
#define va_start_ushort(LIST)	 vacall_start(LIST,__VAushort,   __VA_START_FLAGS)
#define va_start_int(LIST)	 vacall_start(LIST,__VAint,      __VA_START_FLAGS)
#define va_start_uint(LIST)	 vacall_start(LIST,__VAuint,     __VA_START_FLAGS)
#define va_start_long(LIST)	 vacall_start(LIST,__VAlong,     __VA_START_FLAGS)
#define va_start_ulong(LIST)	 vacall_start(LIST,__VAulong,    __VA_START_FLAGS)
#define va_start_longlong(LIST)	 vacall_start(LIST,__VAlonglong, __VA_START_FLAGS)
#define va_start_ulonglong(LIST) vacall_start(LIST,__VAulonglong,__VA_START_FLAGS)
#define va_start_float(LIST)	 vacall_start(LIST,__VAfloat,    __VA_START_FLAGS)
#define va_start_double(LIST)	 vacall_start(LIST,__VAdouble,   __VA_START_FLAGS)
#define va_start_ptr(LIST,TYPE)	 vacall_start(LIST,__VAvoidp,    __VA_START_FLAGS)

/*
 * va_start_struct: Preparing structure return.
 */
extern void vacall_start_struct (va_alist /* LIST */, size_t /* TYPE_SIZE */, size_t /* TYPE_ALIGN */, int /* TYPE_SPLITTABLE */, int /* FLAGS */);

#define va_start_struct(LIST,TYPE,TYPE_SPLITTABLE)  \
  _va_start_struct(LIST,sizeof(TYPE),__VA_alignof(TYPE),TYPE_SPLITTABLE)
/* _va_start_struct() is like va_start_struct(), except that you pass
 * the type's size and alignment instead of the type itself.
 * Undocumented, but used by GNU clisp.
 */
#define _va_start_struct(LIST,TYPE_SIZE,TYPE_ALIGN,TYPE_SPLITTABLE)  \
  vacall_start_struct(LIST,TYPE_SIZE,TYPE_ALIGN,TYPE_SPLITTABLE,__VA_START_FLAGS)


/*
 * Definition of the va_arg_xxx macros.
 */

extern char           vacall_arg_char   (va_alist /* LIST */);
extern signed char    vacall_arg_schar  (va_alist /* LIST */);
extern unsigned char  vacall_arg_uchar  (va_alist /* LIST */);
extern short          vacall_arg_short  (va_alist /* LIST */);
extern unsigned short vacall_arg_ushort (va_alist /* LIST */);
extern int            vacall_arg_int    (va_alist /* LIST */);
extern unsigned int   vacall_arg_uint   (va_alist /* LIST */);
extern long           vacall_arg_long   (va_alist /* LIST */);
extern unsigned long  vacall_arg_ulong  (va_alist /* LIST */);

#define va_arg_char(LIST)	vacall_arg_char(LIST)
#define va_arg_schar(LIST)	vacall_arg_schar(LIST)
#define va_arg_uchar(LIST)	vacall_arg_uchar(LIST)
#define va_arg_short(LIST)	vacall_arg_short(LIST)
#define va_arg_ushort(LIST)	vacall_arg_ushort(LIST)
#define va_arg_int(LIST)	vacall_arg_int(LIST)
#define va_arg_uint(LIST)	vacall_arg_uint(LIST)
#define va_arg_long(LIST)	vacall_arg_long(LIST)
#define va_arg_ulong(LIST)	vacall_arg_ulong(LIST)

extern long long          vacall_arg_longlong  (va_alist /* LIST */);
extern unsigned long long vacall_arg_ulonglong (va_alist /* LIST */);

#define va_arg_longlong(LIST)	vacall_arg_longlong(LIST)
#define va_arg_ulonglong(LIST)	vacall_arg_ulonglong(LIST)

/* Floating point arguments. */

extern float  vacall_arg_float  (va_alist /* LIST */);
extern double vacall_arg_double (va_alist /* LIST */);

#define va_arg_float(LIST)	vacall_arg_float(LIST)
#define va_arg_double(LIST)	vacall_arg_double(LIST)

/* Pointer arguments. */

extern void* vacall_arg_ptr (va_alist /* LIST */);
#define va_arg_ptr(LIST,TYPE)	((TYPE)vacall_arg_ptr(LIST))

/* Structure arguments. */

extern void* vacall_arg_struct (va_alist /* LIST */, size_t /* TYPE_SIZE */, size_t /* TYPE_ALIGN */);

#define va_arg_struct(LIST,TYPE)  \
  *(TYPE*)vacall_arg_struct(LIST,sizeof(TYPE),__VA_alignof(TYPE))
/* _va_arg_struct() is like va_arg_struct(), except that you pass the type's
 * size and alignment instead of the type and get the value's address instead
 * of the value itself.
 * Undocumented, but used by GNU clisp.
 */
#define _va_arg_struct(LIST,TYPE_SIZE,TYPE_ALIGN)  \
  vacall_arg_struct(LIST,TYPE_SIZE,TYPE_ALIGN)


/*
 * Definition of the va_return_xxx macros.
 */

extern void vacall_return_void (va_alist /* LIST */);
#define va_return_void(LIST)		vacall_return_void(LIST)

extern void vacall_return_char (va_alist /* LIST */, char /* VAL */);
extern void vacall_return_schar (va_alist /* LIST */, signed char /* VAL */);
extern void vacall_return_uchar (va_alist /* LIST */, unsigned char /* VAL */);
extern void vacall_return_short (va_alist /* LIST */, short /* VAL */);
extern void vacall_return_ushort (va_alist /* LIST */, unsigned short /* VAL */);
extern void vacall_return_int (va_alist /* LIST */, int /* VAL */);
extern void vacall_return_uint (va_alist /* LIST */, unsigned int /* VAL */);
extern void vacall_return_long (va_alist /* LIST */, long /* VAL */);
extern void vacall_return_ulong (va_alist /* LIST */, unsigned long /* VAL */);
#define va_return_char(LIST,VAL)	vacall_return_char(LIST,VAL)
#define va_return_schar(LIST,VAL)	vacall_return_schar(LIST,VAL)
#define va_return_uchar(LIST,VAL)	vacall_return_uchar(LIST,VAL)
#define va_return_short(LIST,VAL)	vacall_return_short(LIST,VAL)
#define va_return_ushort(LIST,VAL)	vacall_return_ushort(LIST,VAL)
#define va_return_int(LIST,VAL)		vacall_return_int(LIST,VAL)
#define va_return_uint(LIST,VAL)	vacall_return_uint(LIST,VAL)
#define va_return_long(LIST,VAL)	vacall_return_long(LIST,VAL)
#define va_return_ulong(LIST,VAL)	vacall_return_ulong(LIST,VAL)

extern void vacall_return_longlong (va_alist /* LIST */, long long /* VAL */);
extern void vacall_return_ulonglong (va_alist /* LIST */, unsigned long long /* VAL */);
#define va_return_longlong(LIST,VAL)	vacall_return_longlong(LIST,VAL)
#define va_return_ulonglong(LIST,VAL)	vacall_return_ulonglong(LIST,VAL)

extern void vacall_return_float (va_alist /* LIST */, float /* VAL */);
extern void vacall_return_double (va_alist /* LIST */, double /* VAL */);
#define va_return_float(LIST,VAL)	vacall_return_float(LIST,VAL)
#define va_return_double(LIST,VAL)	vacall_return_double(LIST,VAL)

extern void vacall_return_ptr (va_alist /* LIST */, void* /* VAL */);
#define va_return_ptr(LIST,TYPE,VAL)	vacall_return_ptr(LIST,(void*)(TYPE)(VAL))

extern void vacall_return_struct (va_alist /* LIST */, size_t /* TYPE_SIZE */, size_t /* TYPE_ALIGN */, const void* /* VAL_ADDR */);

#define va_return_struct(LIST,TYPE,VAL)  \
  _va_return_struct(LIST,sizeof(TYPE),__VA_alignof(TYPE),&(VAL))
/* Undocumented, but used by GNU clisp. */
#define _va_return_struct(LIST,TYPE_SIZE,TYPE_ALIGN,VAL_ADDR)  \
  vacall_return_struct(LIST,TYPE_SIZE,TYPE_ALIGN,VAL_ADDR)


/* Determine whether a struct type is word-splittable, i.e. whether each of
 * its components fit into a register.
 * The entire computation is done at compile time.
 */
#define va_word_splittable_1(slot1)  \
  (__va_offset1(slot1)/sizeof(__vaword) == (__va_offset1(slot1)+sizeof(slot1)-1)/sizeof(__vaword))
#define va_word_splittable_2(slot1,slot2)  \
  ((__va_offset1(slot1)/sizeof(__vaword) == (__va_offset1(slot1)+sizeof(slot1)-1)/sizeof(__vaword)) \
   && (__va_offset2(slot1,slot2)/sizeof(__vaword) == (__va_offset2(slot1,slot2)+sizeof(slot2)-1)/sizeof(__vaword)) \
  )
#define va_word_splittable_3(slot1,slot2,slot3)  \
  ((__va_offset1(slot1)/sizeof(__vaword) == (__va_offset1(slot1)+sizeof(slot1)-1)/sizeof(__vaword)) \
   && (__va_offset2(slot1,slot2)/sizeof(__vaword) == (__va_offset2(slot1,slot2)+sizeof(slot2)-1)/sizeof(__vaword)) \
   && (__va_offset3(slot1,slot2,slot3)/sizeof(__vaword) == (__va_offset3(slot1,slot2,slot3)+sizeof(slot3)-1)/sizeof(__vaword)) \
  )
#define va_word_splittable_4(slot1,slot2,slot3,slot4)  \
  ((__va_offset1(slot1)/sizeof(__vaword) == (__va_offset1(slot1)+sizeof(slot1)-1)/sizeof(__vaword)) \
   && (__va_offset2(slot1,slot2)/sizeof(__vaword) == (__va_offset2(slot1,slot2)+sizeof(slot2)-1)/sizeof(__vaword)) \
   && (__va_offset3(slot1,slot2,slot3)/sizeof(__vaword) == (__va_offset3(slot1,slot2,slot3)+sizeof(slot3)-1)/sizeof(__vaword)) \
   && (__va_offset4(slot1,slot2,slot3,slot4)/sizeof(__vaword) == (__va_offset4(slot1,slot2,slot3,slot4)+sizeof(slot4)-1)/sizeof(__vaword)) \
  )
#define __va_offset1(slot1)  \
  0
#define __va_offset2(slot1,slot2)  \
  ((__va_offset1(slot1)+sizeof(slot1)+__VA_alignof(slot2)-1) & -(long)__VA_alignof(slot2))
#define __va_offset3(slot1,slot2,slot3)  \
  ((__va_offset2(slot1,slot2)+sizeof(slot2)+__VA_alignof(slot3)-1) & -(long)__VA_alignof(slot3))
#define __va_offset4(slot1,slot2,slot3,slot4)  \
  ((__va_offset3(slot1,slot2,slot3)+sizeof(slot3)+__VA_alignof(slot4)-1) & -(long)__VA_alignof(slot4))


/*
 * Miscellaneous declarations.
 */

#if defined(__sparc__) || defined(__sparc64__)
/* On SPARC, PIC compiled code crashes when used outside of a shared library.
   Therefore, don't use the callback_get_receiver indirection on this platform. */

extern
#ifdef __cplusplus
"C"
#endif
void callback_receiver (); /* Actually it takes arguments and returns values! */

#define callback_get_receiver() (&callback_receiver)

#else

/* A fake type for callback_receiver.
   Actually it takes arguments and returns values. */
typedef void (*__vacall_r_t) (void);

/* This function returns the address of callback_receiver.
   callback_receiver is not a global symbol, because on ELF platforms, functions
   with global visibility cannot accept additional arguments in registers. See
   elf-hack.txt for more details. */
extern
#ifdef __cplusplus
"C"
#endif
__vacall_r_t callback_get_receiver (void);

#endif


#ifdef __cplusplus
}
#endif

#endif /* _VACALL_R_H */
