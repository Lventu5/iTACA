#!/bin/bash

docker pull chromadb/chroma
docker run -p 8000:8000 -d chromadb/chroma

sleep 5 
python3 `pwd`/restoreDb.py
