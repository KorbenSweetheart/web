import { useState, FormEvent } from 'react';
import './AuthPage.css';

interface RegisterPageProps {
  onSwitchToLogin?: () => void;
}

export default function RegisterPage({ onSwitchToLogin }: RegisterPageProps) {
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [acceptedTerms, setAcceptedTerms] = useState(false);
  const [error, setError] = useState('');

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    setError('');

    if (!name || !email || !password || !confirmPassword) {
      setError('Please fill in all fields.');
      return;
    }
    if (password !== confirmPassword) {
      setError('Passwords don\'t match.');
      return;
    }
    if (!acceptedTerms) {
      setError('You need to accept the terms to continue.');
      return;
    }

    // TODO: hook up to auth API
    console.log('Registering', { name, email, password });
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
            <h1 className="auth__title">Create your account</h1>
            <p className="auth__subtitle">Find your next training partner.</p>
          </div>

          <form onSubmit={handleSubmit}>
            <div className="form-group">
              <label className="form-label" htmlFor="name">Full name</label>
              <input
                id="name"
                type="text"
                className="input"
                placeholder="Marcus K."
                value={name}
                onChange={(e) => setName(e.target.value)}
              />
            </div>

            <div className="form-group">
              <label className="form-label" htmlFor="email">Email</label>
              <input
                id="email"
                type="email"
                className="input"
                placeholder="you@example.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>

            <div className="form-group">
              <label className="form-label" htmlFor="password">Password</label>
              <input
                id="password"
                type="password"
                className="input"
                placeholder="••••••••"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
              />
              <p className="form-helper">At least 8 characters.</p>
            </div>

            <div className="form-group">
              <label className="form-label" htmlFor="confirmPassword">Confirm password</label>
              <input
                id="confirmPassword"
                type="password"
                className={`input ${error ? 'input-error' : ''}`}
                placeholder="••••••••"
                value={confirmPassword}
                onChange={(e) => setConfirmPassword(e.target.value)}
              />
            </div>

            <div className="form-group">
              <label className="checkbox-item">
                <input
                  type="checkbox"
                  checked={acceptedTerms}
                  onChange={(e) => setAcceptedTerms(e.target.checked)}
                />
                <span className="text-caption">
                  I agree to the Terms of Service and Privacy Policy
                </span>
              </label>
              {error && <p className="form-error-msg">{error}</p>}
            </div>

            <button type="submit" className="btn btn-primary btn-full btn-large mt-md">
              Sign up
            </button>
          </form>

          <div className="auth__footer">
            Already have an account?{' '}
            <span className="auth__link" onClick={onSwitchToLogin}>
              Log in
            </span>
          </div>
        </div>
      </div>
    </div>
  );
}