#!/bin/sh
cd $1
dlv debug --headless -log -l 0.0.0.0:$2 --api-version=2
