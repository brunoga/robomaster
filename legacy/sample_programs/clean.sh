#!/bin/bash
# Simple script to clean (Remove comments and empty lines of) a Python file.
# This is currently required because the Robomaster S1 Python editor does not
# seem to handle big Python programs very well. I am not sure this would help
# anyway, but it does not hurt.
grep -v '^$\|^\s*\#' $1 > $1.cleaned

