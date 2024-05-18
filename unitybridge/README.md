# unitybridge

This is a high level API for the Robomaster Unity Bridge. It takes care of correctly initializing it and also exposes the key based interface to callers (keys can be read written and executed and those allow controlling the robot). It also exposes a per-event callback system so changes to keys can be monitored.

After this, most of the work is figuring out what each key does and when they should be used.
