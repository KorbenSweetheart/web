import { useState, useEffect } from 'react';
import { Users } from 'lucide-react';
import UserCard from '../components/UserCard';
import type { UserProfile } from '../services/mockUsers';
import {
  getConnections,
  getConnectionRequests,
  respondToConnection,
  deleteConnection,
} from '../services/connections';
import './ConnectionsPage.css';
import { useNavigate, useLocation } from 'react-router-dom';
import { getMyProfile } from '../services/users';
import { openDirectChat } from '../services/chats';
import ProfilePanel from '../components/ProfilePanel';

type Tab = 'received' | 'connected';

export default function ConnectionsPage() {
  const [tab, setTab] = useState<Tab>('received');
  const [received, setReceived] = useState<UserProfile[]>([]);
  const [connected, setConnected] = useState<UserProfile[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [myId, setMyId] = useState<number | null>(null);
  const [openingChat, setOpeningChat] = useState(false);
  const navigate = useNavigate();
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const location = useLocation();

  // Load both lists on mount.
  // Load both lists on mount.
  useEffect(() => {
    Promise.all([getConnectionRequests(), getConnections(), getMyProfile()])
      .then(([requests, conns, me]) => {
        setReceived(requests);
        setConnected(conns);
        setMyId(me.id);

        // Arriving from a chat → open that person's panel (they're a connection).
        const openUserId = (location.state as { openUserId?: number } | null)?.openUserId;
        if (openUserId != null) {
          setTab('connected');
          setSelectedId(openUserId);
        }
      })
      .catch((err) => {
        console.error(err);
        setError('Could not load connections.');
      })
      .finally(() => setLoading(false));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  async function handleAccept(id: number) {
    try {
      await respondToConnection(id, 'accepted');
      // Move the person from received → connected.
      const person = received.find((u) => u.id === id);
      setReceived((prev) => prev.filter((u) => u.id !== id));
      if (person) setConnected((prev) => [...prev, person]);
    } catch (err) {
      console.error(err);
    }
  }

  async function handleDecline(id: number) {
    try {
      await respondToConnection(id, 'declined');
      setReceived((prev) => prev.filter((u) => u.id !== id));
    } catch (err) {
      console.error(err);
    }
  }

  async function handleRemove(id: number) {
    try {
      await deleteConnection(id);
      setConnected((prev) => prev.filter((u) => u.id !== id));
    } catch (err) {
      console.error(err);
    }
  }

  async function handleMessage(id: number) {
  if (myId == null || openingChat) return;   // guard: no dispares sin myId ni doble-click
  setOpeningChat(true);
  try {
    const chat = await openDirectChat(id, myId);
    navigate('/app/chats', { state: { openChatId: chat.id } });
  } catch (err) {
    console.error(err);
  } finally {
    setOpeningChat(false);
  }
}

    const lists: Record<Tab, UserProfile[]> = { received, connected };
  const current = lists[tab];
  const selectedUser = current.find((u) => u.id === selectedId) ?? null;

  return (
    <div className="connections">
      <div className="connections__header">
        <div className="flex items-center gap-sm">
          <Users size={28} strokeWidth={2.5} className="text-accent" />
          <h1 className="text-title">Connections</h1>
        </div>
        <p className="text-body mt-xs">People you've matched with.</p>
      </div>

      <div className="connections__tabs">
        <button className={`connections__tab ${tab === 'received' ? 'is-active' : ''}`}
          onClick={() => { setTab('received'); setSelectedId(null); }}>
          Received {received.length > 0 && <span className="badge">{received.length}</span>}
        </button>
        <button className={`connections__tab ${tab === 'connected' ? 'is-active' : ''}`}
          onClick={() => { setTab('connected'); setSelectedId(null); }}>
          Connected
        </button>
      </div>

      {loading ? (
        <p className="text-body mt-lg">Loading…</p>
      ) : error ? (
        <p className="form-error-msg mt-lg">{error}</p>
      ) : current.length === 0 ? (
        <p className="text-body mt-lg">Nothing here yet.</p>
      ) : (
        <div className={`discover__layout ${selectedUser ? 'has-panel' : ''}`}>
          <div className="discover__grid">
            {current.map((user) => (
              <UserCard
                key={user.id}
                user={user}
                variant={tab === 'received' ? 'received' : 'connected'}
                onAccept={handleAccept}
                onDecline={handleDecline}
                onMessage={handleMessage}
                onCancel={handleRemove}
                onClick={setSelectedId}
              />
            ))}
          </div>

          {selectedUser && (
            <div className="discover__panel">
              <ProfilePanel
                user={selectedUser}
                variant={tab}
                onClose={() => setSelectedId(null)}
                onAccept={handleAccept}
                onDecline={handleDecline}
                onMessage={handleMessage}
                onRemove={handleRemove}
              />
            </div>
          )}
        </div>
      )}
    </div>
  );
}