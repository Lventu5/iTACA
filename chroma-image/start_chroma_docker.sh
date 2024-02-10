#!/bin/bash

docker run -p 8000:8000 -d chroma-server 
sleep 5 
python3 restoreDb.py
