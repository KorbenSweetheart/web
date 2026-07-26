import { useState } from 'react';
import type { FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { SPORTS, MODES, LEVELS, saveProfile } from '../services/profile';
import type { Profile, ProfileActivity } from '../types';

export default function ProfileSetupPage() {
  const [name, setName] = useState('');
  const [age, setAge] = useState('');
  const [bio, setBio] = useState('');
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

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError('');

    if (!name || !age || !bio) {
      setError('Please fill in name, age and bio.');
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
      age: Number(age),
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

  return (
    <div className="container" style={{ maxWidth: 640, paddingTop: '2rem', paddingBottom: '2rem' }}>
      <h1 className="text-title mb-sm">Set up your profile</h1>
      <p className="text-body mt-xs">Tell us a bit about you so we can find your training partners.</p>

      <form onSubmit={handleSubmit} className="mt-lg">
        {/* Name */}
        <div className="form-group">
          <label className="form-label" htmlFor="name">Name</label>
          <input id="name" type="text" className="input" placeholder="Marcus K."
            value={name} onChange={(e) => setName(e.target.value)} disabled={loading} />
        </div>

        {/* Age */}
        <div className="form-group">
          <label className="form-label" htmlFor="age">Age</label>
          <input id="age" type="number" className="input" placeholder="28"
            value={age} onChange={(e) => setAge(e.target.value)} disabled={loading} />
        </div>

        {/* Bio */}
        <div className="form-group">
          <label className="form-label" htmlFor="bio">About me</label>
          <textarea id="bio" className="textarea" placeholder="What are you looking for in a training partner?"
            value={bio} onChange={(e) => setBio(e.target.value)} disabled={loading} />
        </div>

        {/* Sports + level each */}
        <div className="form-group">
          <label className="form-label">Sports &amp; your level</label>
          <p className="form-helper mb-sm">Pick the sports you do. Set your level for each.</p>
          <div style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: 'clamp(0.5rem, 1vw, 0.75rem)' }}>
            {SPORTS.map((sport) => (
              <div key={sport.id} className="card" style={{ padding: '0.75rem' }}>
                <label className="checkbox-item">
                  <input type="checkbox" checked={isSelected(sport.id)}
                    onChange={() => toggleSport(sport.id)} disabled={loading} />
                  <span>{sport.name}</span>
                </label>

                {/* Level selector appears only if this sport is picked */}
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
        </div>

        {/* Training mode */}
        <div className="form-group">
          <label className="form-label">Training mode</label>
          <div className="flex flex-col gap-sm">
            {MODES.map((mode) => (
              <button key={mode.id} type="button"
                className={`training-mode ${modeId === mode.id ? 'active' : ''}`}
                onClick={() => setModeId(mode.id)} disabled={loading}
                style={{ textAlign: 'left', width: '100%' }}>
                <div>
                  <div className="training-mode__label">{mode.title}</div>
                  <div className="training-mode__desc">{mode.desc}</div>
                </div>
              </button>
            ))}
          </div>
        </div>

        {/* Max radius */}
        <div className="form-group">
          <label className="form-label" htmlFor="radius">Max distance: {maxRadius} km</label>
          <input id="radius" type="range" min={1} max={100} value={maxRadius}
            onChange={(e) => setMaxRadius(Number(e.target.value))} disabled={loading}
            style={{ width: '100%' }} />
        </div>

        {error && <p className="form-error-msg">{error}</p>}

        <button type="submit" className="btn btn-primary btn-full btn-large mt-md" disabled={loading}>
          {loading ? 'Saving...' : 'Save and continue'}
        </button>
      </form>
    </div>
  );
}