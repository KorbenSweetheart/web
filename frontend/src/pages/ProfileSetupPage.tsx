import type { FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { SPORTS, MODES, LEVELS, saveProfile } from '../services/profile';
import type { Profile, ProfileActivity } from '../types';
import { User, Dumbbell, SlidersHorizontal, MapPin, Headphones, Users, Sparkles } from 'lucide-react';
import { getMyProfile } from '../services/users';
import './ProfileSetupPage.css';
import { useState, useEffect } from 'react';

export default function ProfileSetupPage() {
  const [name, setName] = useState('');
  const [birthDate, setBirthDate] = useState('');
  const [bio, setBio] = useState('');
  const [pictureUrl, setPictureUrl] = useState('');
  const [modeId, setModeId] = useState<number | null>(null);
  const [maxRadius, setMaxRadius] = useState(10);
  const [activities, setActivities] = useState<ProfileActivity[]>([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();
  const [isEditing, setIsEditing] = useState(false);

  // On mount, load the existing profile (if any) to pre-fill the form.
  // First-time users get an empty form; returning users get their data.
  useEffect(() => {
    getMyProfile()
      .then((data) => {
         // If there's already a name, this is an existing profile → editing
        if (data.name) setIsEditing(true);
        if (data.name) setName(data.name);
        if (data.bio) setBio(data.bio);
        if (data.picture_url) setPictureUrl(data.picture_url);
        if (data.max_radius) setMaxRadius(data.max_radius);
        if (data.interaction_mode) setModeId(data.interaction_mode);
        if (data.birth_date) setBirthDate(data.birth_date);
        if (data.activities && data.activities.length > 0) {
          // Backend gives {id, experience_level, interest_level}
          // We use {activity_id, experience, interest_level} internally, so translate:
          setActivities(
            data.activities.map((a: any) => ({
              activity_id: a.id,
              experience: a.experience_level ?? a.experience ?? 3,
              interest_level: a.interest_level ?? 3,
            }))
          );
        }
      })
      .catch(() => {
        // No profile yet or not logged in — leave the form empty.
      });
  }, []);

  // Toggle a sport on/off
  function toggleSport(sportId: number) {
    setActivities((prev) => {
      const exists = prev.find((a) => a.activity_id === sportId);
      if (exists) {
        return prev.filter((a) => a.activity_id !== sportId);
      }
      return [...prev, { activity_id: sportId, experience: 3, interest_level: 3 }];
    });
  }

  // Change the level of a sport already selected
  function setSportLevel(sportId: number, level: number) {
    setActivities((prev) =>
      prev.map((a) =>
        a.activity_id === sportId ? { ...a, experience: level } : a
      )
    );
  }

  function isSelected(sportId: number) {
    return activities.some((a) => a.activity_id === sportId);
  }

  function getLevel(sportId: number) {
    return activities.find((a) => a.activity_id === sportId)?.experience ?? 3;
  }

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError('');

    if (!name || !birthDate || !bio) {
      setError('Please fill in name, date of birth and bio.');
      return;
    }
    if (activities.length === 0) {
      setError('Pick at least one sport.');
      return;
    }
    if (modeId === null) {
      setError('Choose a training mode.');
      return;
    }

    const profile: Profile = {
      name,
      birth_date: birthDate,
      bio,
      picture_url: pictureUrl,
      interaction_mode_id: modeId,
      activities,
      max_radius: maxRadius,
    };

    setLoading(true);
    try {
      await saveProfile(profile);
      navigate(isEditing ? '/app/profile' : '/app/discover');
    } catch (err) {
      setError('Something went wrong. Please try again.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }

  const MODE_ICONS: Record<number, typeof Headphones> = {
    1: Sparkles,    // Open to anything
    2: Users,       // Social
    3: Headphones,  // Silent
  };

  return (
    <div className="profile-setup">
      {/* Sticky floating header */}
      <div className="profile-setup__header">
        <h1 className="text-title">Set up your profile</h1>
        <p className="text-body mt-xs">Tell us a bit about you so we can find your training partners.</p>
      </div>

      {/* Form card */}
      <form onSubmit={handleSubmit} className="profile-setup__card">

        {/* ===== SECTION: General info ===== */}
        <div className="flex items-center gap-sm mb-md">
            <User size={18} strokeWidth={2.5} className="text-accent" />
            <span className="text-section">General information</span>
        </div>

        <div className="form-group">
          <label className="form-label" htmlFor="name">Name</label>
          <input id="name" type="text" className="input" placeholder="Marcus K."
            value={name} onChange={(e) => setName(e.target.value)} disabled={loading} />
        </div>

        <div className="form-group">
          <label className="form-label">Profile photo</label>
          <div className="flex items-center gap-md">
            {/* Avatar preview: photo or placeholder */}
            <div className="avatar avatar-lg" style={{ overflow: 'hidden', width: '96px', height: '96px' }}>
              {pictureUrl ? (
                <img src={pictureUrl} alt="Profile preview"
                  style={{ width: '100%', height: '100%', objectFit: 'cover' }}
                  onError={(e) => { e.currentTarget.style.display = 'none'; }} />
              ) : (
                <span style={{ fontSize: '2.5rem' }}>👤</span>
              )}
            </div>

            {/* URL input */}
            <div style={{ flex: 1 }}>
              <input type="url" className="input" placeholder="Paste an image URL"
                value={pictureUrl} onChange={(e) => setPictureUrl(e.target.value)}
                disabled={loading} />
              <p className="form-helper mt-xs">Optional — leave empty to use a placeholder.</p>
            </div>
          </div>
        </div>

        <div className="form-group">
          <label className="form-label" htmlFor="birthDate">Date of birth</label>
          <input id="birthDate" type="date" className="input"
            value={birthDate} onChange={(e) => setBirthDate(e.target.value)} disabled={loading}
            max={new Date().toISOString().split('T')[0]} />
        </div>

        <div className="form-group">
          <label className="form-label" htmlFor="bio">About me</label>
          <textarea id="bio" className="textarea" placeholder="What are you looking for in a training partner?"
            value={bio} onChange={(e) => setBio(e.target.value)} disabled={loading} />
        </div>

        {/* ===== SECTION: Sports ===== */}
        <div className="flex items-center gap-sm mt-2xl mb-xs">
            <Dumbbell size={18} strokeWidth={2.5} className="text-accent" />
            <span className="text-section">Your sports</span>
        </div>

        <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 'clamp(0.5rem, 1vw, 0.75rem)' }}>
          {SPORTS.map((sport) => (
            <div key={sport.id} className="card" style={{ padding: '0.75rem' }}>
              <label className="checkbox-item">
                <input type="checkbox" checked={isSelected(sport.id)}
                  onChange={() => toggleSport(sport.id)} disabled={loading} />
                <span>{sport.name}</span>
              </label>

              {isSelected(sport.id) && (
                <div className="flex flex-wrap gap-xs mt-sm">
                  {LEVELS.map((lvl) => (
                    <button key={lvl.value} type="button"
                      className={`pill ${getLevel(sport.id) === lvl.value ? 'pill-active' : ''}`}
                      onClick={() => setSportLevel(sport.id, lvl.value)} disabled={loading}>
                      {lvl.label}
                    </button>
                  ))}
                </div>
              )}
            </div>
          ))}
        </div>

        {/* ===== SECTION: Training mode ===== */}
        <div className="flex items-center gap-sm mt-2xl mb-xs">
            <SlidersHorizontal size={18} strokeWidth={2.5} className="text-accent" />
            <span className="text-section">Training mode</span>
        </div>
        <p className="text-body mb-md">How do you like to train?</p>

        <div className="flex flex-col gap-sm">
          {MODES.map((mode) => {
            const Icon = MODE_ICONS[mode.id];
            return (
              <button key={mode.id} type="button"
                className="card"
                onClick={() => setModeId(mode.id)} disabled={loading}
                style={{
                  textAlign: 'left',
                  width: '100%',
                  cursor: 'pointer',
                  borderColor: modeId === mode.id ? 'var(--accent-hover)' : undefined,
                  borderWidth: modeId === mode.id ? '2px' : undefined,
                }}>
                <div className="flex items-center gap-md">
                  <Icon size={24} strokeWidth={2} className={modeId === mode.id ? 'text-accent' : 'text-muted'} />
                  <div>
                    <div className="text-body-strong">{mode.title}</div>
                    <div className="text-caption mt-xs">{mode.desc}</div>
                  </div>
                </div>
              </button>
            );
          })}
        </div>

        {/* ===== SECTION: Distance ===== */}
        <div className="flex items-center gap-sm mt-2xl mb-md">
            <MapPin size={18} strokeWidth={2.5} className="text-accent" />
            <span className="text-section">Distance</span>
        </div>
        <div className="form-group mt-xs">
          <label className="form-label" htmlFor="radius">How far are you willing to travel?</label>
          <p className="form-helper mb-sm">
            We'll only show you training partners within <strong>{maxRadius} km</strong> of you.
          </p>
          <input id="radius" type="range" min={1} max={100} value={maxRadius}
            onChange={(e) => setMaxRadius(Number(e.target.value))} disabled={loading}
            style={{ width: '100%' }} />
        </div>

        {error && <p className="form-error-msg mt-md">{error}</p>}

        <button type="submit" className="btn btn-primary btn-full btn-large mt-xl" disabled={loading}>
          {loading ? 'Saving...' : 'Save and continue'}
        </button>
      </form>
    </div>
  );
}