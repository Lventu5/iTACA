#!/usr/bin/python3
import chromadb
import sys

HOST = 'localhost'
PORT = 8000
COLL = 'payloads'

def main():
    args = sys.argv[1:]
    client = chromadb.HttpClient(HOST, PORT)
    coll = client.get_collection(COLL)

    query_res = coll.query(query_texts = args, n_results=5)
    res = []
    for i in range(5):
        if query_res['distances'][0][i] < 1.0:
            res.append(query_res['ids'][0][i])
        else:
            res.append('safe')

    if res.count('safe') != 5:
        print(res)

if __name__ == '__main__':
    main()