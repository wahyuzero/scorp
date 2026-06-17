#!/usr/bin/env python3
"""
STT (Speech-to-Text) using faster-whisper.
Usage: python3 stt.py <audio_file> [language]
Outputs transcribed text to stdout.
"""
import sys
import os

audio_file = sys.argv[1] if len(sys.argv) > 1 else ""
language = sys.argv[2] if len(sys.argv) > 2 else None  # None = auto-detect

if not audio_file or not os.path.exists(audio_file):
    print("ERROR: audio file not found", file=sys.stderr)
    sys.exit(1)

from faster_whisper import WhisperModel

# Use "base" model — good balance of speed/accuracy (~150MB)
# Runs on CPU (no GPU on VPS)
model = WhisperModel("base", device="cpu", compute_type="int8")

segments, info = model.transcribe(audio_file, language=language, beam_size=5)

text = " ".join([seg.text.strip() for seg in segments])
print(text.strip())
