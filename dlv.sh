#!/bin/sh
cd /app/cmd
dlv debug --headless -log -l 0.0.0.0:2345 --api-version=2
