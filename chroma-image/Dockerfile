FROM python:latest
WORKDIR /app
COPY requirements.txt /app/
#RUN apt-get update && apt-get install -y python3 && apt-get install -y python3-pip
RUN pip install -r requirements.txt
RUN mkdir /app/vectorDB

#COPY /files/* /app/vectorDB
#COPY restoreDb.py /app/

#RUN python3 restoreDb.py

EXPOSE 8000

CMD chroma run # --path /vectorDB
