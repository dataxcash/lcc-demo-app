#!/bin/bash
# Stop all lcc-demo-app processes

pkill -f 'bin/web' 2>/dev/null
pkill -f 'bin/demo-app' 2>/dev/null
pkill -f 'lcc-web-demo' 2>/dev/null
sleep 0.5

exit 0
