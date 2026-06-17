#!/usr/bin/env python3
"""
RAG Embedding Service - uses sentence-transformers for local embeddings
Run as: python3 rag_embed.py <command> [args]
Commands:
  embed <text>          - Generate embedding for text
  embed_file <path>     - Generate embedding for file content
  similarity <vec1> <vec2> - Compute cosine similarity
"""

import sys
import json
import os
import numpy as np

# Try to import sentence-transformers
try:
    from sentence_transformers import SentenceTransformer
except ImportError:
    print(json.dumps({"error": "sentence-transformers not installed. Run: pip install sentence-transformers"}))
    sys.exit(1)

# Model name (small, fast, good quality)
MODEL_NAME = "sentence-transformers/all-MiniLM-L6-v2"
MODEL = None

def get_model():
    global MODEL
    if MODEL is None:
        MODEL = SentenceTransformer(MODEL_NAME)
    return MODEL

def embed_text(text):
    model = get_model()
    embedding = model.encode(text, convert_to_numpy=True, normalize_embeddings=True)
    return embedding.tolist()

def embed_file(path):
    if not os.path.exists(path):
        return {"error": f"File not found: {path}"}
    
    with open(path, 'r', encoding='utf-8', errors='ignore') as f:
        content = f.read()
    
    # Chunk if too long (max 512 tokens ~ 2000 chars for MiniLM)
    max_chunk = 2000
    if len(content) <= max_chunk:
        return {"embedding": embed_text(content), "chunks": 1}
    
    # Split into chunks
    chunks = [content[i:i+max_chunk] for i in range(0, len(content), max_chunk)]
    embeddings = [embed_text(chunk) for chunk in chunks]
    # Average embeddings
    avg_emb = np.mean(embeddings, axis=0).tolist()
    return {"embedding": avg_emb, "chunks": len(chunks)}

def cosine_similarity(vec1, vec2):
    v1 = np.array(vec1)
    v2 = np.array(vec2)
    return float(np.dot(v1, v2) / (np.linalg.norm(v1) * np.linalg.norm(v2)))

if __name__ == "__main__":
    if len(sys.argv) < 2:
        print(json.dumps({"error": "Usage: rag_embed.py <embed|embed_file|similarity> [args]"}))
        sys.exit(1)
    
    cmd = sys.argv[1]
    
    if cmd == "embed":
        if len(sys.argv) < 3:
            print(json.dumps({"error": "embed requires text argument"}))
            sys.exit(1)
        text = " ".join(sys.argv[2:])
        emb = embed_text(text)
        print(json.dumps({"embedding": emb}))
    
    elif cmd == "embed_file":
        if len(sys.argv) < 3:
            print(json.dumps({"error": "embed_file requires path argument"}))
            sys.exit(1)
        result = embed_file(sys.argv[2])
        print(json.dumps(result))
    
    elif cmd == "similarity":
        if len(sys.argv) < 4:
            print(json.dumps({"error": "similarity requires two vectors"}))
            sys.exit(1)
        vec1 = json.loads(sys.argv[2])
        vec2 = json.loads(sys.argv[3])
        sim = cosine_similarity(vec1, vec2)
        print(json.dumps({"similarity": sim}))
    
    else:
        print(json.dumps({"error": f"Unknown command: {cmd}"}))
        sys.exit(1)