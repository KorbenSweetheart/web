import { useState, useEffect } from 'react';
import { Compass } from 'lucide-react';
import UserCard from '../components/UserCard';
import ProfilePanel from '../components/ProfilePanel';
import { getRecommendations, updateMyLocation } from '../services/users';
import { sendConnectionRequest } from '../services/connections';
import type { UserProfile } from '../services/mockUsers';
import './DiscoverPage.css';

export default function DiscoverPage() {
  const [users, setUsers] = useState<UserProfile[]>([]);
  const [selectedId, setSelectedId] = useState<number | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  // On mount: send fresh location first, then fetch recommendations.
  useEffect(() => {
    updateMyLocation()
      .then(() => getRecommendations())
      .then((data) => setUsers(data))
      .catch((err) => {
        console.error(err);
        setError('Could not load recommendations.');
      })
      .finally(() => setLoading(false));
  }, []);

  const selectedUser = users.find((u) => u.id === selectedId) ?? null;

  async function handleConnect(id: number) {
      try {
      await sendConnectionRequest(id);
      // Remove them from the list once the request is sent.
        setUsers((prev) => prev.filter((u) => u.id !== id));
        setSelectedId(null);
      } catch (err) {
      console.error('Could not send connection request:', err);
    }
  }

  function handleDismiss(id: number) {
    console.log('Dismiss user', id);
    setUsers((prev) => prev.filter((u) => u.id !== id));
    setSelectedId(null);
  }

  return (
    <div className="discover">
      {/* Header card */}
      <div className="discover__header">
        <div className="flex items-center gap-sm">
          <Compass size={28} strokeWidth={2.5} className="text-accent" />
          <h1 className="text-title">Discover</h1>
        </div>
        <p className="text-body mt-xs">People who match how you like to train.</p>
      </div>

      {loading ? (
        <p className="text-body">Loading recommendations…</p>
      ) : error ? (
        <p className="form-error-msg">{error}</p>
      ) : users.length === 0 ? (
        <p className="text-body">No more suggestions right now — check back later!</p>
      ) : (
        <div className={`discover__layout ${selectedUser ? 'has-panel' : ''}`}>
          {/* List of cards */}
          <div className="discover__grid">
            {users.map((user) => (
              <UserCard
                key={user.id}
                user={user}
                onConnect={handleConnect}
                onDismiss={handleDismiss}
                onClick={setSelectedId}
              />
            ))}
          </div>

          {/* Detail panel (desktop right / mobile fullscreen) */}
          {selectedUser && (
            <div className="discover__panel">
              <ProfilePanel
                user={selectedUser}
                onClose={() => setSelectedId(null)}
                onConnect={handleConnect}
                onDismiss={handleDismiss}
              />
            </div>
          )}
        </div>
      )}
    </div>
  );
}