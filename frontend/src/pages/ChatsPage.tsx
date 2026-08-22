import { useState, useEffect, useRef } from 'react';
import { MessageCircle, ArrowLeft, Send } from 'lucide-react';
import { useLocation, useNavigate } from 'react-router-dom';
import { getMyProfile } from '../services/users';
import { getChats, getChatMessages } from '../services/chats';
import type { ChatSummary, Message } from '../services/chats';
import { ChatSocket } from '../services/websocket';

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

  // Mount: who am I → my chats. Preselect if we arrived from "Message".
  useEffect(() => {
    getMyProfile()
      .then((me) => {
        setMyId(me.id);
        return getChats(me.id);
      })
      .then((data) => {
        setChats(data);
        const openId = (location.state as { openChatId?: number } | null)?.openChatId;
        if (openId != null) setSelectedChatId(openId);
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

  // Open the WebSocket once on mount, listen for incoming messages.
useEffect(() => {
  const socket = new ChatSocket();
  socketRef.current = socket;
  socket.connect();

  const off = socket.onMessage((msg) => {
    // Only append if it belongs to the chat we're currently viewing.
    // (Other chats' messages will update unread badges later, in step 4.)
    setMessages((prev) => {
      if (msg.chat_id !== selectedChatIdRef.current) return prev;
      if (prev.some((m) => m.id === msg.id)) return prev; // guard against dupes
      return [...prev, msg];
    });
  });

  return () => {
    off();
    socket.disconnect();
  };
  // eslint-disable-next-line react-hooks/exhaustive-deps
}, []);

  const selectedChatIdRef = useRef<number | null>(null);
  useEffect(() => {
  selectedChatIdRef.current = selectedChatId;
  }, [selectedChatId]);

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
      {/* Header card — igual que Discover/Connections */}
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
                  onClick={() => setSelectedChatId(chat.id)}
                >
                  <div className="chat-row__avatar avatar avatar-sm">
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
                  <span className="chat-row__name">{chat.other_user.name}</span>
                </button>
              ))}
            </div>
          )}
        </div>

        {/* ---- RIGHT: conversation (solo si hay chat abierto) ---- */}
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