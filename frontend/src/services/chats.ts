/* ============================================
   CHAT SERVICE
   ============================================
   Talks to Ivan's REST chat endpoints and shapes
   the data for the UI.

   The backend sends each chat with user_one +
   user_two, but doesn't say which one is "me".
   So we resolve "the other person" here (using my
   own id) and hand the pages a clean ChatSummary
   with other_user ready to render.

   WebSocket (real-time) lives in a separate file —
   this one is REST only (open chat, list, history).
   ============================================ */

import { apiGet, apiPost } from './api';

// ---- Types ----

// A minimal user, as the chat endpoints return it.
export interface ChatUser {
  id: number;
  name: string;
  picture_url: string;
}

// Raw chat exactly as the backend sends it.
// We don't expose this to the pages — we convert it to ChatSummary.
interface RawChat {
  id: number;
  user_one_id: number;
  user_two_id: number;
  created_at: string;
  user_one: ChatUser;
  user_two: ChatUser;
}

// Clean chat for the UI: the id, when it started, and WHO the other
// person is (already resolved). This is what the chat list renders.
export interface ChatSummary {
  id: number;
  created_at: string;
  other_user: ChatUser;
}

// A single message, exactly as the backend sends it.
export interface Message {
  id: number;
  chat_id: number;
  sender_id: number;
  content: string;
  created_at: string;
  is_viewed: boolean;
}

// ---- Helper ----

// Given a raw chat and my own id, figure out who the OTHER person is.
// If I'm user_one, the other is user_two, and vice versa.
function toChatSummary(raw: RawChat, myId: number): ChatSummary {
  const other = raw.user_one.id === myId ? raw.user_two : raw.user_one;
  return {
    id: raw.id,
    created_at: raw.created_at,
    other_user: other,
  };
}

// ---- REST functions ----

// POST /chats/direct → open (or get) the 1-on-1 chat with a user.
// Used when I click "Message" on a connection. Backend returns 400 if
// we don't have an accepted connection (handled by the caller).
export async function openDirectChat(targetUserId: number, myId: number): Promise<ChatSummary> {
  const raw: RawChat = await apiPost('/chats/direct', { target_user_id: targetUserId });
  return toChatSummary(raw, myId);
}

// GET /chats → all my chats, each with the other person resolved.
// This is the chat list (left column).
export async function getChats(myId: number): Promise<ChatSummary[]> {
  const raw: RawChat[] = await apiGet('/chats');
  return (raw ?? []).map((chat) => toChatSummary(chat, myId));
}

// GET /chats/:id/messages → paginated history for one chat.
// Newest first (backend orders by id DESC). last_message_id is the
// cursor: pass the oldest id you already have to fetch older ones.
// Omit it (or 0) for the first page.
export async function getChatMessages(
  chatId: number,
  lastMessageId = 0,
  limit = 15,
): Promise<Message[]> {
  let path = `/chats/${chatId}/messages?limit=${limit}`;
  if (lastMessageId > 0) {
    path += `&last_message_id=${lastMessageId}`;
  }
  return apiGet(path);
}