import { Outlet, NavLink, useNavigate } from 'react-router-dom';
import { Compass, Users, MessageCircle, User, LogOut } from 'lucide-react';
import { logout } from '../services/auth';
import './AppLayout.css';

const NAV = [
  { to: '/app/discover',    Icon: Compass,       label: 'Discover' },
  { to: '/app/connections', Icon: Users,         label: 'Connections' },
  { to: '/app/chats',       Icon: MessageCircle, label: 'Chats' },
  { to: '/app/profile',     Icon: User,          label: 'Profile' },
];

export default function AppLayout() {
  const navigate = useNavigate();

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
      <aside className="sidebar hide-mobile">
        {/* Logo */}
        <div className="logo-mark">P</div>

        {/* Navigation */}
        <div className="sidebar__nav">
          {NAV.map(({ to, Icon, label }) => (
            <NavLink
              key={to}
              to={to}
              className={({ isActive }) =>
                `nav-item ${isActive ? 'nav-item-active' : ''}`
              }
              title={label}
            >
              <Icon size={28} strokeWidth={2.5} />
            </NavLink>
          ))}
        </div>

        {/* Logout (pushed to bottom) */}
        <button className="nav-item sidebar__logout" onClick={handleLogout} title="Log out">
          <LogOut size={28} strokeWidth={2.5} />
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