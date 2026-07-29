import { useState } from 'react';
import type { FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { SPORTS, MODES, LEVELS, saveProfile } from '../services/profile';
import type { Profile, ProfileActivity } from '../types';
import { User, Dumbbell, SlidersHorizontal, MapPin, Headphones, Users, Sparkles } from 'lucide-react';

export default function ProfileSetupPage() {
  const [name, setName] = useState('');
  const [birthDate, setBirthDate] = useState('');  const [bio, setBio] = useState('');
  const [modeId, setModeId] = useState<number | null>(null);
  const [maxRadius, setMaxRadius] = useState(10);
  const [activities, setActivities] = useState<ProfileActivity[]>([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

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

// Turn a birth date into an age number (what the backend expects)
  function calculateAge(birth: string): number {
    const today = new Date();
    const b = new Date(birth);
    let age = today.getFullYear() - b.getFullYear();
    const monthDiff = today.getMonth() - b.getMonth();
    if (monthDiff < 0 || (monthDiff === 0 && today.getDate() < b.getDate())) {
      age--;
    }
    return age;
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
      age: calculateAge(birthDate),
      bio,
      picture_url: '',
      interaction_mode_id: modeId,
      activities,
      max_radius: maxRadius,
    };

    setLoading(true);
    try {
      await saveProfile(profile);
      navigate('/app/discover');
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
    <div className="container" style={{ maxWidth: 640, paddingTop: '2rem', paddingBottom: '2rem' }}>
      <h1 className="text-title">Set up your profile</h1>
      <p className="text-body mt-xs">Tell us a bit about you so we can find your training partners.</p>

      <form onSubmit={handleSubmit} className="mt-xl">

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