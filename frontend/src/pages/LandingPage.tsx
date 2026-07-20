import './LandingPage.css';

const SPORTS = ['Running', 'CrossFit', 'Padel', 'Gym', 'Yoga', 'Climbing', 'Tennis', 'Swimming'];

function MatchScoreRing({ score, size = 40 }: { score: number; size?: number }) {
  const radius = (size / 2) - 4;
  const circumference = 2 * Math.PI * radius;
  const offset = circumference - (score / 100) * circumference;

  return (
    <div className="match-score" style={{ width: size, height: size }}>
      <svg
        width={size}
        height={size}
        viewBox={`0 0 ${size} ${size}`}
        style={{ transform: 'rotate(-90deg)' }}
      >
        <circle
          cx={size / 2}
          cy={size / 2}
          r={radius}
          fill="none"
          stroke="var(--preview-border)"
          strokeWidth="3"
        />
        <circle
          cx={size / 2}
          cy={size / 2}
          r={radius}
          fill="none"
          stroke="var(--accent)"
          strokeWidth="3"
          strokeDasharray={circumference}
          strokeDashoffset={offset}
          strokeLinecap="round"
        />
      </svg>
      <span className="match-score__value">{score}</span>
    </div>
  );
}

export default function LandingPage() {
  return (
    <div className="landing">
      {/* Navigation */}
      <nav className="landing__nav">
        <div className="landing__logo">
          <div className="landing__logo-mark">
            <span className="landing__logo-letter">P</span>
          </div>
          <span className="landing__logo-text">pulse</span>
        </div>

        <div className="landing__nav-actions">
          <button className="btn btn-outline hide-mobile">Log in</button>
          <button className="btn btn-primary">Sign up free</button>
        </div>
      </nav>

      {/* Hero */}
      <section className="landing__hero">
        <div className="landing__hero-content">
          <h1 className="landing__headline">
            Don't train<br />
            alone<span className="landing__headline-accent">.</span>
          </h1>

          <p className="landing__tagline">
            Connect with people who move like you do. Run, train, play.
            Find your teammate — together is better.
          </p>

          <div className="landing__cta">
            <button className="btn btn-primary btn-large btn-full">
              Get started
            </button>
          </div>

          <div className="landing__stats">
            <div>
              <div className="landing__stat-value landing__stat-value--accent">20+</div>
              <div className="landing__stat-label">Sports</div>
            </div>
            <div>
              <div className="landing__stat-value landing__stat-value--white">2.4K</div>
              <div className="landing__stat-label">Athletes</div>
            </div>
            <div>
              <div className="landing__stat-value landing__stat-value--white">500+</div>
              <div className="landing__stat-label">Teams formed</div>
            </div>
          </div>
        </div>

        {/* Preview Panel — Desktop only */}
        <aside className="landing__preview">
          <div className="preview-card">
            <p className="preview-card__label">Preview: your best match</p>

            <div className="preview-card__user">
              <div className="avatar avatar-md preview-card__avatar">MK</div>
              <div className="preview-card__info">
                <div className="preview-card__name">
                  Marcus K. <span className="online-dot" />
                </div>
                <div className="preview-card__sport">Running · 5K · Mornings</div>
              </div>
              <MatchScoreRing score={92} />
            </div>

            <div className="preview-card__tags">
              <span className="tag">Mon-Fri</span>
              <span className="tag">5-7 AM</span>
              <span className="tag">Intermediate</span>
            </div>

            <button className="btn btn-primary btn-full">Connect</button>
          </div>

          <div className="preview-sports">
            {SPORTS.map((sport) => (
              <span key={sport} className="pill">{sport}</span>
            ))}
          </div>
        </aside>
      </section>

      {/* Features */}
      <section className="landing__features">
        <div className="feature">
          <div className="feature__header">
            <div className="feature__icon">🎯</div>
            <h3 className="feature__title">Smart matching</h3>
          </div>
          <p className="feature__description">
            Our algorithm finds people who match your sport, level, schedule, and goals.
          </p>
        </div>

        <div className="feature">
          <div className="feature__header">
            <div className="feature__icon">💬</div>
            <h3 className="feature__title">Real-time chat</h3>
          </div>
          <p className="feature__description">
            Connect and coordinate sessions instantly. No waiting, no refreshing.
          </p>
        </div>

        <div className="feature">
          <div className="feature__header">
            <div className="feature__icon">📍</div>
            <h3 className="feature__title">Location aware</h3>
          </div>
          <p className="feature__description">
            Only see athletes near you. No impractical matches across the globe.
          </p>
        </div>
      </section>
    </div>
  );
}