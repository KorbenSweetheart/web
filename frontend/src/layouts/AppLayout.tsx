import { Outlet, NavLink } from 'react-router-dom';
import './AppLayout.css';

const NAV = [
  { to: '/app/discover',    icon: '🔥', label: 'Discover' },
  { to: '/app/connections', icon: '👥', label: 'Connections' },
  { to: '/app/chats',       icon: '💬', label: 'Chats' },
  { to: '/app/profile',     icon: '👤', label: 'Profile' },
];

export default function AppLayout() {
  return (
    <div className="app-layout">
      {/* Desktop: sidebar */}
      <aside className="sidebar hide-mobile">
        {NAV.map(item => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              `nav-item ${isActive ? 'nav-item-active' : ''}`
            }
            title={item.label}
          >
            {item.icon}
          </NavLink>
        ))}
      </aside>

      {/* Contenido de la página hija */}
      <main className="app-main">
        <Outlet />
      </main>

      {/* Mobile: bottom nav */}
      <nav className="bottom-nav hide-desktop">
        {NAV.map(item => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              `bottom-nav__item ${isActive ? 'bottom-nav__item--active' : ''}`
            }
          >
            <span>{item.icon}</span>
            <span className="bottom-nav__label">{item.label}</span>
          </NavLink>
        ))}
      </nav>
    </div>
  );
}