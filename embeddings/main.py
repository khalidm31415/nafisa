import os
import tensorflow_hub as hub

from dotenv import load_dotenv
from fastapi import FastAPI

config = load_dotenv()

model = os.environ.get('EMBEDDINGS_MODEL', 'https://tfhub.dev/google/universal-sentence-encoder/4')
print(f'Model used: {model}')

print("Loading pre-trained embeddings from tensorflow hub...")
use_model = hub.load(model)
print("Done.")

app = FastAPI()

@app.get('/ping')
def ping():
    return 'pong!'

@app.get('/api/embeddings')
def embeddings(text: str):
    embeddings = use_model([text])
    return embeddings.numpy().tolist()[0]
