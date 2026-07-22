FROM python:3.11-slim

RUN pip install --pre -U -f https://mlc.ai/wheels mlc-llm-nightly-cu121 mlc-ai-nightly-cu121

EXPOSE 8000

ENV MLC_MODEL=HF://mlc-ai/gemma-2b-it-q4f16_1-MLC
ENV MLC_HOST=0.0.0.0
ENV MLC_PORT=8000

ENTRYPOINT python -m mlc_llm serve $MLC_MODEL --host $MLC_HOST --port $MLC_PORT --mode interactive
