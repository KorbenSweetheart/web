import { useState, useEffect } from 'react';
import { getMyProfile } from '../services/users';
import type { UserProfile } from '../services/mockUsers';
import { User, Headphones, Users, Sparkles, MapPin, Pencil } from 'lucide-react';
import { EXP_LABELS, INTEREST_LABELS } from '../services/labels';
import './ProfilePage.css';
import { useNavigate } from 'react-router-dom';

const MODES: Record<number, { icon: typeof Headphones; label: string }> = {
  1: { icon: Sparkles,   label: 'Open to anything' },
  2: { icon: Users,      label: 'Social' },
  3: { icon: Headphones, label: 'Silent' },
};

export default function ProfilePage() {
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    getMyProfile()
      .then((data) => setProfile(data))
      .catch((err) => {
        console.error(err);
        setError('Could not load your profile.');
      })
      .finally(() => setLoading(false));
  }, []);

  if (loading) {
    return (
      <div className="profile-page">
        <div className="profile-page__card">
          <p className="text-body">Loading your profile…</p>
        </div>
      </div>
    );
  }

  if (error || !profile) {
    return (
      <div className="profile-page">
        <div className="profile-page__card">
          <p className="form-error-msg">{error || 'No profile found.'}</p>
        </div>
      </div>
    );
  }

  const ModeIcon = MODES[profile.interaction_mode]?.icon ?? Sparkles;
  const hasProfile = profile.activities.length > 0 || profile.bio;

  return (
    <div className="profile-page">
      {/* Header card — same style as setup/discover */}
      <div className="profile-page__header">
        <div className="flex items-center gap-sm">
          <User size={28} strokeWidth={2.5} className="text-accent" />
          <h1 className="text-title">My profile</h1>
        </div>
        <button className="btn btn-outline" onClick={() => navigate('/app/profile-setup')}>
          <Pencil size={16} strokeWidth={2.5} /> Edit
        </button>
      </div>

      {/* Content card */}
      <div className="profile-page__card">
        {/* Photo + name */}
        <div className="flex items-center gap-md mb-xl">
          <div className="avatar avatar-lg" style={{ width: 96, height: 96, overflow: 'hidden' }}>
            {profile.picture_url ? (
              <img src={profile.picture_url} alt={profile.name}
                style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                onError={(e) => { e.currentTarget.style.display = 'none'; }} />
            ) : (
              <span style={{ fontSize: '2.5rem' }}>👤</span>
            )}
          </div>
          <div>
            <h2 className="text-subtitle">{profile.name || 'No name yet'}</h2>
            {profile.age > 0 && <p className="text-body">{profile.age} years old</p>}
          </div>
        </div>

        {!hasProfile ? (
          <div className="profile-page__empty">
            <p className="text-body mb-md">Your profile is looking a bit empty.</p>
            <p className="text-caption mb-lg">Complete your profile to add your sports, bio and training preferences — and to start seeing recommendations.</p>
            <button className="btn btn-primary" onClick={() => navigate('/app/profile-setup')}>
              Complete my profile
            </button>
          </div>
        ) : (
          <>
            {profile.bio && (
              <div className="mb-xl">
                <p className="text-section mb-xs">About</p>
                <p className="text-body">{profile.bio}</p>
              </div>
            )}

            <div className="mb-xl">
              <p className="text-section mb-xs">Training mode</p>
              <p className="text-body">
                <ModeIcon size={16} strokeWidth={2.5} style={{ verticalAlign: '-2px' }} />
                {' '}{MODES[profile.interaction_mode]?.label}
              </p>
            </div>

            {profile.activities.length > 0 && (
              <div className="mb-xl">
                <p className="text-section mb-sm">Sports</p>
                <div className="flex flex-wrap gap-sm">
                  {profile.activities.map((a) => (
                    <span key={a.id} className="sport-tag">
                      <span className="sport-tag__name">{a.title}</span>
                      <span className="sport-tag__meta"> · {EXP_LABELS[a.experience]}</span>
                      {a.interest_level != null && (
                        <span className="sport-tag__meta"> · {INTEREST_LABELS[a.interest_level]}</span>
                      )}
                    </span>
                  ))}
                </div>
              </div>
            )}

            <div>
              <p className="text-section mb-xs">Distance</p>
              <p className="text-body">
                <MapPin size={16} strokeWidth={2.5} style={{ verticalAlign: '-2px' }} />
                {' '}Within {profile.max_radius} km
              </p>
            </div>
          </>
        )}
      </div>
    </div>
  );
}