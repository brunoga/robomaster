/*
 * Copyright 1997-2017 Bruno Haible <bruno@clisp.org>
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

#ifndef _CALLBACK_H
#define _CALLBACK_H

#include "ffcall-version.h"

/* Defines the type 'va_alist' and the va_* macros. */
#include "vacall_r.h"


#ifdef __cplusplus
extern "C" {
#endif


/* This type denotes an opaque function pointer.
   You need to cast it to an actual function pointer type (with correct return
   type) before you can actually invoke it. */
#ifdef __cplusplus
typedef int (*callback_t) (...);
#else
typedef int (*callback_t) ();
#endif
/* A deprecated alias of this type. */
typedef callback_t __TR_function;

/* This type denotes a callback implementation.
   DATA is the pointer that was passed to alloc_callback().
   ALIST allows to iterate over the argument list. */
typedef void (*callback_function_t) (void* /* DATA */, va_alist /* ALIST */);


/* Allocates a callback.
   It returns a function pointer whose signature depends on the behaviour
   of ADDRESS.
   When invoked, it passes DATA as first argument to ADDRESS and the actual
   arguments as a va_alist to ADDRESS. It returns the value passed to a
   va_return_* macro by ADDRESS.
   The callback has indefinite extent. It can be accessed until a call to
   free_callback().
 */
extern callback_t alloc_callback (callback_function_t /* ADDRESS */, void* /* DATA */);

/* Frees the memory used by a callback.
   CALLBACK must be the result of an alloc_callback() invocation.
   After this call, CALLBACK must not be used any more - neither invoked,
   not used as an argument to other functions.
 */
extern void free_callback (callback_t /* CALLBACK */);


/* Tests whether a given pointer is a function pointer returned by
   alloc_callback(). Returns 1 for yes, 0 for no.
   If yes, it can be cast to callback_t.
 */
extern int is_callback (void* /* CALLBACK */);

/* Returns the ADDRESS argument passed to the alloc_callback() invocation.
   CALLBACK must be the result of an alloc_callback() invocation.
 */
extern callback_function_t callback_address (callback_t /* CALLBACK */);

/* Returns the DATA argument passed to the alloc_callback() invocation.
   CALLBACK must be the result of an alloc_callback() invocation.
 */
extern void* callback_data (callback_t /* CALLBACK */);


#ifdef __cplusplus
}
#endif


#endif /* _CALLBACK_H */
