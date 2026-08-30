import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { Compass, Users, MessageCircle, User, LogOut, PanelLeft, PanelLeftClose } from 'lucide-react';
import { logout } from '../services/auth';
import { useState, useEffect } from 'react';
import { updateMyLocation } from '../services/users';
import './AppLayout.css';

const NAV = [
  { to: '/app/discover',    Icon: Compass,       label: 'Discover' },
  { to: '/app/connections', Icon: Users,         label: 'Connections' },
  { to: '/app/chats',       Icon: MessageCircle, label: 'Chats' },
  { to: '/app/profile',     Icon: User,          label: 'Profile' },
];

export default function AppLayout() {
  const navigate = useNavigate();
  const [expanded, setExpanded] = useState(false);
  // Send the user's location once, when they enter the app — not before
  // every Discover fetch (that added a 1-5s browser-geolocation delay to
  // each visit). AppLayout mounts once for the whole logged-in area, so
  // this runs a single time per session.
  useEffect(() => {
    updateMyLocation();
  }, []);

  async function handleLogout() {
    await logout();
    navigate('/login');
  }

  return (
    <div className="app-layout">
       {/* Mobile: top header (logo + logout) */}
      <header className="app-header hide-desktop">
        <div className="logo">
          <div className="logo-mark">P</div>
        </div>
        <button className="app-header__logout" onClick={handleLogout} title="Log out">
          <LogOut size={22} strokeWidth={2.5} />
        </button>
      </header>

      {/* Desktop: sidebar */}
      <aside className={`sidebar hide-mobile ${expanded ? 'is-expanded' : ''}`}>
         {/* Logo */}
        <div className="logo-mark">P</div>
        {/* Toggle button */}
        <button
          className="sidebar__toggle"
          onClick={() => setExpanded((v) => !v)}
          title={expanded ? 'Collapse menu' : 'Expand menu'}
          aria-label={expanded ? 'Collapse menu' : 'Expand menu'}
        >
          {expanded ? <PanelLeftClose size={24} strokeWidth={2.5} /> : <PanelLeft size={24} strokeWidth={2.5} />}
        </button>

        {/* Navigation */}
        <div className="sidebar__nav">
          {NAV.map(({ to, Icon, label }) => (
            <NavLink
              key={to}
              to={to}
              className={({ isActive }) => `nav-item ${isActive ? 'nav-item-active' : ''}`}
              title={label}
            >
              <Icon size={28} strokeWidth={2.5} />
              <span className="nav-item__label">{label}</span>
            </NavLink>
          ))}
        </div>

        {/* Logout (pushed to bottom) */}
        <button className="nav-item sidebar__logout" onClick={handleLogout} title="Log out">
          <LogOut size={28} strokeWidth={2.5} />
          <span className="nav-item__label">Log out</span>
        </button>
      </aside>

      {/* Contenido de la página hija */}
      <main className="app-main">
        <Outlet />
      </main>

      {/* Mobile: bottom nav */}
      <nav className="bottom-nav hide-desktop">
        {NAV.map(({ to, Icon, label }) => (
          <NavLink
            key={to}
            to={to}
            className={({ isActive }) =>
              `bottom-nav__item ${isActive ? 'bottom-nav__item--active' : ''}`
            }
          >
            <Icon size={22} strokeWidth={2.5} />
            <span className="bottom-nav__label">{label}</span>
          </NavLink>
        ))}
      </nav>
    </div>
  );
}