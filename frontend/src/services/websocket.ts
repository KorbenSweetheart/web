/* ============================================
   CHAT WEBSOCKET CLIENT
   ============================================
   Single connection to Ivan's /ws endpoint.
   Auth is by session cookie (the upgrade reads
   user_id from the request context), so we just
   connect — no auth message needed.

   Wire protocol (both directions): { type, payload }.
   The server may pack several events into one frame
   separated by '\n' (see WritePump), so on receive
   we split by newline and parse each line.
   ============================================ */

import type { Message } from './chats';

// ---- Event type constants (must match the Go wire protocol) ----
export const WsEvent = {
  ChatMessage: 'chat:message',
  ChatTyping: 'chat:typing',
  ChatRead: 'chat:read',
  PresenceCheck: 'presence:check',
  PresenceBatch: 'presence:batch',
  Error: 'error',
} as const;

// Generic envelope, matching InboundEvent / OutboundEvent on the server.
interface WsEnvelope {
  type: string;
  payload: unknown;
}

// What the server sends when the other person types.
// Note the extra user_id (who is typing) — matches Ivan's TypingBroadcastPayload.
export interface TypingBroadcast {
  chat_id: number;
  user_id: number;
  is_typing: boolean;
}
type TypingHandler = (payload: TypingBroadcast) => void;

// A new-message payload (server → client) has the same shape as a REST Message.
export type IncomingMessage = Message;

type MessageHandler = (msg: Message) => void;

// Build the ws:// (or wss://) URL from the current origin so it works
// through the Vite proxy in dev and behind TLS in prod.
function wsUrl(): string {
  const proto = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  return `${proto}//${window.location.host}/ws`;
}

export class ChatSocket {
  private ws: WebSocket | null = null;
  private messageHandlers = new Set<MessageHandler>();
  private typingHandlers = new Set<TypingHandler>();
  private shouldReconnect = true;
  private reconnectDelay = 1000; // grows on repeated failures, capped below

  connect() {
    this.shouldReconnect = true;
    // Don't open a second socket if one is already open or connecting.
    if (this.ws && (this.ws.readyState === WebSocket.OPEN || this.ws.readyState === WebSocket.CONNECTING)) {
      return;
    }
    this.open();
  }

  private open() {
    const ws = new WebSocket(wsUrl());
    this.ws = ws;

    ws.onopen = () => {
      this.reconnectDelay = 1000; // reset backoff on a clean connection
    };

    ws.onmessage = (e) => {
      // The server may pack multiple events into one frame, newline-separated.
      for (const line of String(e.data).split('\n')) {
        if (!line.trim()) continue;
        this.handleFrame(line);
      }
    };

    ws.onclose = () => {
      this.ws = null;
      if (this.shouldReconnect) {
        setTimeout(() => this.open(), this.reconnectDelay);
        this.reconnectDelay = Math.min(this.reconnectDelay * 2, 15000);
      }
    };

    ws.onerror = () => {
      ws.close(); // triggers onclose → reconnect
    };
  }

  private handleFrame(raw: string) {
    let env: WsEnvelope;
    try {
      env = JSON.parse(raw);
    } catch {
      return; // ignore malformed frame
    }

    if (env.type === WsEvent.ChatMessage) {
      const msg = env.payload as Message;
      this.messageHandlers.forEach((h) => h(msg));
    }
    // Other event types (typing, read, presence) come in steps 4.
    
    if (env.type === WsEvent.ChatTyping) {
      const payload = env.payload as TypingBroadcast;
      this.typingHandlers.forEach((h) => h(payload));
    }
  }

  // Send a chat message. Server expects { chat_id, content } in the payload.
  sendMessage(chatId: number, content: string) {
    if (this.ws?.readyState !== WebSocket.OPEN) return false;
    this.ws.send(
      JSON.stringify({
        type: WsEvent.ChatMessage,
        payload: { chat_id: chatId, content },
      }),
    );
    return true;
  }

  // Tell the server I've read this chat, so it marks the messages as viewed.
  // Payload matches the wire protocol: { chat_id }.
  sendRead(chatId: number) {
    if (this.ws?.readyState !== WebSocket.OPEN) return false;
    this.ws.send(
      JSON.stringify({
        type: WsEvent.ChatRead,
        payload: { chat_id: chatId },
      }),
    );
    return true;
  }

  // Tell the server I'm typing (or stopped) in this chat.
  // Server relays it only to the other participant. Payload: { chat_id, is_typing }.
  sendTyping(chatId: number, isTyping: boolean) {
    if (this.ws?.readyState !== WebSocket.OPEN) return false;
    this.ws.send(
      JSON.stringify({
        type: WsEvent.ChatTyping,
        payload: { chat_id: chatId, is_typing: isTyping },
      }),
    );
    return true;
  }

  // Subscribe to incoming messages. Returns an unsubscribe function.
  onMessage(handler: MessageHandler): () => void {
    this.messageHandlers.add(handler);
    return () => this.messageHandlers.delete(handler);
  }

    // Subscribe to typing events from the other participant. Returns unsubscribe.
  onTyping(handler: TypingHandler): () => void {
    this.typingHandlers.add(handler);
    return () => this.typingHandlers.delete(handler);
  }

  disconnect() {
    this.shouldReconnect = false;
    this.ws?.close();
    this.ws = null;
    this.messageHandlers.clear();
    this.typingHandlers.clear();
  }
}


// Shared single connection for the whole app. Using this everywhere
// (instead of `new ChatSocket()`) guarantees only ONE socket exists,
// even when React re-mounts components in development.
let sharedSocket: ChatSocket | null = null;

export function getChatSocket(): ChatSocket {
  if (!sharedSocket) {
    sharedSocket = new ChatSocket();
    sharedSocket.connect();
  }
  return sharedSocket;
}