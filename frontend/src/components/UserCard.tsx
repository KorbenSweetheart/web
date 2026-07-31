import type { UserProfile } from '../services/mockUsers';
import './UserCard.css';
import { Headphones, Users, Sparkles, ChevronRight } from 'lucide-react';

// Map interaction_mode → icon + label
const MODES: Record<number, { icon: typeof Headphones; label: string }> = {
  1: { icon: Sparkles,   label: 'Open to anything' },
  2: { icon: Users,      label: 'Social' },
  3: { icon: Headphones, label: 'Silent' },
};

// Map experience number → readable label
const EXP_LABELS: Record<number, string> = {
  1: 'Beginner',
  2: 'Active Novice',
  3: 'Intermediate',
  4: 'Advanced',
  5: 'Professional',
};

interface UserCardProps {
  user: UserProfile;
  onConnect?: (id: number) => void;
  onDismiss?: (id: number) => void;
  onClick?: (id: number) => void;
}

export default function UserCard({ user, onConnect, onDismiss, onClick }: UserCardProps) {
  // Initials for the avatar fallback
  const initials = user.name
    .split(' ')
    .map((w) => w[0])
    .join('')
    .slice(0, 2)
    .toUpperCase();

  // Ring math: how much of the circle to fill based on score
  const radius = 26;
  const circumference = 2 * Math.PI * radius;
  const offset = circumference - (user.match_score / 100) * circumference;

  return (
    <div className="user-card" onClick={() => onClick?.(user.id)} style={{ cursor: 'pointer' }}>
      {/* Top: avatar + name + score ring */}
      <div className="user-card__top">
        <div className="user-card__avatar">
          {user.picture_url ? (
            <img src={user.picture_url} alt={user.name}
              onError={(e) => { e.currentTarget.style.display = 'none'; }} />
          ) : (
            <span>{initials}</span>
          )}
        </div>

        <div className="user-card__info">
          <div className="user-card__name-row">
            <h3 className="user-card__name">{user.name}</h3>
            <span className={`user-card__status ${user.is_online ? 'is-online' : ''}`} />
          </div>
          <p className="user-card__meta">
            {(() => {
              const ModeIcon = MODES[user.interaction_mode].icon;
              return (
                <>
                  {user.age} · <ModeIcon size={14} strokeWidth={2.5} className="user-card__mode-icon" />
                  {MODES[user.interaction_mode].label}
                </>
              );
            })()}
          </p>
        </div>

        {/* Match score ring */}
        <div className="user-card__ring">
          <svg width="64" height="64" viewBox="0 0 64 64">
            <circle cx="32" cy="32" r={radius} className="user-card__ring-bg" />
            <circle cx="32" cy="32" r={radius} className="user-card__ring-fill"
              strokeDasharray={circumference} strokeDashoffset={offset}
              transform="rotate(-90 32 32)" />
          </svg>
          <span className="user-card__ring-text">{user.match_score}</span>
        </div>
      </div>

      {/* Bio */}
      <p className="user-card__bio">{user.bio}</p>

      {/* Sport tags */}
      <div className="user-card__tags">
        {user.activities.slice(0, 3).map((a) => (
          <span key={a.id} className="sport-tag">
            {a.title} · {EXP_LABELS[a.experience]}
          </span>
        ))}
        {user.activities.length > 3 && (
          <span className="sport-tag">+{user.activities.length - 3}</span>
        )}
      </div>

      {/* Actions */}
      <div className="user-card__footer">
        <button className="user-card__more"
          onClick={(e) => { e.stopPropagation(); onClick?.(user.id); }}>
          View full profile <ChevronRight size={16} strokeWidth={2.5} />
        </button>
        <div className="user-card__actions">
          <button className="btn btn-outline"
            onClick={(e) => { e.stopPropagation(); onDismiss?.(user.id); }}>
            Dismiss
          </button>
          <button className="btn btn-primary"
            onClick={(e) => { e.stopPropagation(); onConnect?.(user.id); }}>
            Connect
          </button>
        </div>
      </div>
    </div>  
  );
}