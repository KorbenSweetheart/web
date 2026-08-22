import { useState, useEffect } from 'react';
import { MessageCircle } from 'lucide-react';
import { getMyProfile } from '../services/users';
import { getChats } from '../services/chats';
import type { ChatSummary } from '../services/chats';
import './ChatsPage.css';

export default function ChatsPage() {
  const [chats, setChats] = useState<ChatSummary[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  // On mount: first find out who I am (for resolving "the other user"),
  // then load my chats. getChats needs myId to pick the other person.
  useEffect(() => {
    getMyProfile()
      .then((me) => getChats(me.id))
      .then((data) => setChats(data))
      .catch((err) => {
        console.error(err);
        setError('Could not load chats.');
      })
      .finally(() => setLoading(false));
  }, []);

  return (
    <div className="chats">
      <div className="chats__header">
        <div className="flex items-center gap-sm">
          <MessageCircle size={28} strokeWidth={2.5} className="text-accent" />
          <h1 className="text-title">Chats</h1>
        </div>
        <p className="text-body mt-xs">Your conversations.</p>
      </div>

      {loading ? (
        <p className="text-body mt-lg">Loading…</p>
      ) : error ? (
        <p className="form-error-msg mt-lg">{error}</p>
      ) : chats.length === 0 ? (
        <p className="text-body mt-lg">No conversations yet. Connect with someone and say hi!</p>
      ) : (
        <div className="chats__list">
          {chats.map((chat) => (
            <button key={chat.id} className="chat-row" onClick={() => console.log('Open chat', chat.id)}>
              <div className="chat-row__avatar">
                {chat.other_user.picture_url ? (
                  <img src={chat.other_user.picture_url} alt={chat.other_user.name}
                    onError={(e) => { e.currentTarget.style.display = 'none'; }} />
                ) : (
                  <span>{chat.other_user.name.charAt(0).toUpperCase()}</span>
                )}
              </div>
              <div className="chat-row__info">
                <span className="chat-row__name">{chat.other_user.name}</span>
              </div>
            </button>
          ))}
        </div>
      )}
    </div>
  );
}