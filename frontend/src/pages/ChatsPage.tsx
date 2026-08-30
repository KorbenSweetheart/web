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
  // Pagination: are there older messages to fetch, and are we mid-fetch?
  const [hasMore, setHasMore] = useState(false);
  const [loadingOlder, setLoadingOlder] = useState(false);

  const location = useLocation();
  const threadRef = useRef<HTMLDivElement>(null);
  const socketRef = useRef<ChatSocket | null>(null);
  const navigate = useNavigate();
  // chatId -> does this chat have unread messages? (icon, not a count)
  const [hasUnread, setHasUnread] = useState<Record<number, boolean>>({});
  // Is the other person currently typing in the open chat?
  const [otherTyping, setOtherTyping] = useState(false);
  // userId -> is that person online right now?
  const [online, setOnline] = useState<Record<number, boolean>>({});
  // Timer that sends "stopped typing" after a pause.
  const typingTimer = useRef<ReturnType<typeof setTimeout> | null>(null);
  // Ref to prevent auto-scrolling to bottom when prepending older messages
  const isPrependingRef = useRef(false);

  useEffect(() => {
    window.dispatchEvent(new CustomEvent('chats:unread-changed'));
  }, [hasUnread]);


  const selectedChatIdRef = useRef<number | null>(selectedChatId);
  selectedChatIdRef.current = selectedChatId;

  const myIdRef = useRef<number | null>(myId);
  myIdRef.current = myId;

  // Mount: who am I → my chats → initial unread icons from message history.
  useEffect(() => {
    getMyProfile()
      .then((me) => {
        setMyId(me.id);
        myIdRef.current = me.id;
        return getChats(me.id).then((data) => ({ me, data }));
      })
      .then(async ({ me, data }) => {
        setChats(data);
        const openId = (location.state as { openChatId?: number } | null)?.openChatId;
        if (openId != null) {
          selectedChatIdRef.current = openId;
          setSelectedChatId(openId);
        }

        // For each chat, peek at its latest messages: if any message I didn't
        // send is still unviewed, that chat gets the icon.
        const entries = await Promise.all(
          data.map(async (chat) => {
            const msgs = await getChatMessages(chat.id, 0, 50);
            const unread = msgs.some((m) => m.sender_id !== me.id && !m.is_viewed);
            return [chat.id, unread] as const;
          }),
        );
        setHasUnread((prev) => {
          const map = Object.fromEntries(entries);
          const activeId = selectedChatIdRef.current ?? openId;
          if (activeId != null) {
            map[activeId] = false;
          }
          return { ...prev, ...map };
        });
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
    setOtherTyping(false); // reset when switching chats
    // Any time a chat opens (row click OR arriving from "Message"), clear its
    // icon and tell the backend it's read. Covers both entry points.
    setHasUnread((prev) => ({ ...prev, [selectedChatId]: false }));
    socketRef.current?.sendRead(selectedChatId);
    setMsgLoading(true);
    setMsgError('');
    getChatMessages(selectedChatId)
      .then((data) => {
        setMessages([...data].reverse()); // backend DESC → chronological
        // A full first page (15) means there are probably older messages.
        // A short page means we already have the whole history.
        setHasMore(data.length >= 15);
      })
      .catch((err) => {
        console.error(err);
        setMsgError('Could not load messages.');
      })
      .finally(() => setMsgLoading(false));
  }, [selectedChatId]);

  // Stick to bottom on new messages / chat switch — but NOT when we're
  // prepending older history (that path manages its own scroll).
  useEffect(() => {
    if (isPrependingRef.current) return;
    threadRef.current?.scrollTo(0, threadRef.current.scrollHeight);
  }, [messages]);

  // Ask who's online among my chat partners — on load and every 20s after.
  // (Ivan's backend answers presence:check but doesn't push live changes,
  //  so we poll gently to keep the dots reasonably fresh.)
  useEffect(() => {
    if (chats.length === 0) return;
    const socket = socketRef.current;
    if (!socket) return;

    const ids = chats.map((c) => c.other_user.id);
    const ask = () => socket.checkPresence(ids);

    ask(); // immediate first check
    const interval = setInterval(ask, 20000);
    return () => clearInterval(interval);
  }, [chats]);

  // Subscribe to incoming messages on the shared socket.
  useEffect(() => {
    const socket = getChatSocket();
    socketRef.current = socket;

    const off = socket.onMessage((msg) => {
      const isOpenChat = msg.chat_id === selectedChatIdRef.current;
      const isMine = msg.sender_id === myIdRef.current;

      // 1. If it's the chat we're viewing, append it to the thread and acknowledge read.
      if (isOpenChat) {
        setMessages((prev) => {
          if (prev.some((m) => m.id === msg.id)) return prev; // guard against dupes
          return [...prev, msg];
        });
        // We're looking at this chat, so the message is seen immediately.
        // Tell the backend to mark it read — otherwise it stays is_viewed=false.
        if (!isMine) {
          socket.sendRead(msg.chat_id);
          setHasUnread((prev) => ({ ...prev, [msg.chat_id]: false }));
        }
      }
      // 2. Flag the chat as unread if the message is from the other person
      //    AND we're not currently looking at that chat.
      if (!isMine && !isOpenChat) {
        setHasUnread((prev) => ({ ...prev, [msg.chat_id]: true }));
      }

      // 3. Move the chat to the top of the list (most recent first).
      setChats((prev) => {
        const idx = prev.findIndex((c) => c.id === msg.chat_id);
        if (idx === 0) return prev; // already on top
        if (idx > 0) {
          const next = [...prev];
          const [moved] = next.splice(idx, 1);
          next.unshift(moved);
          return next;
        }
        // If not in our list yet (e.g. new chat initiated), refresh list from server
        if (myIdRef.current != null) {
          getChats(myIdRef.current).then((updated) => {
            setChats(updated);
          }).catch(console.error);
        }
        return prev;
      });
    });

    // Subscribe to the other person's typing state.
    const offTyping = socket.onTyping((payload) => {
      // Only react if it's the chat we're viewing AND it's not my own echo.
      if (payload.chat_id === selectedChatIdRef.current && payload.user_id !== myIdRef.current) {
        setOtherTyping(payload.is_typing);
      }
    });

    // Subscribe to presence replies — merge the batch into our online map.
    const offPresence = socket.onPresence((payload) => {
      setOnline((prev) => ({ ...prev, ...payload.statuses }));
    });

    // Only unsubscribe on unmount — keep the shared socket alive.
    return () => { off(); offTyping(); offPresence(); };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const selectedChat = chats.find((c) => c.id === selectedChatId) ?? null;

  // Load the previous page of messages when scrolling to the top or clicking 'Load earlier messages'.
  // Preserves scroll position so the view doesn't jump when older
  // messages are prepended.
  async function loadOlder() {
    if (selectedChatId == null || loadingOlder || !hasMore) return;
    if (messages.length === 0) return;

    const thread = threadRef.current;
    if (!thread) return;
    const prevHeight = thread.scrollHeight;

    // The oldest message we currently have is the cursor.
    const oldestId = messages[0].id;

    setLoadingOlder(true);
    isPrependingRef.current = true;
    try {
      const older = await getChatMessages(selectedChatId, oldestId);
      if (older.length === 0) {
        setHasMore(false);
        isPrependingRef.current = false;
        return;
      }
      // Backend is DESC → reverse to chronological, then prepend.
      setMessages((prev) => [...[...older].reverse(), ...prev]);
      setHasMore(older.length >= 15);

      // Restore scroll: keep the user looking at the same message.
      requestAnimationFrame(() => {
        if (thread) {
          thread.scrollTop = thread.scrollHeight - prevHeight;
        }
        requestAnimationFrame(() => {
          isPrependingRef.current = false;
        });
      });
    } catch (err) {
      console.error('Could not load older messages:', err);
      isPrependingRef.current = false;
    } finally {
      setLoadingOlder(false);
    }
  }

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
                  {online[chat.other_user.id] && (
                    <span className="chat-row__online" aria-label="Online" />
                  )}
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

            <div
              className="convo-thread"
              ref={threadRef}
              onScroll={(e) => {
                // Near the top → pull the previous page.
                if (e.currentTarget.scrollTop < 80) loadOlder();
              }}
            >
              {hasMore && !loadingOlder && (
                <div className="convo-load-older">
                  <button
                    type="button"
                    className="convo-load-more-btn"
                    onClick={loadOlder}
                  >
                    Load more
                  </button>
                </div>
              )}
              {loadingOlder && (
                <div className="convo-load-older">
                  <p className="text-dim text-center" style={{ margin: 0, fontSize: '0.8125rem' }}>
                    Loading…
                  </p>
                </div>
              )}
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

            {otherTyping && (
              <div className="typing-indicator" style={{ margin: '0 1rem 0.5rem' }}>
                <span className="typing-dot" />
                <span className="typing-dot" />
                <span className="typing-dot" />
              </div>
            )}

            <div className="chat-input-wrapper">
              <input
                className="chat-input"
                placeholder="Type a message…"
                value={draft}
                onChange={(e) => {
                  setDraft(e.target.value);
                  if (selectedChatId == null) return;
                  // Tell the other person I'm typing.
                  socketRef.current?.sendTyping(selectedChatId, true);
                  // Reset the 3s "stopped typing" timer on every keystroke.
                  if (typingTimer.current) clearTimeout(typingTimer.current);
                  typingTimer.current = setTimeout(() => {
                    socketRef.current?.sendTyping(selectedChatId, false);
                  }, 3000);
                }}
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