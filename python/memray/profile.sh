#!/bin/bash
memray run -o memray-example.py.1.bin example.py
memray flamegraph memray-example.py.1.bin
# Open in xdg-browser
open memray-flamegraph-example.py.1.html
