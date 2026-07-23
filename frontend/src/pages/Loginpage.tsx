import { useState, FormEvent } from 'react';
import './AuthPage.css';

interface LoginPageProps {
  onSwitchToRegister?: () => void;
}

export default function LoginPage({ onSwitchToRegister }: LoginPageProps) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError('');

    if (!email || !password) {
      setError('Please fill in both fields.');
      return;
    }

    // TODO: hook up to auth API
    console.log('Logging in with', { email, password });
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
              />
              {error && <p className="form-error-msg">{error}</p>}
            </div>

            <button type="submit" className="btn btn-primary btn-full btn-large mt-md">
              Log in
            </button>
          </form>

          <div className="auth__footer">
            Don't have an account?{' '}
            <span className="auth__link" onClick={onSwitchToRegister}>
              Sign up
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}