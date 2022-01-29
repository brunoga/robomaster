### Unity Bridge

The Unity Bridge wraps up a Windows DLL and, due to that, only works on Windows 
(more specifically, Windows 64 bits). To work around that limitation, this code
uses a a bridge Windows program that wraps the DLL and runs under Wine using some
trickery to communicate with a Linux program through pipes. This makes sure the
Windows surface is as small as possible (unless you are actually running everything
on Windows, that is).  

###### Compiling on Linux

This requires 2 steps:

1 - Compile winewrapper.

This is the program that hosts the Windows DLL. It should be compiled and copied to
the same directory where the actual Linux executable is. To compile it you will need
the Windows C cross-compiler (mingw-w64, used through CGO) and to compile using Go in
cross-compilation mode.

CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build internal/wrapper/winewrapper

This should create a winewrapper.exe file in the current directory.

2 - Compile the Linux code.

This is as usual (go build program.go) but to run the resulting executable, as
mentioned above, you need both unitybridge.dll and winewrapper.exe to be in the
same directory as the resulting executable.

###### Running on Windows

Windows does not require the winewrapper at all so you just need to compile your
program as usual.

###### Running on other platforms

Technically this can run on any platform that is supported by Go and that has a
Wine port. It will not work out of the box because of the use of conditional
compilation based on file names but getting it to work onm other platforms should
be trivial.

DJI also makes available an equivalent of the unitybridge.dll file for MacOS so,
eventually, MacOS should also be able to run this without the need for the
winewrapper (if you want to volunteer to do this work, let me know).

