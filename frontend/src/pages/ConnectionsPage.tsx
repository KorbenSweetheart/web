import { useState } from 'react';
import { Users } from 'lucide-react';
import UserCard from '../components/UserCard';
import { RECEIVED, SENT, CONNECTED } from '../services/mockConnections';
import './ConnectionsPage.css';

type Tab = 'received' | 'sent' | 'connected';

export default function ConnectionsPage() {
  const [tab, setTab] = useState<Tab>('received');
  const [received, setReceived] = useState(RECEIVED);
  const [sent, setSent] = useState(SENT);
  const [connected] = useState(CONNECTED);

  function handleAccept(id: number) {
    console.log('Accept', id);
    setReceived((prev) => prev.filter((u) => u.id !== id));
    // later: move into `connected` once backend confirms
  }

  function handleDecline(id: number) {
    console.log('Decline', id);
    setReceived((prev) => prev.filter((u) => u.id !== id));
  }

  function handleCancel(id: number) {
    console.log('Cancel request', id);
    setSent((prev) => prev.filter((u) => u.id !== id));
  }

  function handleMessage(id: number) {
    console.log('Open chat with', id);
    // later: navigate(`/app/chats/${id}`)
  }

  const lists: Record<Tab, typeof received> = { received, sent, connected };
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
        <button className={`connections__tab ${tab === 'sent' ? 'is-active' : ''}`}
          onClick={() => setTab('sent')}>
          Sent
        </button>
        <button className={`connections__tab ${tab === 'connected' ? 'is-active' : ''}`}
          onClick={() => setTab('connected')}>
          Connected
        </button>
      </div>

      {current.length === 0 ? (
        <p className="text-body mt-lg">Nothing here yet.</p>
      ) : (
        <div className="connections__grid">
          {current.map((user) => (
            <UserCard
              key={user.id}
              user={user}
              variant={tab}
              onAccept={handleAccept}
              onDecline={handleDecline}
              onCancel={handleCancel}
              onMessage={handleMessage}
            />
          ))}
        </div>
      )}
    </div>
  );
}