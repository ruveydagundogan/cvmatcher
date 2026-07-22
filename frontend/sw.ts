import { ServiceWorkerMLCEngineHandler } from "@mlc-ai/web-llm";

let handler: ServiceWorkerMLCEngineHandler;

self.addEventListener("activate", (event) => {
  event.waitUntil(self.clients.claim());
  handler = new ServiceWorkerMLCEngineHandler();
});
