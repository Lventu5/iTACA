import chromadb
import os
import time
#from chromadb.config import Settings

#create a persistent vector database stored in the vectorDB directory
client = chromadb.PersistentClient(path="vectorDB/")

#create a collection where embeddings, documents and metadatas will be stored
#elements added to the collection will be retrieved from files in the /payloads directory 
collection = client.create_collection(name="payloads")

#create a list of all files containing valuable payloads
#in each file, each row represents a payload to exploit a specific vulnerability
files = os.listdir("payloads/payloadFiles/")
print(f'{files}\n\n')

countFiles=1
for f in files:
	print(f'Adding payloads contained in file {f} to the vectorDB...')
	timeBegin = time.time()
	docu = []
	meta = []
	identifiers = []
	fileName = 'payloads/payloadFiles/'+f
	fTemp = open(fileName, 'r')
	count = 1
	for line in fTemp:
		docu.append(line.strip())
		meta.append({"file" : f})
		identifiers.append(f.split('-')[0].strip()+"-"+str(countFiles)+'.'+str(count))
		count+=1
	fTemp.close()
	try:
		collection.add(
			documents=docu,
			metadatas=meta,
			ids=identifiers
		)
	except:
		print(f'Something went wrong while working on file {f}')
		continue
	timeEnd = time.time()
	print(f'Payloads in file {f} correctly added to the vectorDB in {timeEnd - timeBegin} seconds')
	countFiles+=1
		
