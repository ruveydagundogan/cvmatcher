#!/bin/bash
# Start the WebSocket tunnel agent to connect Render backend with local Ollama
# Usage: ./start-tunnel.sh [server-url]

SERVER_URL="${1:-wss://llm-decision-score-api.onrender.com/tunnel/ws}"
OLLAMA_URL="${2:-http://localhost:11434}"

echo "Starting tunnel agent..."
echo "  Server: $SERVER_URL"
echo "  Ollama: $OLLAMA_URL"

exec ./tunnel-agent -server "$SERVER_URL" -ollama "$OLLAMA_URL"
