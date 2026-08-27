import { useState, useEffect, useRef } from 'react';
import { MessageCircle, ArrowLeft, Send, MailPlus } from 'lucide-react';
import { useLocation, useNavigate } from 'react-router-dom';
import { getMyProfile } from '../services/users';
import { getChats, getChatMessages } from '../services/chats';
import type { ChatSummary, Message } from '../services/chats';
import { getChatSocket, type ChatSocket } from '../services/websocket';

import './ChatsPage.css';

// HH:MM from an ISO timestamp
function fmtTime(iso: string) {
  return new Date(iso).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

// Today / Yesterday / locale date, for day separators
function dayLabel(iso: string) {
  const d = new Date(iso);
  const now = new Date();
  const sameDay = (a: Date, b: Date) => a.toDateString() === b.toDateString();
  if (sameDay(d, now)) return 'Today';
  const y = new Date(now);
  y.setDate(now.getDate() - 1);
  if (sameDay(d, y)) return 'Yesterday';
  return d.toLocaleDateString();
}

export default function ChatsPage() {
  const [chats, setChats] = useState<ChatSummary[]>([]);
  const [myId, setMyId] = useState<number | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const [selectedChatId, setSelectedChatId] = useState<number | null>(null);
  const [messages, setMessages] = useState<Message[]>([]);
  const [msgLoading, setMsgLoading] = useState(false);
  const [msgError, setMsgError] = useState('');
  const [draft, setDraft] = useState('');

  const location = useLocation();
  const threadRef = useRef<HTMLDivElement>(null);
  const socketRef = useRef<ChatSocket | null>(null);
  const navigate = useNavigate();
  // chatId -> does this chat have unread messages? (icon, not a count)
  const [hasUnread, setHasUnread] = useState<Record<number, boolean>>({});

  // Mount: who am I → my chats → initial unread icons from message history.
  useEffect(() => {
    getMyProfile()
      .then((me) => {
        setMyId(me.id);
        return getChats(me.id).then((data) => ({ me, data }));
      })
      .then(async ({ me, data }) => {
        setChats(data);
        const openId = (location.state as { openChatId?: number } | null)?.openChatId;
        if (openId != null) setSelectedChatId(openId);

        // For each chat, peek at its latest messages: if any message I didn't
        // send is still unviewed, that chat gets the icon. This is the reliable
        // source — the WebSocket only flips it on/off live on top of this.
        const entries = await Promise.all(
          data.map(async (chat) => {
            const msgs = await getChatMessages(chat.id, 0, 15);
            const unread = msgs.some((m) => m.sender_id !== me.id && !m.is_viewed);
            return [chat.id, unread] as const;
          }),
        );
        setHasUnread(Object.fromEntries(entries));
      })
      .catch((err) => {
        console.error(err);
        setError('Could not load chats.');
      })
      .finally(() => setLoading(false));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  // Load messages when a chat is selected (first page, no pagination yet).
  useEffect(() => {
    if (selectedChatId == null) return;
    // Any time a chat opens (row click OR arriving from "Message"), clear its
    // icon and tell the backend it's read. Covers both entry points.
    setHasUnread((prev) => ({ ...prev, [selectedChatId]: false }));
    socketRef.current?.sendRead(selectedChatId);
    setMsgLoading(true);
    setMsgError('');
    getChatMessages(selectedChatId)
      .then((data) => setMessages([...data].reverse())) // backend DESC → chronological
      .catch((err) => {
        console.error(err);
        setMsgError('Could not load messages.');
      })
      .finally(() => setMsgLoading(false));
  }, [selectedChatId]);

  // Stick to bottom on new messages / chat switch.
  useEffect(() => {
    threadRef.current?.scrollTo(0, threadRef.current.scrollHeight);
  }, [messages]);

  // Subscribe to incoming messages on the shared socket.
  useEffect(() => {
    const socket = getChatSocket();
    socketRef.current = socket;

    const off = socket.onMessage((msg) => {
      const isOpenChat = msg.chat_id === selectedChatIdRef.current;
      const isMine = msg.sender_id === myIdRef.current;

      // 1. If it's the chat we're viewing, append it to the thread.
      if (isOpenChat) {
        setMessages((prev) => {
          if (prev.some((m) => m.id === msg.id)) return prev; // guard against dupes
          return [...prev, msg];
        });
      }

      // 2. Flag the chat as unread if the message is from the other person
      //    AND we're not currently looking at that chat.
      if (!isMine && !isOpenChat) {
        setHasUnread((prev) => ({ ...prev, [msg.chat_id]: true }));
      }

      // 3. Move the chat to the top of the list (most recent first).
      setChats((prev) => {
        const idx = prev.findIndex((c) => c.id === msg.chat_id);
        if (idx <= 0) return prev; // already on top, or not in our list
        const next = [...prev];
        const [moved] = next.splice(idx, 1);
        next.unshift(moved);
        return next;
      });
    });

    // Only unsubscribe on unmount — keep the shared socket alive.
    return () => { off(); };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const selectedChatIdRef = useRef<number | null>(null);
  useEffect(() => {
    selectedChatIdRef.current = selectedChatId;
  }, [selectedChatId]);

  const myIdRef = useRef<number | null>(null);
  useEffect(() => {
    myIdRef.current = myId;
  }, [myId]);

  const selectedChat = chats.find((c) => c.id === selectedChatId) ?? null;

  function handleSend() {
    const text = draft.trim();
    if (!text || selectedChatId == null) return;
    const sent = socketRef.current?.sendMessage(selectedChatId, text);
    if (sent) setDraft(''); // clear only if it actually went out
    // The message will appear when the server echoes it back via onMessage.
  }

  return (
    <div className="chats">
      {/* Header card — same style as Discover/Connections */}
      <div className="chats__header">
        <div className="flex items-center gap-sm">
          <MessageCircle size={28} strokeWidth={2.5} className="text-accent" />
          <h1 className="text-title">Messages</h1>
        </div>
        <p className="text-body mt-xs">Your conversations.</p>
      </div>

      <div className={`chats-layout ${selectedChatId != null ? 'has-selection' : ''}`}>
        {/* ---- LEFT: list ---- */}
        <div className="chats-list-pane">
          {loading ? (
            <p className="text-body mt-lg">Loading…</p>
          ) : error ? (
            <p className="form-error-msg mt-lg">{error}</p>
          ) : chats.length === 0 ? (
            <p className="text-body mt-lg">No conversations yet. Connect with someone and say hi!</p>
          ) : (
            <div className="chats__list">
              {chats.map((chat) => (
                <button
                  key={chat.id}
                  className={`chat-row ${chat.id === selectedChatId ? 'is-active' : ''}`}
                  onClick={() => {
                    setSelectedChatId(chat.id);
                    setHasUnread((prev) => ({ ...prev, [chat.id]: false }));
                    socketRef.current?.sendRead(chat.id); // tell backend it's read
                  }}
                >
                  <div className="chat-row__avatar-wrap">
                    <div className="chat-row__avatar avatar avatar-md">
                      {chat.other_user.picture_url ? (
                        <img
                          src={chat.other_user.picture_url}
                          alt={chat.other_user.name}
                          onError={(e) => { e.currentTarget.style.display = 'none'; }}
                        />
                      ) : (
                        <span>{chat.other_user.name.charAt(0).toUpperCase()}</span>
                      )}
                    </div>
                    {hasUnread[chat.id] && (
                      <span className="chat-row__badge" aria-label="Unread messages">
                        <MailPlus size={14} strokeWidth={2.5} />
                      </span>
                    )}
                  </div>
                  <span className="chat-row__name">{chat.other_user.name}</span>
                </button>
              ))}
            </div>
          )}
        </div>

        {/* ---- RIGHT: conversation (only if a chat is open) ---- */}
        {selectedChat != null && (
          <div className="chats-convo-pane">
            <div className="convo-header">
              <button
                className="convo-back"
                onClick={() => setSelectedChatId(null)}
                aria-label="Back"
              >
                <ArrowLeft size={22} />
              </button>
              {/* Tapping the person opens their panel in Connections */}
              <button
                className="convo-header__person"
                onClick={() => navigate('/app/connections', { state: { openUserId: selectedChat.other_user.id } })}
                title="View profile"
              >
                <div className="avatar avatar-md">
                  {selectedChat.other_user.picture_url ? (
                    <img src={selectedChat.other_user.picture_url} alt={selectedChat.other_user.name} />
                  ) : (
                    <span>{selectedChat.other_user.name.charAt(0).toUpperCase()}</span>
                  )}
                </div>
                <span className="text-body-strong">{selectedChat.other_user.name}</span>
              </button>
            </div>

            <div className="convo-thread" ref={threadRef}>
              {msgLoading ? (
                <p className="text-body">Loading…</p>
              ) : msgError ? (
                <p className="form-error-msg">{msgError}</p>
              ) : messages.length === 0 ? (
                <p className="text-dim text-center mt-lg">No messages yet. Say hi!</p>
              ) : (
                messages.map((m, i) => {
                  const mine = m.sender_id === myId;
                  const prev = messages[i - 1];
                  const showDay = !prev || dayLabel(prev.created_at) !== dayLabel(m.created_at);
                  return (
                    <div key={m.id}>
                      {showDay && (
                        <div className="convo-day"><span>{dayLabel(m.created_at)}</span></div>
                      )}
                      <div className={`message ${mine ? 'message-sent' : 'message-received'}`}>
                        {m.content}
                      </div>
                      <div className={`message-time ${mine ? 'message-time-sent' : ''}`}>
                        {fmtTime(m.created_at)}
                      </div>
                    </div>
                  );
                })
              )}
            </div>

            <div className="chat-input-wrapper">
              <input
                className="chat-input"
                placeholder="Type a message…"
                value={draft}
                onChange={(e) => setDraft(e.target.value)}
                onKeyDown={(e) => { if (e.key === 'Enter') handleSend(); }}
              />
              <button className="btn-icon" onClick={handleSend} aria-label="Send">
                <Send size={18} />
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}