import './LandingPage.css';
import { Link } from "react-router-dom";
import { Target, MessageCircle, MapPin } from 'lucide-react';

function MatchScoreRing({ score, size = 44 }: { score: number; size?: number }) {
  const radius = (size / 2) - 4;
  const circumference = 2 * Math.PI * radius;
  const offset = circumference - (score / 100) * circumference;

  return (
    <div className="match-score" style={{ width: size, height: size }}>
      <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`} style={{ transform: 'rotate(-90deg)' }}>
        <circle cx={size / 2} cy={size / 2} r={radius} fill="none" stroke="rgba(255,255,255,0.1)" strokeWidth="3" />
        <circle cx={size / 2} cy={size / 2} r={radius} fill="none" stroke="#F97316" strokeWidth="3"
          strokeDasharray={circumference} strokeDashoffset={offset} strokeLinecap="round" />
      </svg>
      <span className="match-score__value">{score}</span>
    </div>
  );
}

export default function LandingPage() {
  return (
    <div className="landing">

      {/* ===== NAVBAR ===== */}
      <nav className="landing-nav">
        <div className="logo">
          <div className="logo-mark">P</div>
          <span className="logo-text">pulse</span>
        </div>
        <div className="flex items-center gap-sm">
          <Link to="/login" className="btn btn-outline hide-mobile">Log in</Link>
          <Link to="/register" className="btn btn-primary">Sign up free</Link>
        </div>
      </nav>

      {/* ===== HERO: two equal columns ===== */}
      <section className="landing-hero">
        <div className="landing-hero__content">
          <h1 className="text-hero">
            Don't train<br />alone<span className="text-accent">.</span>
          </h1>

          <p className="text-body mt-md landing-tagline">
            It's like dating apps, but the only thing you're committing to is a 6am run. Matched on sport, level, and vibe.
          </p>

          <div className="landing-hero__cta">
            <Link to="/register" className="btn btn-primary btn-full">Get started — it's free</Link>
          </div>

          <div className="landing-stats">
            <div>
              <div className="stat-value stat-value--accent">20+</div>
              <div className="stat-label">Sports</div>
            </div>
            <div>
              <div className="stat-value text-white">2.4K</div>
              <div className="stat-label">Athletes</div>
            </div>
            <div>
              <div className="stat-value text-white">500+</div>
              <div className="stat-label">Teams formed</div>
            </div>
          </div>
        </div>

        {/* ===== PREVIEW ===== */}
        <aside className="landing-preview">
          <div className="preview-card">
            <p className="preview-label">Preview: your best match</p>

            <div className="preview-user">
              <div className="preview-avatar">MK</div>
              <div style={{ flex: 1 }}>
                <div className="preview-name">
                  Marcus K. <span className="online-dot" />
                </div>
                <div className="preview-sport">Running · 5K · Mornings</div>
              </div>
              <MatchScoreRing score={92} />
            </div>

            <div className="preview-tags">
              <span className="preview-tag">🎧 Silent</span>
              <span className="preview-tag">Mon-Fri</span>
              <span className="preview-tag">Long-term</span>
            </div>

            <button className="btn btn-primary btn-full">Connect</button>

            {/* Bio data points as subtle pills inside the card */}
            <div className="preview-bio">
              <div className="preview-bio__pill"><span>🏃</span> Running</div>
              <div className="preview-bio__pill"><span>📊</span> Intermediate</div>
              <div className="preview-bio__pill"><span>📍</span> Kuopio</div>
              <div className="preview-bio__pill"><span>🎧</span> Silent</div>
              <div className="preview-bio__pill"><span>🔁</span> Long-term</div>
            </div>
          </div>
        </aside>
      </section>

      {/* ===== FEATURES: icon + title inline ===== */}
      <section className="landing-features">
        <div className="feature">
          <div className="feature-header">
            <div className="feature-icon feature-icon--match">
              <Target size={22} strokeWidth={2.5} />
            </div>
            <h3 className="feature-title">Matched on what matters</h3>
          </div>
          <p className="feature-desc">
            Sport, skill, training mode, location, commitment. Five signals
            working together — so every match is someone you'd actually train with.
          </p>
        </div>

        <div className="feature">
          <div className="feature-header">
            <div className="feature-icon feature-icon--chat">
              <MessageCircle size={22} strokeWidth={2.5} />
            </div>
            <h3 className="feature-title">Chat in real time</h3>
          </div>
          <p className="feature-desc">
            Connected? Start talking. Plan your next session, swap routes, or just
            hype each other up. Messages land the instant they're sent.
          </p>
        </div>

        <div className="feature">
          <div className="feature-header">
            <div className="feature-icon feature-icon--location">
              <MapPin size={22} strokeWidth={2.5} />
            </div>
            <h3 className="feature-title">People near you</h3>
          </div>
          <p className="feature-desc">
            Set your radius and only meet people close enough to actually show up.
            No matches stranded on the other side of the map.
          </p>
        </div>
      </section>
    </div>
  );
}