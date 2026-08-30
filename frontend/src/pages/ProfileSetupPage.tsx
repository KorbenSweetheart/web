import type { FormEvent, ChangeEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { SPORTS, MODES, LEVELS, saveProfile } from '../services/profile';
import type { Profile, ProfileActivity } from '../types';
import { User, Dumbbell, SlidersHorizontal, MapPin, Headphones, Users, Sparkles } from 'lucide-react';
import { getMyProfile, uploadProfilePicture, deleteProfilePicture } from '../services/users';
import './ProfileSetupPage.css';
import { useState, useEffect } from 'react';

export default function ProfileSetupPage() {
  const [name, setName] = useState('');
  const [age, setAge] = useState('');
  const [bio, setBio] = useState('');
  const [pictureUrl, setPictureUrl] = useState('');
  const [modeId, setModeId] = useState<number | null>(null);
  const [maxRadius, setMaxRadius] = useState(10);
  const [activities, setActivities] = useState<ProfileActivity[]>([]);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [uploadingPhoto, setUploadingPhoto] = useState(false);
  const navigate = useNavigate();
  const [isEditing, setIsEditing] = useState(false);
  const [lat, setLat] = useState<number | null>(null);
  const [lon, setLon] = useState<number | null>(null);
  const INTEREST_LEVELS = [
    { value: 1, label: 'Curious' },
    { value: 2, label: 'Casual' },
    { value: 3, label: 'Into it' },
    { value: 4, label: 'Keen' },
    { value: 5, label: 'Obsessed' },
  ];

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
        if (data.age) setAge(String(data.age));
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

  // Ask the browser for the user's location (for distance matching).
  useEffect(() => {
    if (!navigator.geolocation) return;
    navigator.geolocation.getCurrentPosition(
      (pos) => { setLat(pos.coords.latitude); setLon(pos.coords.longitude); },
      (err) => { console.warn('Location not available:', err.message); }
    );
  }, []);

  // Upload a picture as soon as the user picks a file.
  // The backend stores it and returns the final URL for the avatar preview.
  async function handlePhotoChange(e: ChangeEvent<HTMLInputElement>) {
    const file = e.target.files?.[0];
    if (!file) return;

    // Quick client-side checks (backend validates too, but this is instant feedback)
    const allowedTypes = ['image/jpeg', 'image/png'];
    if (!allowedTypes.includes(file.type)) {
      setError('Please choose a JPEG or PNG image.');
      e.target.value = '';
      return;
    }
    if (file.size > 1024 * 1024) {
      setError('Image must be 1MB or smaller.');
      e.target.value = '';
      return;
    }

    setError('');
    setUploadingPhoto(true);
    try {
      const url = await uploadProfilePicture(file);
      setPictureUrl(url);
    } catch (err) {
      setError('Could not upload the photo. Please try again.');
      console.error(err);
    } finally {
      setUploadingPhoto(false);
      e.target.value = ''; // reset so the same file can be picked again
    }
  }

  // Remove the current picture and fall back to the placeholder.
  async function handleRemovePhoto() {
    try {
      await deleteProfilePicture();
      setPictureUrl('');
    } catch (err) {
      setError('Could not remove the photo. Please try again.');
      console.error(err);
    }
  }

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

  // Change the interest level of a sport already selected
  function setSportInterest(sportId: number, interest: number) {
    setActivities((prev) =>
      prev.map((a) =>
        a.activity_id === sportId ? { ...a, interest_level: interest } : a
      )
    );
  }

  function getInterest(sportId: number) {
    return activities.find((a) => a.activity_id === sportId)?.interest_level ?? 3;
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

    if (!name) {
      setError('Please enter your name.');
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
      picture_url: pictureUrl,
      interaction_mode_id: modeId,
      activities,
      max_radius: maxRadius,
      lat: lat ?? 0,
      lon: lon ?? 0,
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
        <div className="flex items-baseline gap-sm mb-md">
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
          <div className="flex items-baseline gap-md">
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

            {/* Upload controls */}
            <div style={{ flex: 1 }}>
              {/* The label acts as a button; the real file input is hidden inside it. */}
              <label className="btn btn-outline"
                style={{
                  pointerEvents: (uploadingPhoto || loading) ? 'none' : 'auto',
                  opacity: (uploadingPhoto || loading) ? 0.6 : 1,
                }}>
                {uploadingPhoto ? 'Uploading...' : (pictureUrl ? 'Change photo' : 'Upload photo')}
                <input type="file" accept="image/jpeg,image/png" className="visually-hidden"
                  onChange={handlePhotoChange} disabled={uploadingPhoto || loading} />
              </label>

              {/* Remove button: only when there's a photo */}
              {pictureUrl && !uploadingPhoto && (
                <button type="button" className="btn btn-ghost btn-small mt-xs"
                  onClick={handleRemovePhoto} disabled={loading}>
                  Remove photo
                </button>
              )}

              <p className="form-helper mt-xs">JPEG or PNG, max 1MB. Optional — a placeholder is used if empty.</p>
            </div>
          </div>
        </div>

        <div className="form-group">
          <label className="form-label" htmlFor="age">Age</label>
          <input id="age" type="number" className="input" placeholder="28"
            value={age} onChange={(e) => setAge(e.target.value)} disabled={loading}
            min={16} max={120} />
        </div>

        <div className="form-group">
          <label className="form-label" htmlFor="bio">About me</label>
          <textarea id="bio" className="textarea" placeholder="What are you looking for in a training partner?"
            value={bio} onChange={(e) => setBio(e.target.value)} disabled={loading} />
        </div>

        {/* ===== SECTION: Sports ===== */}
        <div className="flex items-baseline gap-sm mt-2xl mb-xs">
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
                <>
                  <p className="form-helper mt-sm mb-xs">Level</p>
                  <div className="flex flex-wrap gap-xs">
                    {LEVELS.map((lvl) => (
                      <button key={lvl.value} type="button"
                        className={`pill ${getLevel(sport.id) === lvl.value ? 'pill-active' : ''}`}
                        onClick={() => setSportLevel(sport.id, lvl.value)} disabled={loading}>
                        {lvl.label}
                      </button>
                    ))}
                  </div>

                  <p className="form-helper mt-sm mb-xs">How into it are you?</p>
                  <div className="flex flex-wrap gap-xs">
                    {INTEREST_LEVELS.map((lvl) => (
                      <button key={lvl.value} type="button"
                        className={`pill ${getInterest(sport.id) === lvl.value ? 'pill-active' : ''}`}
                        onClick={() => setSportInterest(sport.id, lvl.value)} disabled={loading}>
                        {lvl.label}
                      </button>
                    ))}
                  </div>
                </>
              )}
            </div>
          ))}
        </div>

        {/* ===== SECTION: Training mode ===== */}
        <div className="flex items-baseline gap-sm mt-2xl mb-xs">
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
                <div className="flex items-baseline gap-md">
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
        <div className="flex items-baseline gap-sm mt-2xl mb-md">
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