docker run -d --rm -v ./ollama:/root/.ollama -p 11430:11430 --name ollama ollama/ollama
docker exec -it -d ollama ollama run mxbai-embed-text
