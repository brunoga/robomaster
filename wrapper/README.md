# wrapper
Go wrapper for Robomaster's Unity Bridge library

This is the low level interface with some niceties included (like supporting Go callbacks). It is not very useful unless you know how to use it as everything is based on specific events supported by it and used to communicate with the robot.

Supported platforms:

|Platform            |Support                                          |Test Status |
|--------------------|-------------------------------------------------|------------|
|Windows (AMD64)     |Native                                           |Tested      |
|MacOS (AMD64, ARM64)|Native for AMD64 and through Rosetta for ARM64   |Tested      |
|Linux (AMD64)       |Through a Wine bridge that hosts the Windows DLL |Tested      |
|iOS (ARM64)         |Native                                           |Tested      |
|Android (ARM, ARM64)|Native                                           |Not Tested  |
