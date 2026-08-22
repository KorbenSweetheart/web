import type { UserProfile } from '../services/mockUsers';
import { X, Headphones, Users, Sparkles, MapPin } from 'lucide-react';
import './ProfilePanel.css';

const MODES: Record<number, { icon: typeof Headphones; label: string }> = {
  1: { icon: Sparkles,   label: 'Open to anything' },
  2: { icon: Users,      label: 'Social' },
  3: { icon: Headphones, label: 'Silent' },
};

const EXP_LABELS: Record<number, string> = {
  1: 'Beginner', 2: 'Active Novice', 3: 'Intermediate', 4: 'Advanced', 5: 'Professional',
};

// Which context the panel is shown in — decides the footer buttons.
type PanelVariant = 'discover' | 'received' | 'connected';

interface ProfilePanelProps {
  user: UserProfile;
  onClose: () => void;
  variant?: PanelVariant;
  // All optional: each variant only wires the handlers it needs.
  onConnect?: (id: number) => void;
  onDismiss?: (id: number) => void;
  onAccept?: (id: number) => void;
  onDecline?: (id: number) => void;
  onMessage?: (id: number) => void;
  onRemove?: (id: number) => void;
}

export default function ProfilePanel({
  user,
  onClose,
  variant = 'discover',
  onConnect,
  onDismiss,
  onAccept,
  onDecline,
  onMessage,
  onRemove,
}: ProfilePanelProps) {
  const ModeIcon = MODES[user.interaction_mode].icon;
  const initials = user.name.split(' ').map((w) => w[0]).join('').slice(0, 2).toUpperCase();

  // Footer buttons per context, mirroring UserCard's variant logic.
  function renderActions() {
    switch (variant) {
      case 'received':
        return (
          <>
            <button className="btn btn-outline" onClick={() => onDecline?.(user.id)}>Decline</button>
            <button className="btn btn-primary" onClick={() => onAccept?.(user.id)}>Accept</button>
          </>
        );
      case 'connected':
        return (
          <>
            <button className="btn btn-danger" onClick={() => onRemove?.(user.id)}>Remove</button>
            <button className="btn btn-primary" onClick={() => onMessage?.(user.id)}>Message</button>
          </>
        );
      default: // 'discover'
        return (
          <>
            <button className="btn btn-outline" onClick={() => onDismiss?.(user.id)}>Dismiss</button>
            <button className="btn btn-primary" onClick={() => onConnect?.(user.id)}>Connect</button>
          </>
        );
    }
  }

  return (
    <div className="profile-panel">
      {/* Close button */}
      <button className="profile-panel__close" onClick={onClose} title="Close">
        <X size={22} strokeWidth={2.5} />
      </button>

      {/* Header: photo left, info middle, ring right */}
      <div className="profile-panel__header">
        <div className="profile-panel__avatar">
          {user.picture_url ? (
            <img src={user.picture_url} alt={user.name}
              onError={(e) => { e.currentTarget.style.display = 'none'; }} />
          ) : (
            <span>{initials}</span>
          )}
        </div>

        <div className="profile-panel__header-info">
          <div className="profile-panel__name-row">
            <h2 className="text-title">{user.name}</h2>
            <span className={`profile-panel__status ${user.is_online ? 'is-online' : ''}`} />
          </div>
          <p className="text-body">
            {user.age} · <ModeIcon size={15} strokeWidth={2.5} style={{ verticalAlign: '-2px' }} />
            {' '}{MODES[user.interaction_mode].label}
          </p>
          <p className="text-caption mt-xs">
            <MapPin size={13} strokeWidth={2.5} style={{ verticalAlign: '-2px' }} />
            {' '}Within {user.max_radius} km
          </p>
        </div>

        {/* Ring */}
        <div className="profile-panel__ring">
          <svg width="72" height="72" viewBox="0 0 72 72">
            <circle cx="36" cy="36" r="30" className="profile-panel__ring-bg" />
            <circle cx="36" cy="36" r="30" className="profile-panel__ring-fill"
              strokeDasharray={2 * Math.PI * 30}
              strokeDashoffset={(2 * Math.PI * 30) * (1 - user.match_score / 100)}
              transform="rotate(-90 36 36)" />
          </svg>
          <span className="profile-panel__ring-text">{user.match_score}</span>
        </div>
      </div>

      {/* Bio (full, not truncated) */}
      <div className="profile-panel__section">
        <p className="text-label mb-xs">About</p>
        <p className="text-body">{user.bio}</p>
      </div>

      {/* All sports */}
      <div className="profile-panel__section">
        <p className="text-label mb-sm">Sports</p>
        <div className="profile-panel__tags">
          {user.activities.map((a) => (
            <span key={a.id} className="sport-tag">
              {a.title} · {EXP_LABELS[a.experience]}
            </span>
          ))}
        </div>
      </div>

      {/* Actions (per variant) */}
      <div className="profile-panel__actions">
        {renderActions()}
      </div>
    </div>
  );
}