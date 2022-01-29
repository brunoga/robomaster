###### Wine Bridge

This code implements support for a Linux and Windows (through Wine) program to
communicate through file descriptors (usually pipes).

The Linux code forks and then execs Wine to run the Windows program. It takes a
list of *os.File that need to be made available to the Windows program, taking
care of disabling close-on-exec on the associated file descriptors and passing
them over the command-line. On success the given files will be closed (on the
Linux side only, of course).

The Windows code parses the file descriptors from the command-line and uses
some Wine trickery to make them available to the Windows program as file
handles.
 