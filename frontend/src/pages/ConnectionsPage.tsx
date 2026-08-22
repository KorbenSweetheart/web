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
import { useNavigate } from 'react-router-dom';
import { getMyProfile } from '../services/users';
import { openDirectChat } from '../services/chats';

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

  // Load both lists on mount.
  useEffect(() => {
    Promise.all([getConnectionRequests(), getConnections(), getMyProfile()])
      .then(([requests, conns, me]) => {
        setReceived(requests);
        setConnected(conns);
        setMyId(me.id);
      })
      .catch((err) => {
        console.error(err);
        setError('Could not load connections.');
      })
      .finally(() => setLoading(false));
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
          onClick={() => setTab('received')}>
          Received {received.length > 0 && <span className="badge">{received.length}</span>}
        </button>
        <button className={`connections__tab ${tab === 'connected' ? 'is-active' : ''}`}
          onClick={() => setTab('connected')}>
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
        <div className="connections__grid">
          {current.map((user) => (
            <UserCard
              key={user.id}
              user={user}
              variant={tab === 'received' ? 'received' : 'connected'}
              onAccept={handleAccept}
              onDecline={handleDecline}
              onMessage={handleMessage}
              onCancel={handleRemove}
            />
          ))}
        </div>
      )}
    </div>
  );
}