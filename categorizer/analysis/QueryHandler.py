#!/usr/bin/python3
import chromadb
import sys

def main():
    HOST = sys.argv[1]
    PORT = int(sys.argv[2])
    COLL = sys.argv[3]
    args = sys.argv[4:]
    client = chromadb.HttpClient(HOST, PORT)
    coll = client.get_collection(COLL)

    query_res = coll.query(query_texts = args, n_results=5)
    res = []
    for i in range(5):
        if query_res['distances'][0][i] < 1.5:
            res.append(query_res['ids'][0][i])
        else:
            res.append('safe')

    #if res.count('safe') != 5:
    print(res)

if __name__ == '__main__':
    main()