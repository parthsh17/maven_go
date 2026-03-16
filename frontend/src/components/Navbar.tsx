import { Link, useLocation, useNavigate } from 'react-router-dom';

const navLinks = [
  { to: '/dashboard', label: 'Dashboard' },
  { to: '/orders', label: 'Orders' },
  { to: '/submit', label: 'Submit Order' },
  { to: '/metrics', label: 'Metrics' },
];

export function Navbar() {
  const location = useLocation();
  const navigate = useNavigate();

  const handleLogout = () => {
    localStorage.removeItem('maven_auth');
    navigate('/login');
  };

  return (
    <nav className="navbar">
      <div className="navbar-brand">
        <span className="brand-icon">◈</span>
        <span className="brand-name">Maven</span>
      </div>
      <ul className="navbar-links">
        {navLinks.map((l) => (
          <li key={l.to}>
            <Link
              to={l.to}
              className={`nav-link ${location.pathname.startsWith(l.to) ? 'active' : ''}`}
            >
              {l.label}
            </Link>
          </li>
        ))}
      </ul>
      <button className="btn btn-ghost logout-btn" onClick={handleLogout}>
        Sign Out
      </button>
    </nav>
  );
}
