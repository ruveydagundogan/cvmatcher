import { api } from "@/lib/api";

export interface MCPMessage {
  role: "user" | "assistant" | "system";
  content: string;
}

export interface MCPRequest {
  model?: string;
  adapter?: string;
  messages: MCPMessage[];
  system_prompt?: string;
  max_tokens?: number;
  temperature?: number;
  top_p?: number;
}

export interface MCPResponse {
  response: {
    id: string;
    model: string;
    choices: Array<{
      index: number;
      message: { role: string; content: string };
      finish_reason: string;
    }>;
    usage: { prompt_tokens: number; completion_tokens: number; total_tokens: number };
  };
  rich_result: {
    text: string;
    type: string;
    sections: Array<{ title: string; content: string; type: string }>;
    metadata?: Record<string, unknown>;
  };
}

export async function mcpQuery(req: MCPRequest): Promise<MCPResponse> {
  return api.post("/api/v1/mcp/query", req);
}

export async function getMCPAdapters(): Promise<Array<{ name: string; description: string; active: boolean }>> {
  return api.get("/api/v1/mcp/adapters");
}
