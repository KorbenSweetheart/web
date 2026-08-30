import { Outlet, NavLink, useNavigate, useLocation } from 'react-router-dom';
import { Compass, Users, MessageCircle, User, LogOut, PanelLeft, PanelLeftClose } from 'lucide-react';
import { logout } from '../services/auth';
import { useState, useEffect, useCallback, useRef } from 'react';
import { updateMyLocation, getMyProfile } from '../services/users';
import { getConnectionRequestsCount } from '../services/connections';
import { getUnreadMessagesCount } from '../services/chats';
import { getChatSocket } from '../services/websocket';
import './AppLayout.css';

export default function AppLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const [expanded, setExpanded] = useState(false);
  const [pendingRequestsCount, setPendingRequestsCount] = useState(0);
  const [unreadMessagesCount, setUnreadMessagesCount] = useState(0);
  const myIdRef = useRef<number | null>(null);

  const fetchRequestCount = useCallback(async () => {
    try {
      const count = await getConnectionRequestsCount();
      setPendingRequestsCount(count);
    } catch {
      // silently ignore if unauthenticated or network error
    }
  }, []);

  const fetchUnreadChats = useCallback(async () => {
    try {
      const count = await getUnreadMessagesCount(myIdRef.current ?? undefined);
      setUnreadMessagesCount(count);
    } catch {
      // silently ignore if unauthenticated or network error
    }
  }, []);

  // Send the user's location once, when they enter the app — not before
  // every Discover fetch (that added a 1-5s browser-geolocation delay to
  // each visit). AppLayout mounts once for the whole logged-in area, so
  // this runs a single time per session.
  useEffect(() => {
    updateMyLocation();
    getMyProfile().then((me) => {
      myIdRef.current = me.id;
      fetchUnreadChats();
    }).catch(() => {});
  }, [fetchUnreadChats]);

  useEffect(() => {
    fetchRequestCount();
    fetchUnreadChats();
  }, [fetchRequestCount, fetchUnreadChats, location.pathname]);

  useEffect(() => {
    const handleUpdate = () => {
      fetchRequestCount();
    };

    const handleChatsUpdate = () => {
      fetchUnreadChats();
    };

    window.addEventListener('connections:updated', handleUpdate);
    window.addEventListener('chats:unread-changed', handleChatsUpdate);
    window.addEventListener('focus', handleUpdate);
    window.addEventListener('focus', handleChatsUpdate);

    const socket = getChatSocket();
    const offSocket = socket.onMessage((msg) => {
      if (myIdRef.current != null && msg.sender_id !== myIdRef.current) {
        fetchUnreadChats();
      }
    });

    const interval = setInterval(() => {
      fetchRequestCount();
      fetchUnreadChats();
    }, 30000);

    return () => {
      window.removeEventListener('connections:updated', handleUpdate);
      window.removeEventListener('chats:unread-changed', handleChatsUpdate);
      window.removeEventListener('focus', handleUpdate);
      window.removeEventListener('focus', handleChatsUpdate);
      offSocket();
      clearInterval(interval);
    };
  }, [fetchRequestCount, fetchUnreadChats]);

  async function handleLogout() {
    await logout();
    navigate('/login');
  }

  const navItems = [
    { to: '/app/discover',    Icon: Compass,       label: 'Discover' },
    {
      to: '/app/connections',
      Icon: Users,
      label: 'Connections',
      badge: pendingRequestsCount > 0 ? pendingRequestsCount : undefined,
    },
    {
      to: '/app/chats',
      Icon: MessageCircle,
      label: 'Chats',
      badge: unreadMessagesCount > 0 ? unreadMessagesCount : undefined,
      badgeVariant: 'accent',
    },
    { to: '/app/profile',     Icon: User,          label: 'Profile' },
  ];

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
          {navItems.map(({ to, Icon, label, badge, badgeVariant }) => (
            <NavLink
              key={to}
              to={to}
              className={({ isActive }) => `nav-item ${isActive ? 'nav-item-active' : ''}`}
              title={label}
            >
              <div className="nav-item__icon-wrap">
                <Icon size={28} strokeWidth={2.5} />
                {badge != null && (
                  <span
                    className={`nav-item__badge ${badgeVariant ? `nav-item__badge--${badgeVariant}` : ''}`}
                    aria-label={`${badge} unread items`}
                  >
                    {badge}
                  </span>
                )}
              </div>
              <span className="nav-item__label">{label}</span>
              {badge != null && (
                <span
                  className={`nav-item__badge-expanded ${badgeVariant ? `nav-item__badge-expanded--${badgeVariant}` : ''}`}
                  aria-label={`${badge} unread items`}
                >
                  {badge}
                </span>
              )}
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
        {navItems.map(({ to, Icon, label, badge, badgeVariant }) => (
          <NavLink
            key={to}
            to={to}
            className={({ isActive }) =>
              `bottom-nav__item ${isActive ? 'bottom-nav__item--active' : ''}`
            }
          >
            <div className="bottom-nav__icon-wrap">
              <Icon size={22} strokeWidth={2.5} />
              {badge != null && (
                <span
                  className={`bottom-nav__badge ${badgeVariant ? `bottom-nav__badge--${badgeVariant}` : ''}`}
                  aria-label={`${badge} unread items`}
                >
                  {badge}
                </span>
              )}
            </div>
            <span className="bottom-nav__label">{label}</span>
          </NavLink>
        ))}
      </nav>
    </div>
  );
}