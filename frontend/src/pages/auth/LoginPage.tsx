import { useState } from 'react';
import { getMyProfile, isProfileComplete } from '../../services/users';
import type { FormEvent } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { loginUser } from '../../services/auth';
import './AuthPage.css';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  async function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError('');

    if (!email || !password) {
      setError('Please fill in both fields.');
      return;
    }

    setLoading(true);
    try {
      await loginUser(email, password);

      // Check if the profile is complete, and route accordingly:
      // complete → Discover, incomplete → profile setup.
      const profile = await getMyProfile();
      if (isProfileComplete(profile)) {
        navigate('/app/discover');
      } else {
        navigate('/app/profile-setup');
      }
    } catch (err) {
      setError('Invalid email or password.');
      console.error(err);
    } finally {
      setLoading(false);
    }
  }

  return (
    <div className="auth">
      <div className="auth__wrapper">
        <div className="auth__logo">
          <div className="auth__logo-mark">P</div>
          <span className="auth__logo-text">pulse</span>
        </div>

        <div className="card auth__card">
          <div className="auth__header">
            <h1 className="auth__title">Welcome back</h1>
            <p className="auth__subtitle">Log in to keep training together.</p>
          </div>

          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label className="form-label" htmlFor="email">Email</label>
              <input
                id="email"
                type="email"
                className={`input ${error ? 'input-error' : ''}`}
                placeholder="you@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                disabled={loading}
              />
            </div>

            <div className="form-group">
              <div className="auth__row">
                <label className="form-label" htmlFor="password" style={{ marginBottom: 0 }}>
                  Password
                </label>
                <span className="btn btn-ghost btn-small">Forgot password?</span>
              </div>
              <input
                id="password"
                type="password"
                className={`input ${error ? 'input-error' : ''}`}
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                disabled={loading}
              />
              {error && <p className="form-error-msg">{error}</p>}
            </div>

            <button type="submit" className="btn btn-primary btn-full btn-large mt-md" disabled={loading}>
              {loading ? 'Logging in...' : 'Log in'}

            </button>
          </form>

          <div className="auth__footer">
            Don't have an account?{' '}
            <Link to="/register" className="auth__link">Sign up</Link>
          </div>
        </div>
      </div>
    </div>
  );
}