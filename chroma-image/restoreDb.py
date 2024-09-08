import chromadb
from chroma_datasets.utils import import_chroma_exported_hf_dataset_from_disk

cl = chromadb.HttpClient("0.0.0.0", 8000)
collection = import_chroma_exported_hf_dataset_from_disk(cl, "../exportedDB", "payloads")
